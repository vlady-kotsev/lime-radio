package repository

import "github.com/jmoiron/sqlx"

type Storager interface {
	Beginx() (*sqlx.Tx, error)
	Select(dest any, query string, args ...any) error
}
