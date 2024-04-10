package postgres

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

const schemaMigrationTable = "schema_migrations"

func MigrateUp(db *gorm.DB, dir string, up int) error {
	m, err := newMigrate(db, dir)
	if err != nil {
		return err
	}
	if up == 0 {
		return m.Up()
	}
	return m.Steps(up)
}

func MigrateDown(db *gorm.DB, dir string, down int) error {
	m, err := newMigrate(db, dir)
	if err != nil {
		return err
	}
	if down == 0 {
		return m.Down()
	}
	return m.Steps(down)
}

func newMigrate(db *gorm.DB, dir string) (*migrate.Migrate, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{
		MigrationsTable: schemaMigrationTable,
	})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		dir,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}
