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
	txrepo "github.com/vuquang23/poseidon/internal/pkg/repository/tx"
	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
	asynqpkg "github.com/vuquang23/poseidon/pkg/asynq"
	"github.com/vuquang23/poseidon/pkg/eth"
	"github.com/vuquang23/poseidon/pkg/logger"
	"github.com/vuquang23/poseidon/pkg/uniswapv3"
)

type TaskService struct {
	config      Config
	poolRepo    poolrepo.IPoolRepository
	txRepo      txrepo.ITxRepository
	ethClient   eth.IClient
	asynqClient asynqpkg.IAsynqClient
}

func New(
	config Config,
	poolRepo poolrepo.IPoolRepository,
	txRepo txrepo.ITxRepository,
	ethClient eth.IClient,
	asynqClient asynqpkg.IAsynqClient,
) *TaskService {
	return &TaskService{
		config:      config,
		poolRepo:    poolRepo,
		txRepo:      txRepo,
		ethClient:   ethClient,
		asynqClient: asynqClient,
	}
}

func (s *TaskService) GetPeriodicTaskConfigs(ctx context.Context) ([]*asynq.PeriodicTaskConfig, error) {
	pools, err := s.poolRepo.GetPools(ctx)
	if err != nil {
		return nil, err
	}

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

	logs, _, toBlock, err := s.getLogs(ctx, cursor.BlockNumber, poolAddress)
	if err != nil {
		return err
	}

	blockHashes := uniqueBlockHashes(logs)
	blocks, err := s.getBlocks(ctx, blockHashes)
	if err != nil {
		return err
	}

	if err := s.enqueueTaskGetETHUSDTKlines(ctx, blocks); err != nil {
		return err
	}

	txHashes := uniqueTxHashes(logs)
	txReceipts, err := s.getTxs(ctx, txHashes)
	if err != nil {
		return err
	}

	txs, err := initTxs(ctx, poolID, blocks, txReceipts)
	if err != nil {
		return err
	}

	swapEvents, err := initSwapEvents(ctx, poolID, token0Decimals, token1Decimals, logs)
	if err != nil {
		return err
	}

	return s.txRepo.UpdateDataScanner(ctx, cursor.ID, toBlock+1, txs, swapEvents)
}

func (s *TaskService) getLogs(ctx context.Context, cursorBlockNbr uint64, poolAddress string) ([]types.Log, uint64, uint64, error) {
	latestBlock, err := s.ethClient.GetLatestBlock(ctx)
	if err != nil {
		return nil, 0, 0, err
	}
	latestBlockNbr := latestBlock.Number().Uint64()

	fromBlock := cursorBlockNbr
	toBlock := fromBlock + s.config.BlockBatchSize
	if toBlock > latestBlockNbr {
		toBlock = latestBlockNbr
	}

	if fromBlock > toBlock {
		logger.WithFields(ctx, logger.Fields{
			"poolAddress":    poolAddress,
			"fromBlock":      fromBlock,
			"toBlock":        toBlock,
			"batchSize":      s.config.BlockBatchSize,
			"latestBlockNbr": latestBlock,
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
		receipts  = make([]*types.Receipt, 0, len(txHashes))
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

	for i := 0; i < len(txHashes); i++ {
		r, _ := resultMap.Load(i)
		receipt := r.(*types.Receipt)
		receipts[i] = receipt
	}

	return receipts, nil
}

func (s *TaskService) getBlocks(ctx context.Context, blockHashes []common.Hash) ([]*types.Block, error) {
	var (
		wg        errgroup.Group
		resultMap sync.Map
		blocks    = make([]*types.Block, 0, len(blockHashes))
	)

	for idx, blockHash := range blockHashes {
		_idx, _blockHash := idx, blockHash

		wg.Go(func() error {
			receipt, err := s.ethClient.GetBlockByHash(ctx, _blockHash)
			if err != nil {
				return err
			}
			resultMap.Store(_idx, receipt)
			return nil
		})
	}

	for i := 0; i < len(blockHashes); i++ {
		r, _ := resultMap.Load(i)
		block := r.(*types.Block)
		blocks[i] = block
	}

	return blocks, nil
}

func (s *TaskService) enqueueTaskGetETHUSDTKlines(ctx context.Context, blocks []*types.Block) error {
	exists := map[uint64]struct{}{}
	for _, b := range blocks {
		if _, ok := exists[b.Time()]; ok {
			continue
		}

		if err := s.asynqClient.EnqueueTask(
			ctx,
			valueobject.TaskTypeGetETHUSDTKline,
			"", "",
			valueobject.TaskGetETHUSDTKlinePayload{Time: b.Time()},
			-1,
		); err != nil {
			return err
		}

		exists[b.Time()] = struct{}{}
	}

	return nil
}

func (s *TaskService) GetLatestFinalizedBlockNumber(ctx context.Context) (uint64, error) {
	block, err := s.ethClient.GetLatestBlock(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		return 0, err
	}

	t := time.Unix(int64(block.Time()), 0)
	if t.Before(time.Now().Add(-s.config.BlockTimeDelayThreshold)) {
		logger.WithFields(ctx, logger.Fields{
			"blockTime":      block.Time(),
			"delayThreshold": s.config.BlockTimeDelayThreshold,
		}).Error(ErrInvalidLatestBlock.Error())
		return 0, ErrInvalidLatestBlock
	}

	return block.NumberU64() - s.config.BlockFinalization, nil
}

func initTxs(ctx context.Context, poolID uint64, blocks []*types.Block, txReceipts []*types.Receipt) ([]*entity.Tx, error) {
	blockByBlockHash := map[string]*types.Block{}
	for _, b := range blocks {
		h := strings.ToLower(b.Hash().Hex())
		blockByBlockHash[h] = b
	}

	txs := make([]*entity.Tx, 0, len(txReceipts))

	for _, receipt := range txReceipts {
		txHash := strings.ToLower(receipt.TxHash.Hex())

		block, ok := blockByBlockHash[strings.ToLower(receipt.BlockHash.Hex())]
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
			BlockNumber: block.NumberU64(),
			BlockTime:   block.Time(),
			Gas:         receipt.CumulativeGasUsed,
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
		if log.Topics[0] != uniswapv3.EventSwap {
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
		price := amount0Dec.Div(amount1Dec)

		swapEvent := entity.SwapEvent{
			PoolID:  poolID,
			TxHash:  strings.ToLower(log.TxHash.Hex()),
			Amount0: amount0,
			Amount1: amount1,
			Price:   price,
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
		blockHashes = append(blockHashes, l.TxHash)
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
