package asynq

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hibiken/asynq"

	"github.com/vuquang23/poseidon/pkg/logger"
	"github.com/vuquang23/poseidon/pkg/redis"
)

var ErrEmptyRedisHost = errors.New("redis host is empty")

type IAsynqClient interface {
	EnqueueTask(
		ctx context.Context,
		taskType string,
		taskID string,
		queueID string,
		payload interface{},
		maxRetry int,
	) error
}

type AsynqClient struct {
	client *asynq.Client
}

func NewClient(config redis.Config) (*AsynqClient, error) {
	redisConnOpt, err := GetAsynqRedisConnectionOption(config)
	if err != nil {
		return nil, err
	}

	return &AsynqClient{client: asynq.NewClient(redisConnOpt)}, nil
}

func (c *AsynqClient) EnqueueTask(
	ctx context.Context,
	taskType string,
	taskID string,
	queueID string,
	payload interface{},
	maxRetry int,
) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	opts := []asynq.Option{}
	if taskID != "" {
		opts = append(opts, asynq.TaskID(taskID))
	}
	if queueID != "" {
		opts = append(opts, asynq.Queue(queueID))
	}
	if maxRetry >= 0 {
		opts = append(opts, asynq.MaxRetry(maxRetry))
	}

	task := asynq.NewTask(taskType, payloadBytes, opts...)
	_, err = c.client.Enqueue(task)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	return nil
}

func GetAsynqRedisConnectionOption(cfg redis.Config) (asynq.RedisConnOpt, error) {
	redisAddresses := cfg.Addresses
	if cfg.MasterName != "" {
		return asynq.RedisFailoverClientOpt{
			SentinelAddrs: redisAddresses,
			MasterName:    cfg.MasterName,
			Password:      cfg.Password,
		}, nil

	}

	if len(redisAddresses) == 0 {
		redisAddresses = append(redisAddresses, "")
	}
	if len(redisAddresses) == 1 {
		return asynq.RedisClientOpt{
			Addr:     redisAddresses[0],
			Password: cfg.Password,
		}, nil
	}

	return asynq.RedisClusterClientOpt{
		Addrs:    redisAddresses,
		Password: cfg.Password,
	}, nil
}
