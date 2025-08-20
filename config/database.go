package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func NewDatabase(config *Config, logger *logrus.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.Username,
		config.Database.Password,
		config.Database.DbName,
	)
	db, err := sql.Open(config.Database.DriverName, dsn)
	if err != nil {
		logger.WithError(err).Error("Failed to open database")
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.WithError(err).Error("Failed to ping database")
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
