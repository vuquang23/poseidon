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
	poolsvc "github.com/vuquang23/poseidon/internal/pkg/service/pool"
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

	poolRepo *poolrepo.PoolRepository
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

	// service
	poolSvc := poolsvc.New(poolRepo)

	// server
	ginEngine = gin.New()

	api.RegisterRoutes(&conf, ginEngine, poolSvc)
}

func (suite *TestSuite) SetupTest() {
	// http
	httpmock.Reset()

	// postgres
	err := db.Exec("TRUNCATE TABLE pools CASCADE").Error
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
