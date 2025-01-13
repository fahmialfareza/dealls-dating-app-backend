package postgres

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type IPostgres interface {
	Begin(ctx context.Context, opts *sql.TxOptions) (context.Context, *gorm.DB, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	GORM(ctx context.Context) *gorm.DB
}

type Postgres struct {
	db *gorm.DB
}

func NewPostgres(setLimits bool) IPostgres {
	db, err := openPostgres(setLimits)
	if err != nil {
		panic(err)
	}

	return &Postgres{
		db: db,
	}
}
