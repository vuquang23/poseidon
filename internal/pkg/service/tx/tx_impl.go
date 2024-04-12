package tx

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
	pricerepo "github.com/vuquang23/poseidon/internal/pkg/repository/price"
	txrepo "github.com/vuquang23/poseidon/internal/pkg/repository/tx"
	"github.com/vuquang23/poseidon/internal/pkg/service/dto"
	timepkg "github.com/vuquang23/poseidon/internal/pkg/util/time"
)

type TxService struct {
	txRepo    txrepo.ITxRepository
	priceRepo pricerepo.IPriceRepository
}

func New(txRepo txrepo.ITxRepository, priceRepo pricerepo.IPriceRepository) *TxService {
	return &TxService{
		txRepo:    txRepo,
		priceRepo: priceRepo,
	}
}

func (s *TxService) GetTxFeeUSDT(ctx context.Context, txHash string) (decimal.Decimal, error) {
	tx, err := s.txRepo.GetTxByHash(ctx, txHash)
	if err != nil {
		return decimal.Decimal{}, err
	}

	blockTime := tx.BlockTime
	openTimeNSec := 1000 * timepkg.RoundDown(int64(blockTime), time.Minute)
	kline, err := s.priceRepo.GetKlineByOpenTime(ctx, openTimeNSec)
	if err != nil {
		return decimal.Decimal{}, err
	}

	usdtValue := decimal.NewFromInt(int64(tx.Gas)).
		Mul(tx.GasPrice).
		Mul(kline.OHLC4).
		Div(decimal.New(1, 18)).
		Round(6)

	return usdtValue, nil
}

func (s *TxService) GetSwapEventsByTxHash(ctx context.Context, txHash string) ([]*entity.SwapEvent, error) {
	swapEvents, err := s.txRepo.GetSwapEventsByTxHash(ctx, txHash)
	if err != nil {
		return nil, err
	}

	return swapEvents, nil
}

func (s *TxService) GetTxs(ctx context.Context, query dto.GetTxsQuery) ([]*entity.Tx, int64, error) {
	limit := query.PageSize
	offset := (query.Page - 1) * query.PageSize

	return s.txRepo.GetTxs(ctx, query.PoolAddress, offset, limit)
}
