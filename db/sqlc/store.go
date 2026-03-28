package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	ExecuteTx(ctx context.Context, run func(*Queries) error) error
}

type SQLStore struct {
	*Queries
	DB *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		DB:      db,
		Queries: New(db),
	}
}

func (s *SQLStore) ExecuteTx(ctx context.Context, run func(*Queries) error) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	if err = run(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
