package repository

import (
	"database/sql"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(databaseURL string) (*gorm.DB, error) {
	gdb, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}
	configurePool(sqlDB)
	return gdb, nil
}

func configurePool(sqlDB *sql.DB) {
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
}

func DatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}
