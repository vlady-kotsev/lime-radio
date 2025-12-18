package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

//go:embed sql/schema.sql
var schemaSQL string

type Storage struct {
	db *sql.DB
}

func NewStorage(lc fx.Lifecycle, logger *zap.Logger) (*Storage, error) {
	path := "internal/repository/data/data.db" // or get from config

	db, err := openDB(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &Storage{db: db}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Creating database schema")
			_, err = storage.db.Exec(schemaSQL)
			return err
		},
		OnStop: func(context.Context) error {
			logger.Info("Closing database connection")
			return db.Close()
		},
	})

	return storage, nil
}

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// SQLite works best with a single writer
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	return db, nil
}
