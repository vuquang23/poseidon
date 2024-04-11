package chain

import (
	"context"
	"time"

	"github.com/vuquang23/poseidon/pkg/eth"
	"github.com/vuquang23/poseidon/pkg/logger"
)

type ChainRepository struct {
	config    Config
	ethClient eth.IClient
}

func New(config Config, ethClient eth.IClient) *ChainRepository {
	return &ChainRepository{
		config:    config,
		ethClient: ethClient,
	}
}

func (r *ChainRepository) GetLatestFinalizedBlockNumber(ctx context.Context) (uint64, error) {
	block, err := r.ethClient.GetLatestBlock(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		return 0, err
	}

	t := time.Unix(int64(block.Time()), 0)
	if t.Before(time.Now().Add(-r.config.BlockTimeDelayThreshold)) {
		logger.WithFields(ctx, logger.Fields{
			"blockTime":      block.Time(),
			"delayThreshold": r.config.BlockTimeDelayThreshold,
		}).Error(ErrInvalidLatestBlock.Error())
		return 0, ErrInvalidLatestBlock
	}

	return block.NumberU64() - r.config.BlockFinalization, nil
}
