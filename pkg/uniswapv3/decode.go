package uniswapv3

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	EventSwap     common.Hash
	univ3Filterer *PoolFilterer
)

func init() {
	filterer, err := NewPoolFilterer(common.Address{}, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	univ3Filterer = filterer

	EventSwap = UniswapV3PoolABI.Events["Swap"].ID
}

func DecodeSwap(event types.Log) (*PoolSwap, error) {
	return univ3Filterer.ParseSwap(event)
}
