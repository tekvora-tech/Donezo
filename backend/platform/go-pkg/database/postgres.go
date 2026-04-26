package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresConfig konfigurasi koneksi PostgreSQL
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string // disable, require, verify-ca, verify-full
	// Pool settings
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// NewPostgres buat koneksi PostgreSQL dengan pgxpool
func NewPostgres(cfg PostgresConfig) (*pgxpool.Pool, error) {
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	// Pool settings
	if cfg.MaxConns > 0 {
		poolConfig.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 {
		poolConfig.MinConns = cfg.MinConns
	}
	if cfg.MaxConnLifetime > 0 {
		poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	}
	if cfg.MaxConnIdleTime > 0 {
		poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	// Test koneksi
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}