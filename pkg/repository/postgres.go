package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	groupsTable      = "groupss"
	songsTable       = "songs"
	songDetailsTable = "songdetails"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	logrus.Debug("Connecting to PostgreSQL database")

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)

	logrus.WithFields(logrus.Fields{
		"host":     cfg.Host,
		"port":     cfg.Port,
		"username": cfg.Username,
		"dbname":   cfg.DBName,
	}).Info("Connecting to database")

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		logrus.WithError(err).Error("Failed to open database connection")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logrus.WithError(err).Error("Failed to ping database")
		return nil, err
	}

	logrus.Info("Database connection established successfully")
	return db, nil
}
