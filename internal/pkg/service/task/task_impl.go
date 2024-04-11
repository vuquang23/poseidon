package task

import (
	"context"
	"encoding/json"
	"errors"

	"gorm.io/datatypes"

	"github.com/hibiken/asynq"
	"github.com/vuquang23/poseidon/internal/pkg/entity"
	chainrepo "github.com/vuquang23/poseidon/internal/pkg/repository/chain"
	poolrepo "github.com/vuquang23/poseidon/internal/pkg/repository/pool"
	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
	"github.com/vuquang23/poseidon/pkg/logger"
)

type TaskService struct {
	poolRepo  poolrepo.IPoolRepository
	chainRepo chainrepo.IChainRepository
}

func New(poolRepo poolrepo.IPoolRepository, chainRepo chainrepo.IChainRepository) *TaskService {
	return &TaskService{
		poolRepo:  poolRepo,
		chainRepo: chainRepo,
	}
}

func (s *TaskService) GetPeriodicTaskConfigs(ctx context.Context) ([]*asynq.PeriodicTaskConfig, error) {
	return nil, errors.New("not implemented")
}

func (s *TaskService) HandlePoolCreated(ctx context.Context, poolAddress string) error {
	pool, err := s.poolRepo.GetPoolByAddress(ctx, poolAddress)
	if err != nil {
		return err
	}

	// scanner cursor
	scannerCursor := entity.BlockCursor{
		PoolID:      pool.ID,
		Type:        valueobject.BlockCursorTypeScanner,
		BlockNumber: pool.StartBlock,
	}

	// finalizer cursor
	latestFinalizedBlockNbr, err := s.chainRepo.GetLatestFinalizedBlockNumber(ctx)
	if err != nil {
		return err
	}
	extra := valueobject.BlockCursorFinalizerExtra{
		CreatedAtFinalizedBlock: latestFinalizedBlockNbr,
	}
	extraBytes, err := json.Marshal(extra)
	if err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	var extraJSON datatypes.JSON
	if err := extraJSON.UnmarshalJSON(extraBytes); err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	finalizerCursor := entity.BlockCursor{
		PoolID:      pool.ID,
		Type:        valueobject.BlockCursorTypeFinalizer,
		BlockNumber: pool.StartBlock,
		Extra:       extraJSON,
	}

	// persist cursors to DB
	return s.poolRepo.CreateBlockCursors(ctx, []*entity.BlockCursor{&scannerCursor, &finalizerCursor})
}
