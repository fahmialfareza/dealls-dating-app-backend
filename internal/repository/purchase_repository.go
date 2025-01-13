package repository

import (
	"context"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
)

type IPurchaseRepository interface {
	CreatePurchase(ctx context.Context, data *domain.Purchase) error
}

// CreatePurchase implements IRepository.
func (r *Repository) CreatePurchase(ctx context.Context, data *domain.Purchase) error {
	segment := logger.StartSegment(ctx, "Repository.CreatePurchase")
	defer segment.End()

	// get from postgres
	if err := r.postgres.InsertPurchase(ctx, data); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"date": data,
		})
		return err
	}

	return nil
}
