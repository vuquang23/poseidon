package task

import "time"

type Config struct {
	BlockBatchSize          uint64
	BlockTimeDelayThreshold time.Duration
	BlockFinalization       uint64 // Number of blocks to reach finalization
}
