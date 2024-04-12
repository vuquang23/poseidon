package task

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hibiken/asynq"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
	"gorm.io/datatypes"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
	poolrepo "github.com/vuquang23/poseidon/internal/pkg/repository/pool"
	pricerepo "github.com/vuquang23/poseidon/internal/pkg/repository/price"
	txrepo "github.com/vuquang23/poseidon/internal/pkg/repository/tx"
	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
	asynqpkg "github.com/vuquang23/poseidon/pkg/asynq"
	"github.com/vuquang23/poseidon/pkg/binance"
	"github.com/vuquang23/poseidon/pkg/eth"
	"github.com/vuquang23/poseidon/pkg/logger"
	"github.com/vuquang23/poseidon/pkg/uniswapv3"
)

type TaskService struct {
	config        Config
	poolRepo      poolrepo.IPoolRepository
	txRepo        txrepo.ITxRepository
	priceRepo     pricerepo.IPriceRepository
	ethClient     eth.IClient
	asynqClient   asynqpkg.IAsynqClient
	binanceClient binance.IClient
}

func New(
	config Config,
	poolRepo poolrepo.IPoolRepository,
	txRepo txrepo.ITxRepository,
	priceRepo pricerepo.IPriceRepository,
	ethClient eth.IClient,
	asynqClient asynqpkg.IAsynqClient,
	binanceClient binance.IClient,
) *TaskService {
	return &TaskService{
		config:        config,
		poolRepo:      poolRepo,
		txRepo:        txRepo,
		priceRepo:     priceRepo,
		ethClient:     ethClient,
		asynqClient:   asynqClient,
		binanceClient: binanceClient,
	}
}

// SetEthClient sets mock eth client for testing.
func (s *TaskService) SetEthClient(c eth.IClient) {
	s.ethClient = c
}

// SetAsynqClient sets mock asynq client for testing.
func (s *TaskService) SetAsynqClient(c asynqpkg.IAsynqClient) {
	s.asynqClient = c
}

func (s *TaskService) GetPeriodicTaskConfigs(ctx context.Context) ([]*asynq.PeriodicTaskConfig, error) {
	pools, err := s.poolRepo.GetPools(ctx)
	if err != nil {
		return nil, err
	}

	var configs []*asynq.PeriodicTaskConfig

	scanTxsTaskConfigs, err := s.initScanTxsTaskConfigs(ctx, pools)
	if err != nil {
		return nil, err
	}

	finalizerTxsTaskConfigs, err := s.initFinalizeTxsTaskConfigs(ctx, pools)
	if err != nil {
		return nil, err
	}

	configs = append(configs, scanTxsTaskConfigs...)
	configs = append(configs, finalizerTxsTaskConfigs...)

	return configs, nil
}

func (s *TaskService) initFinalizeTxsTaskConfigs(ctx context.Context, pools []*entity.Pool) ([]*asynq.PeriodicTaskConfig, error) {
	var configs []*asynq.PeriodicTaskConfig
	for _, p := range pools {
		payload := valueobject.TaskFinalizeTxsPayload{
			PoolID:         p.ID,
			PoolAddress:    p.Address,
			Token0Decimals: p.Token0Decimals,
			Token1Decimals: p.Token1Decimals,
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			logger.Error(ctx, err.Error())
			return nil, err
		}

		t := asynq.NewTask(
			valueobject.TaskTypeFinalizeTxs,
			payloadBytes,
			asynq.TaskID(fmt.Sprintf("%s:%s", valueobject.TaskTypeFinalizeTxs, p.Address)),
		)

		configs = append(configs, &asynq.PeriodicTaskConfig{
			Cronspec: s.config.Cronspec,
			Task:     t,
		})
	}

	return configs, nil
}

func (s *TaskService) initScanTxsTaskConfigs(ctx context.Context, pools []*entity.Pool) ([]*asynq.PeriodicTaskConfig, error) {
	var configs []*asynq.PeriodicTaskConfig
	for _, p := range pools {
		payload := valueobject.TaskScanTxsPayload{
			PoolID:         p.ID,
			PoolAddress:    p.Address,
			Token0Decimals: p.Token0Decimals,
			Token1Decimals: p.Token1Decimals,
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			logger.Error(ctx, err.Error())
			return nil, err
		}

		t := asynq.NewTask(
			valueobject.TaskTypeScanTxs,
			payloadBytes,
			asynq.TaskID(fmt.Sprintf("%s:%s", valueobject.TaskTypeScanTxs, p.Address)),
		)

		configs = append(configs, &asynq.PeriodicTaskConfig{
			Cronspec: s.config.Cronspec,
			Task:     t,
		})
	}

	return configs, nil
}

func (s *TaskService) HandlePoolCreated(ctx context.Context, poolAddress string) error {
	pool, err := s.poolRepo.GetPoolByAddress(ctx, poolAddress)
	if err != nil {
		return err
	}

	// scanner cursor
	scannerCursor := entity.BlockCursor{
		PoolID:      pool.ID,
		Type:        valueobject.BlockCursorTypeScanner,
		BlockNumber: pool.StartBlock,
	}

	// finalizer cursor
	latestFinalizedBlockNbr, err := s.GetLatestFinalizedBlockNumber(ctx)
	if err != nil {
		return err
	}
	extra := valueobject.BlockCursorFinalizerExtra{
		CreatedAtFinalizedBlock: latestFinalizedBlockNbr,
	}
	extraBytes, err := json.Marshal(extra)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	var extraJSON datatypes.JSON
	if err := extraJSON.UnmarshalJSON(extraBytes); err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	finalizerCursor := entity.BlockCursor{
		PoolID:      pool.ID,
		Type:        valueobject.BlockCursorTypeFinalizer,
		BlockNumber: pool.StartBlock,
		Extra:       extraJSON,
	}

	// persist cursors to DB
	return s.txRepo.CreateBlockCursors(ctx, []*entity.BlockCursor{&scannerCursor, &finalizerCursor})
}

func (s *TaskService) ScanTxs(ctx context.Context, task valueobject.TaskScanTxsPayload) error {
	var (
		poolID         = task.PoolID
		poolAddress    = task.PoolAddress
		token0Decimals = task.Token0Decimals
		token1Decimals = task.Token1Decimals
	)

	cursor, err := s.txRepo.GetCursorByPoolIDAndType(ctx, poolID, valueobject.BlockCursorTypeScanner)
	if err != nil {
		return err
	}

	logs, fromBlock, toBlock, err := s.getLogs(ctx, cursor.BlockNumber, poolAddress)
	if err != nil {
		return err
	}

	logger.WithFields(ctx, logger.Fields{
		"fromBlock": fromBlock,
		"toBlock":   toBlock,
		"pool":      poolAddress,
		"logs":      len(logs),
	}).Info("retrieve logs")

	blockHashes := uniqueBlockHashes(logs)
	headers, err := s.getBlockHeaders(ctx, blockHashes)
	if err != nil {
		return err
	}

	if err := s.enqueueTaskGetETHUSDTKlines(ctx, headers); err != nil {
		return err
	}

	txHashes := uniqueTxHashes(logs)
	txReceipts, err := s.getTxs(ctx, txHashes)
	if err != nil {
		return err
	}

	txs, err := initTxs(ctx, poolID, headers, txReceipts)
	if err != nil {
		return err
	}

	swapEvents, err := initSwapEvents(ctx, poolID, token0Decimals, token1Decimals, logs)
	if err != nil {
		return err
	}

	return s.txRepo.UpdateDataScanner(ctx, cursor.ID, toBlock+1, txs, swapEvents)
}

func (s *TaskService) GetETHUSDTKline(ctx context.Context, payload valueobject.TaskGetETHUSDTKlinePayload) error {
	starTimeNsec := int64(payload.Time) * 1000
	klines, err := s.binanceClient.GetKlines(ctx, starTimeNsec, 0, 1, "ETHUSDT", "1m")
	if err != nil {
		return err
	}

	if len(klines) == 0 {
		logger.WithFields(ctx, logger.Fields{
			"startTimeNsec": starTimeNsec,
		}).Error(ErrEmptyKlines.Error())
		return ErrEmptyKlines
	}

	kline := klines[0]

	ohlc4 := decimal.NewFromInt(0).
		Add(decimal.RequireFromString(kline.Open)).
		Add(decimal.RequireFromString(kline.High)).
		Add(decimal.RequireFromString(kline.Low)).
		Add(decimal.RequireFromString(kline.Close)).
		Div(decimal.NewFromInt(4)).Round(6)

	e := entity.ETHUSDTKline{
		OpenTime:   uint64(kline.OpenTime),
		CloseTime:  uint64(kline.CloseTime),
		OpenPrice:  kline.Open,
		HighPrice:  kline.High,
		LowPrice:   kline.Low,
		ClosePrice: kline.Close,
		OHLC4:      ohlc4,
	}

	return s.priceRepo.CreateKline(ctx, &e)
}

func (s *TaskService) FinalizeTxs(ctx context.Context, payload valueobject.TaskFinalizeTxsPayload) error {
	var (
		poolID         = payload.PoolID
		poolAddress    = payload.PoolAddress
		token0Decimals = payload.Token0Decimals
		token1Decimals = payload.Token1Decimals
	)

	finalizerCursor, err := s.txRepo.GetCursorByPoolIDAndType(ctx, poolID, valueobject.BlockCursorTypeFinalizer)
	if err != nil {
		return err
	}
	scannerCursor, err := s.txRepo.GetCursorByPoolIDAndType(ctx, poolID, valueobject.BlockCursorTypeScanner)
	if err != nil {
		return err
	}

	fromBlock, toBlock, finalizerCreatedAtFinalizedBlockNbr, err := s.initFinalizerBlockRange(ctx, finalizerCursor, scannerCursor)
	if err != nil {
		return err
	}
	if toBlock <= finalizerCreatedAtFinalizedBlockNbr {
		return s.txRepo.UpdateDataFinalizer(ctx, poolID, finalizerCursor.ID, fromBlock, toBlock, nil, nil)
	}

	logger.WithFields(ctx, logger.Fields{
		"poolId":    poolID,
		"fromBlock": fromBlock,
		"toBlock":   toBlock,
	}).Info("finalize txs")

	logs, err := s.ethClient.GetLogs(ctx, fromBlock, toBlock, []common.Address{common.HexToAddress(poolAddress)})
	if err != nil {
		return err
	}

	txHashes := uniqueTxHashes(logs)

	existingTxs, err := s.txRepo.GetTxsByPoolIDAndBlockRange(ctx, poolID, fromBlock, toBlock)
	if err != nil {
		return err
	}

	reorg := compareFinalizedTxsWithExistingTxs(txHashes, existingTxs)
	if !reorg {
		return s.txRepo.UpdateDataFinalizer(
			ctx, poolID, finalizerCursor.ID, fromBlock, toBlock, nil, nil,
		)
	}

	blockHashes := uniqueBlockHashes(logs)
	headers, err := s.getBlockHeaders(ctx, blockHashes)
	if err != nil {
		return err
	}

	receipts, err := s.getTxs(ctx, txHashes)
	if err != nil {
		return err
	}

	txs, err := initTxs(ctx, poolID, headers, receipts)
	if err != nil {
		return err
	}
	for _, tx := range txs {
		tx.IsFinalized = true
	}

	swapEvents, err := initSwapEvents(ctx, poolID, token0Decimals, token1Decimals, logs)
	if err != nil {
		return err
	}

	return s.txRepo.UpdateDataFinalizer(ctx, poolID, finalizerCursor.ID, fromBlock, toBlock, txs, swapEvents)
}

func (s *TaskService) initFinalizerBlockRange(ctx context.Context, finalizerCursor, scannerCursor *entity.BlockCursor) (uint64, uint64, uint64, error) {
	var extra valueobject.BlockCursorFinalizerExtra
	extraBytes, err := finalizerCursor.Extra.MarshalJSON()
	if err != nil {
		logger.Error(ctx, err.Error())
		return 0, 0, 0, err
	}
	if err := json.Unmarshal(extraBytes, &extra); err != nil {
		logger.Error(ctx, err.Error())
		return 0, 0, 0, err
	}

	fromBlock := finalizerCursor.BlockNumber
	toBlock := fromBlock + s.config.BlockBatchSize - 1

	latestFinalizedBlockNbr, err := s.GetLatestFinalizedBlockNumber(ctx)
	if err != nil {
		return 0, 0, 0, err
	}
	if toBlock > latestFinalizedBlockNbr {
		toBlock = latestFinalizedBlockNbr
	}

	if toBlock > scannerCursor.BlockNumber {
		toBlock = scannerCursor.BlockNumber
	}

	if fromBlock > toBlock { // node got issues
		logger.WithFields(ctx, logger.Fields{
			"fromBlock": fromBlock,
			"toBlock":   toBlock,
		}).Error(ErrInvalidBlockRange.Error())
		return 0, 0, 0, ErrInvalidBlockRange
	}

	if fromBlock <= extra.CreatedAtFinalizedBlock && extra.CreatedAtFinalizedBlock <= toBlock {
		toBlock = extra.CreatedAtFinalizedBlock
	}

	return fromBlock, toBlock, extra.CreatedAtFinalizedBlock, nil
}

// compareFinalizedTxsWithExistingTxs checks whether reorg occured.
func compareFinalizedTxsWithExistingTxs(finalizedTxHashes []common.Hash, existingTxs []*entity.Tx) bool {
	if len(finalizedTxHashes) != len(existingTxs) {
		return true
	}

	m := make(map[string]struct{})
	for _, h := range finalizedTxHashes {
		m[strings.ToLower(h.Hex())] = struct{}{}
	}

	for _, tx := range existingTxs {
		if _, ok := m[tx.TxHash]; !ok {
			return true
		}
	}

	return false
}

func (s *TaskService) getLogs(ctx context.Context, cursorBlockNbr uint64, poolAddress string) ([]types.Log, uint64, uint64, error) {
	latestBlockHeader, err := s.ethClient.GetLatestBlockHeader(ctx)
	if err != nil {
		return nil, 0, 0, err
	}
	latestBlockNbr := latestBlockHeader.Number.Uint64()

	fromBlock := cursorBlockNbr
	toBlock := fromBlock + s.config.BlockBatchSize - 1
	if toBlock > latestBlockNbr {
		toBlock = latestBlockNbr
	}

	if fromBlock > toBlock {
		logger.WithFields(ctx, logger.Fields{
			"poolAddress":    poolAddress,
			"fromBlock":      fromBlock,
			"toBlock":        toBlock,
			"batchSize":      s.config.BlockBatchSize,
			"latestBlockNbr": latestBlockNbr,
		}).Warn("invalid block range")
		return nil, 0, 0, ErrInvalidBlockRange
	}

	logs, err := s.ethClient.GetLogs(ctx, fromBlock, toBlock, []common.Address{common.HexToAddress(poolAddress)})
	if err != nil {
		return nil, 0, 0, err
	}

	return logs, fromBlock, toBlock, nil
}

func (s *TaskService) getTxs(ctx context.Context, txHashes []common.Hash) ([]*types.Receipt, error) {
	var (
		wg        errgroup.Group
		resultMap sync.Map
		receipts  = make([]*types.Receipt, len(txHashes))
	)

	for idx, txHash := range txHashes {
		_idx, _txHash := idx, txHash

		wg.Go(func() error {
			receipt, err := s.ethClient.GetTxReceipt(ctx, _txHash)
			if err != nil {
				return err
			}
			resultMap.Store(_idx, receipt)
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		logger.Error(ctx, err.Error())
		return nil, err
	}

	for i := 0; i < len(txHashes); i++ {
		r, _ := resultMap.Load(i)
		receipt := r.(*types.Receipt)
		receipts[i] = receipt
	}

	return receipts, nil
}

func (s *TaskService) getBlockHeaders(ctx context.Context, blockHashes []common.Hash) ([]*types.Header, error) {
	var (
		wg        errgroup.Group
		resultMap sync.Map
		headers   = make([]*types.Header, len(blockHashes))
	)

	for idx, blockHash := range blockHashes {
		_idx, _blockHash := idx, blockHash

		wg.Go(func() error {
			header, err := s.ethClient.HeaderByHash(ctx, _blockHash)
			if err != nil {
				return err
			}
			resultMap.Store(_idx, header)
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		logger.Error(ctx, err.Error())
		return nil, err
	}

	for i := 0; i < len(blockHashes); i++ {
		r, _ := resultMap.Load(i)
		header := r.(*types.Header)
		headers[i] = header
	}

	return headers, nil
}

func (s *TaskService) enqueueTaskGetETHUSDTKlines(ctx context.Context, blockHeaders []*types.Header) error {
	exists := map[uint64]struct{}{}
	for _, b := range blockHeaders {
		// round down to the timestamp starting this minute
		t := uint64(time.Unix(int64(b.Time), 0).Truncate(time.Minute).Unix())

		if _, ok := exists[t]; ok {
			continue
		}

		if err := s.asynqClient.EnqueueTask(
			ctx,
			valueobject.TaskTypeGetETHUSDTKline,
			"", "",
			valueobject.TaskGetETHUSDTKlinePayload{Time: t},
			-1,
		); err != nil {
			return err
		}

		exists[t] = struct{}{}
	}

	return nil
}

func (s *TaskService) GetLatestFinalizedBlockNumber(ctx context.Context) (uint64, error) {
	header, err := s.ethClient.GetLatestBlockHeader(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		return 0, err
	}

	t := time.Unix(int64(header.Time), 0)
	if t.Before(time.Now().Add(-s.config.BlockTimeDelayThreshold)) {
		logger.WithFields(ctx, logger.Fields{
			"blockTime":      header.Time,
			"delayThreshold": s.config.BlockTimeDelayThreshold,
		}).Error(ErrInvalidLatestBlock.Error())
		return 0, ErrInvalidLatestBlock
	}

	return header.Number.Uint64() - s.config.BlockFinalization, nil
}

func initTxs(ctx context.Context, poolID uint64, blockHeaders []*types.Header, txReceipts []*types.Receipt) ([]*entity.Tx, error) {
	headerByHash := map[string]*types.Header{}
	for _, b := range blockHeaders {
		h := strings.ToLower(b.Hash().Hex())
		headerByHash[h] = b
	}

	txs := make([]*entity.Tx, 0, len(txReceipts))

	for _, receipt := range txReceipts {
		txHash := strings.ToLower(receipt.TxHash.Hex())

		header, ok := headerByHash[strings.ToLower(receipt.BlockHash.Hex())]
		if !ok {
			logger.WithFields(ctx, logger.Fields{
				"blockHash": receipt.BlockHash.Hex(),
			}).Error(ErrBlockNotFound.Error())
			return nil, ErrBlockNotFound
		}

		receiptBytes, err := json.Marshal(receipt)
		if err != nil {
			logger.Error(ctx, err.Error())
			return nil, err
		}
		var receiptJSON datatypes.JSON
		if err := receiptJSON.UnmarshalJSON(receiptBytes); err != nil {
			logger.Error(ctx, err.Error())
			return nil, err
		}

		tx := entity.Tx{
			PoolID:      poolID,
			TxHash:      txHash,
			BlockNumber: header.Number.Uint64(),
			BlockTime:   header.Time,
			Gas:         receipt.GasUsed,
			GasPrice:    decimal.NewFromBigInt(receipt.EffectiveGasPrice, 0),
			Receipt:     receiptJSON,
			IsFinalized: false,
		}

		txs = append(txs, &tx)
	}

	return txs, nil
}

func initSwapEvents(ctx context.Context, poolID uint64, token0Decimals, token1Decimals uint, logs []types.Log) ([]*entity.SwapEvent, error) {
	var swapEvents []*entity.SwapEvent

	for _, log := range logs {
		if len(log.Topics) == 0 || log.Topics[0] != uniswapv3.EventSwap {
			continue
		}

		event, err := uniswapv3.DecodeSwap(log)
		if err != nil {
			logger.WithFields(ctx, logger.Fields{
				"address":  log.Address.Hex(),
				"tx":       log.TxHash.Hex(),
				"logIndex": log.Index,
				"topic0":   log.Topics[0].Hex(),
			}).Warn("cannot decode swap event")
			continue
		}

		amount0 := decimal.NewFromBigInt(event.Amount0, 0)
		amount1 := decimal.NewFromBigInt(event.Amount1, 0)

		amount0Dec := amount0.Abs().Div(decimal.NewFromBigInt(big.NewInt(10), int32(token0Decimals)))
		amount1Dec := amount1.Abs().Div(decimal.NewFromBigInt(big.NewInt(10), int32(token1Decimals)))
		price := amount0Dec.Div(amount1Dec).Round(6)

		swapEvent := entity.SwapEvent{
			PoolID:  poolID,
			TxHash:  strings.ToLower(log.TxHash.Hex()),
			Amount0: amount0.String(),
			Amount1: amount1.String(),
			Price:   price.String(),
		}
		swapEvents = append(swapEvents, &swapEvent)
	}

	return swapEvents, nil
}

func uniqueBlockHashes(logs []types.Log) []common.Hash {
	m := map[string]struct{}{}
	blockHashes := []common.Hash{}
	for _, l := range logs {
		h := strings.ToLower(l.BlockHash.Hex())
		if _, ok := m[h]; ok {
			continue
		}
		m[h] = struct{}{}
		blockHashes = append(blockHashes, l.BlockHash)
	}

	return blockHashes
}

func uniqueTxHashes(logs []types.Log) []common.Hash {
	m := map[string]struct{}{}
	txHashes := []common.Hash{}
	for _, l := range logs {
		h := strings.ToLower(l.TxHash.Hex())
		if _, ok := m[h]; ok {
			continue
		}
		m[h] = struct{}{}
		txHashes = append(txHashes, l.TxHash)
	}

	return txHashes
}
