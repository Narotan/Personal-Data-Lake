package db

import (
	"context"

	activitywatch_db "DataLake/internal/db/activitywatch"
	googlecalendar_db "DataLake/internal/db/googlecalendar"
	googlefit_db "DataLake/internal/db/googlefit"
	wakatime_db "DataLake/internal/db/wakatime"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	WakaTime       *wakatime_db.Queries
	GoogleFit      *googlefit_db.Queries
	ActivityWatch  *activitywatch_db.Queries
	GoogleCalendar *googlecalendar_db.Queries
	db             *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		WakaTime:       wakatime_db.New(pool),
		ActivityWatch:  activitywatch_db.New(pool),
		GoogleFit:      googlefit_db.New(pool),
		GoogleCalendar: googlecalendar_db.New(pool),
		db:             pool,
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

func (s *Store) ExecTxGoogleFit(ctx context.Context, fn func(*googlefit_db.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := googlefit_db.New(tx)

	err = fn(q)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) ExecTxGoogleCalendar(ctx context.Context, fn func(*googlecalendar_db.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := googlecalendar_db.New(tx)

	err = fn(q)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
