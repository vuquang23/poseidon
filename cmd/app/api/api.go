package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli/v2"

	"github.com/vuquang23/poseidon/internal/pkg/api"
	"github.com/vuquang23/poseidon/internal/pkg/config"
	poolrepo "github.com/vuquang23/poseidon/internal/pkg/repository/pool"
	"github.com/vuquang23/poseidon/internal/pkg/server"
	"github.com/vuquang23/poseidon/internal/pkg/server/middleware"
	poolsvc "github.com/vuquang23/poseidon/internal/pkg/service/pool"
	"github.com/vuquang23/poseidon/pkg/asynq"
	"github.com/vuquang23/poseidon/pkg/logger"
	"github.com/vuquang23/poseidon/pkg/postgres"
)

func RunAPI(c *cli.Context) error {
	conf := config.New()
	if err := conf.Load(c.String("config")); err != nil {
		return err
	}

	// logger
	_, err := logger.Init(conf.Log, logger.LoggerBackendZap)
	if err != nil {
		return err
	}

	// postgres
	db, err := postgres.New(conf.Postgres)
	if err != nil {
		return err
	}

	// asynq client
	asynqClient, err := asynq.NewClient(conf.Redis)
	if err != nil {
		return err
	}

	// auto migration
	if err := postgres.MigrateUp(db, "file://./migration/postgres", 0); err != nil && err != migrate.ErrNoChange {
		return err
	}

	// repository
	poolRepo := poolrepo.New(db, asynqClient)

	// service
	poolSvc := poolsvc.New(poolRepo)

	// server
	server := server.GinEngine(conf.Http, conf.Log, logger.LoggerBackendZap)
	router := server.Group("/api/v1")

	// api

	/// health
	router.GET("health/live", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, "OK") })
	router.GET("health/ready", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, "OK") })

	/// pool
	router.POST("pools", middleware.NewAuthMiddleware(conf.Common.APIKey), api.CreatePool(poolSvc))

	return server.Run(conf.Http.BindAddress)
}
