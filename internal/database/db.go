package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/coddemn/base-structure-go/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func New(ctx context.Context, cfg config.DBConfig) (*sql.DB, error) {
	connStr := cfg.DSN()

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parse pgxpool config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	// test connect
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	sqlDB := stdlib.OpenDBFromPool(pool)
	sqlDB.SetMaxOpenConns(0)
	sqlDB.SetMaxIdleConns(0)

	return sqlDB, nil
}
