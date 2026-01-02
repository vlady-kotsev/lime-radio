package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

//go:embed schema/schema.sql
var schemaSQL string

type Storage struct {
	DB *sqlx.DB
}

var _ Storager = (*Storage)(nil)

func NewStorage(lc fx.Lifecycle, logger *zap.Logger) (*Storage, error) {
	path := "internal/repository/data/data.db" // or get from config

	db, err := openDB(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &Storage{DB: db}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Creating database schema")
			_, err = storage.DB.Exec(schemaSQL)
			return err
		},
		OnStop: func(context.Context) error {
			logger.Info("Closing database connection")
			return db.Close()
		},
	})

	return storage, nil
}

func (s *Storage) Beginx() (*sqlx.Tx, error) {
	return s.DB.Beginx()
}

func (s *Storage) Select(dest any, query string, args ...any) error {
	return s.DB.Select(dest, query, args...)
}

func (s *Storage) Get(dest any, query string, args ...any) error {
	return s.DB.Get(dest, query, args...)
}

func openDB(path string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	return db, nil
}
