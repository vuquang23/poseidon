package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/vuquang23/poseidon/pkg/logger"
)

type IClient interface {
	GetLatestBlock(ctx context.Context) (*types.Block, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	GetLogs(ctx context.Context, fromBlock, toBlock uint64, addresses []common.Address) ([]types.Log, error)
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

func (c *Client) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	blockHeader, err := c.ethClient.HeaderByHash(ctx, hash)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, err
	}
	return blockHeader, err
}

func (c *Client) GetLogs(ctx context.Context, fromBlock, toBlock uint64, addresses []common.Address) ([]types.Log, error) {
	q := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: addresses,
	}

	logs, err := c.ethClient.FilterLogs(ctx, q)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, err
	}

	return logs, nil
}

func (c *Client) GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := c.ethClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, err
	}

	return receipt, nil
}
