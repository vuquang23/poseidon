package uniswapv3

import (
	"bytes"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	UniswapV3PoolABI abi.ABI
)

func init() {
	builder := []struct {
		ABI  *abi.ABI
		data []byte
	}{
		{&UniswapV3PoolABI, uniswapV3PoolJson},
	}

	for _, b := range builder {
		var err error
		*b.ABI, err = abi.JSON(bytes.NewReader(b.data))
		if err != nil {
			panic(err)
		}
	}
}
