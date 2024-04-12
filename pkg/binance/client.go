package binance

import (
	"context"

	"github.com/adshao/go-binance/v2"

	"github.com/vuquang23/poseidon/pkg/logger"
)

//go:generate mockgen -destination=client_mock.go -package binance . IClient
type IClient interface {
	GetKlines(ctx context.Context, startTime, endTime int64, limit int, symbol, interval string) ([]*binance.Kline, error)
}

type Client struct {
	spot *binance.Client
}

func NewClient() *Client {
	spot := binance.NewClient("", "")

	return &Client{
		spot: spot,
	}
}

func (c *Client) GetKlines(ctx context.Context, startTime, endTime int64, limit int, symbol, interval string) ([]*binance.Kline, error) {
	req := c.spot.NewKlinesService().StartTime(startTime).Symbol(symbol).Interval(interval)

	if endTime > 0 {
		req.EndTime(endTime)
	}

	if limit > 0 {
		req.Limit(int(limit))
	}

	klines, err := req.Do(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		return nil, err
	}

	return klines, nil
}
