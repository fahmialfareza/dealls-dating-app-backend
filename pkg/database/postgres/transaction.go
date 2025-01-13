package postgres

import (
	"context"
	"database/sql"

	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
)

func StartTransaction(ctx context.Context, opts *sql.TxOptions) (context.Context, error) {
	segment := logger.StartSegment(ctx, "postgres.StartTransaction")
	defer segment.End()

	ctx, _, err := PostgresClient.Begin(ctx, opts)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"opts": opts,
		})
		return ctx, err
	}

	return ctx, nil
}

func Commit(ctx context.Context) error {
	segment := logger.StartSegment(ctx, "postgres.Commit")
	defer segment.End()

	if err := PostgresClient.Commit(ctx); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		return err
	}

	return nil
}

func Rollback(ctx context.Context) error {
	segment := logger.StartSegment(ctx, "postgres.Rollback")
	defer segment.End()

	if err := PostgresClient.Rollback(ctx); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		return err
	}

	return nil
}
