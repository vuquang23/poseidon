package chain

import "context"

type IChainRepository interface {
	GetLatestFinalizedBlockNumber(ctx context.Context) (uint64, error)
}
