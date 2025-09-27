package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/esuEdu/AWS-Go/internal/config"
)

func NewPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS items (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE OR REPLACE FUNCTION set_updated_at()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = NOW();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		DO $$ BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_trigger WHERE tgname = 'items_set_updated_at'
		) THEN
			CREATE TRIGGER items_set_updated_at
			BEFORE UPDATE ON items
			FOR EACH ROW EXECUTE PROCEDURE set_updated_at();
		END IF;
		END $$;
	`)
	return err
}
