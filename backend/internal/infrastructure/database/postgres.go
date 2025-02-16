package database

import (
	"context"
	"time"

	"github.com/Sparker0i/cactro-polls/internal/infrastructure/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
	cfg  *config.DatabaseConfig
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.ConnectionString())
	if err != nil {
		return nil, err
	}

	// Configure pool
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	// Create pool
	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return &Database{
		pool: pool,
		cfg:  cfg,
	}, nil
}

func (db *Database) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *Database) Stats() *pgxpool.Stat {
	return db.pool.Stat()
}

// Transaction management
type Tx struct {
	pgx.Tx
}

func (db *Database) BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

func (tx *Tx) Commit(ctx context.Context) error {
	return tx.Tx.Commit(ctx)
}

func (tx *Tx) Rollback(ctx context.Context) error {
	return tx.Tx.Rollback(ctx)
}
