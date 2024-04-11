package eth

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/vuquang23/poseidon/pkg/logger"
)

type IClient interface {
	GetLatestBlock(ctx context.Context) (*types.Block, error)
}

type Client struct {
	config    Config
	ethClient *ethclient.Client
}

func NewClient(config Config) (*Client, error) {
	ethClient, err := ethclient.Dial(config.RPC)
	if err != nil {
		return nil, err
	}

	return &Client{
		config:    config,
		ethClient: ethClient,
	}, nil
}

func (c *Client) GetLatestBlock(ctx context.Context) (*types.Block, error) {
	block, err := c.ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, err
	}

	return block, nil
}
