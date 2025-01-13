package repository

import (
	"context"
	"time"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
)

type ISwiperRepository interface {
	CreateSwipe(ctx context.Context, data *domain.Swipe) error
	GetSwipe(ctx context.Context, swiperID uint, date time.Time) ([]domain.Swipe, error)
}

// CreateSwipe implements IRepository.
func (r *Repository) CreateSwipe(ctx context.Context, data *domain.Swipe) error {
	segment := logger.StartSegment(ctx, "Repository.CreateSwipe")
	defer segment.End()

	// get from postgres
	if err := r.postgres.InsertSwipe(ctx, data); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"date": data,
		})
		return err
	}

	return nil
}

// GetSwipe implements IRepository.
func (r *Repository) GetSwipe(ctx context.Context, swiperID uint, date time.Time) ([]domain.Swipe, error) {
	segment := logger.StartSegment(ctx, "Repository.GetSwipe")
	defer segment.End()

	// get from postgres
	data, err := r.postgres.GetSwipe(ctx, swiperID, date)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"swiper_id": swiperID,
			"date":      data,
		})
		return data, err
	}

	return data, nil
}
