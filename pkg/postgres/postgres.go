package postgres

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(cfg Config) (*gorm.DB, error) {
	gormDB, err := gorm.Open(
		postgres.Open(dsn(cfg)),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(cfg.LogLevel)),
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		},
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnLifeTime) * time.Second)

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return gormDB, nil
}

func dsn(c Config) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}
