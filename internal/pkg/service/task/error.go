package task

import "errors"

var (
	ErrBlockNotFound      = errors.New("block not found")
	ErrInvalidBlockRange  = errors.New("invalid block range")
	ErrInvalidLatestBlock = errors.New("invalid latest block")
)
