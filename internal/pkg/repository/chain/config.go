package chain

import "time"

type Config struct {
	BlockTimeDelayThreshold time.Duration
	BlockFinalization       uint64 // Number of blocks to reach finalization
}
