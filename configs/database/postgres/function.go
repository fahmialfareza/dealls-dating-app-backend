package postgres

import (
	"context"
	"database/sql"

	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
	"gorm.io/gorm"
)

func (m *Postgres) Begin(ctx context.Context, opts *sql.TxOptions) (context.Context, *gorm.DB, error) {
	segment := logger.StartSegment(ctx, "Postgres.Begin")
	defer segment.End()

	result := m.db.WithContext(ctx).Begin(opts)
	if result.Error != nil {
		logger.PrintErrorLog(ctx, result.Error, logger.GetErrorFileLine(), map[string]interface{}{
			"opts": opts,
		})
		return ctx, result, result.Error
	}

	ctx = context.WithValue(ctx, postgresTxContext, result)

	return ctx, result, nil
}

// Commit implements IPostgres.
func (m *Postgres) Commit(ctx context.Context) error {
	segment := logger.StartSegment(ctx, "Postgres.Commit")
	defer segment.End()

	// get transaction
	tx, isTx := getTransaction(ctx)
	if isTx {
		result := tx.Commit()
		if result.Error != nil {
			logger.PrintErrorLog(ctx, result.Error, logger.GetErrorFileLine(), nil)
			return result.Error
		}
		return nil
	}

	return nil
}

// Rollback implements IPostgres.
func (m *Postgres) Rollback(ctx context.Context) error {
	segment := logger.StartSegment(ctx, "Postgres.Rollback")
	defer segment.End()

	// get transaction
	tx, isTx := getTransaction(ctx)
	if isTx {
		result := tx.WithContext(ctx).Rollback()
		if result.Error != nil {
			logger.PrintErrorLog(ctx, result.Error, logger.GetErrorFileLine(), nil)
			return result.Error
		}
	}

	return nil
}

// GetGORMDB implements IPostgres.
func (m *Postgres) GORM(ctx context.Context) *gorm.DB {
	segment := logger.StartSegment(ctx, "Postgres.GORM")
	defer segment.End()

	// get transaction
	result, isTx := getTransaction(ctx)
	if isTx {
		return result.WithContext(ctx)
	}

	return m.db.WithContext(ctx)
}

func getTransaction(ctx context.Context) (*gorm.DB, bool) {
	segment := logger.StartSegment(ctx, "mysql.GetTransaction")
	defer segment.End()

	// Retrieve the transaction from context
	txCtx := ctx.Value(postgresTxContext)
	tx, ok := txCtx.(*gorm.DB) // GORM transaction is of type *gorm.DB
	if !ok {
		return tx, false
	}

	return tx, true
}
