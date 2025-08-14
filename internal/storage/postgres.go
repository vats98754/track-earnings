package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct { DB *sqlx.DB }

func NewPostgres(dsn string) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	return &Postgres{DB: db}, nil
}

func (p *Postgres) Close() error { return p.DB.Close() }

func (p *Postgres) Ping(ctx context.Context) error { return p.DB.PingContext(ctx) }
