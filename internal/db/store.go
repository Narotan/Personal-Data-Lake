package db

import (
	"context"

	wakatime_db "DataLake/internal/db/wakatime"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	WakaTime *wakatime_db.Queries
	db       *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		WakaTime: wakatime_db.New(pool),
		db:       pool,
	}
}

func (s *Store) ExecTx(ctx context.Context, fn func(*wakatime_db.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := wakatime_db.New(tx)

	err = fn(q)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
