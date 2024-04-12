package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/vuquang23/poseidon/internal/pkg/api"
	"github.com/vuquang23/poseidon/internal/pkg/config"
	poolrepo "github.com/vuquang23/poseidon/internal/pkg/repository/pool"
	pricerepo "github.com/vuquang23/poseidon/internal/pkg/repository/price"
	txrepo "github.com/vuquang23/poseidon/internal/pkg/repository/tx"
	poolsvc "github.com/vuquang23/poseidon/internal/pkg/service/pool"
	tasksvc "github.com/vuquang23/poseidon/internal/pkg/service/task"
	txsvc "github.com/vuquang23/poseidon/internal/pkg/service/tx"
	"github.com/vuquang23/poseidon/pkg/asynq"
	"github.com/vuquang23/poseidon/pkg/logger"
	"github.com/vuquang23/poseidon/pkg/postgres"
)

var (
	migrationDir = "file://../../../migration/postgres"
	configFile   = "../config/test.yaml"

	XAPIKey   string
	ginEngine *gin.Engine
	db        *gorm.DB
	taskSvc   *tasksvc.TaskService
	poolRepo  *poolrepo.PoolRepository
)

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) SetupSuite() {
	conf := config.New()
	if err := conf.Load(configFile); err != nil {
		suite.FailNow(err.Error())
	}

	XAPIKey = conf.Common.APIKey

	// logger
	_, err := logger.Init(conf.Log, logger.LoggerBackendZap)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// postgres
	db, err = postgres.New(conf.Postgres)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// asynq client
	asynqClient, err := asynq.NewClient(conf.Redis)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// auto migration
	if err := postgres.MigrateUp(db, migrationDir, 0); err != nil && err != migrate.ErrNoChange {
		suite.FailNow(err.Error())
	}

	// repository
	poolRepo = poolrepo.New(db, asynqClient)
	txRepo := txrepo.New(db)
	priceRepo := pricerepo.New(db)

	// service
	poolSvc := poolsvc.New(poolRepo)
	taskSvc = tasksvc.New(conf.Service.Task, poolRepo, txRepo, priceRepo, nil, nil, nil)
	txSvc := txsvc.New(txRepo, priceRepo)

	// server
	ginEngine = gin.New()

	api.RegisterRoutes(&conf, ginEngine, poolSvc, txSvc)
}

func (suite *TestSuite) SetupTest() {
	// http
	httpmock.Reset()

	// postgres
	err := db.Exec("TRUNCATE TABLE pools CASCADE").Error
	if err != nil {
		suite.FailNow(err.Error())
	}

	err = db.Exec("TRUNCATE TABLE txs CASCADE").Error
	if err != nil {
		suite.FailNow(err.Error())
	}

	err = db.Exec("TRUNCATE TABLE swap_events CASCADE").Error
	if err != nil {
		suite.FailNow(err.Error())
	}

	err = db.Exec("TRUNCATE TABLE block_cursors CASCADE").Error
	if err != nil {
		suite.FailNow(err.Error())
	}

	err = db.Exec("TRUNCATE TABLE ethusdt_klines CASCADE").Error
	if err != nil {
		suite.FailNow(err.Error())
	}

	// redis
}

func (suite *TestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
