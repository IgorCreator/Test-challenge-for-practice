package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool   *pgxpool.Pool
	crypto Crypto
}

type Crypto interface {
	Encrypt(plaintext string) ([]byte, []byte, error)
	Decrypt(ciphertext, nonce []byte) (string, error)
}

func NewStore(ctx context.Context, dsn string, crypto Crypto) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}
	cfg.MaxConns = 10
	cfg.MaxConnLifetime = 5 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return &Store{pool: pool, crypto: crypto}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}
