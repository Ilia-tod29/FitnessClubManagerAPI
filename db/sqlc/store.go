package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// SQLStore provides all functions to execute SQL queries //and transactions
type SQLStore struct {
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		Queries: New(connPool),
	}
}
