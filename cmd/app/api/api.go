package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"

	"github.com/vuquang23/poseidon/internal/pkg/config"
	"github.com/vuquang23/poseidon/internal/pkg/server"
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

	// auto migration
	if err := postgres.MigrateUp(db, "file://./migration/postgres", 0); err != nil {
		return err
	}

	// server
	server := server.GinEngine(conf.Http, conf.Log, logger.LoggerBackendZap)
	router := server.Group("/api/v1")

	// api

	// setup routes

	/// health
	router.GET("health/live", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, "OK") })
	router.GET("health/ready", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, "OK") })

	return server.Run(conf.Http.BindAddress)
}
