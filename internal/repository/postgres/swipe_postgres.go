package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
	"gorm.io/gorm"
)

type ISwipePostgresRepository interface {
	InsertSwipe(ctx context.Context, data *domain.Swipe) error
	GetSwipe(ctx context.Context, swiperID uint, date time.Time) ([]domain.Swipe, error)
}

// InsertSwipe implements IPostgresRepository.
func (u *PostgresRepository) InsertSwipe(ctx context.Context, data *domain.Swipe) error {
	segment := logger.StartSegment(ctx, "PostgresRepository.InsertSwipe")
	defer segment.End()

	if err := hystrix.DoC(ctx, constant.HystrixPostgres, func(ctx context.Context) error {
		// create
		if err := u.db.GORM(ctx).Create(data).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				err = constant.SwipeHasAlreadyExist
			}

			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"data": data,
			})
			return err
		}

		return nil
	}, nil); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return err
	}

	return nil
}

// GetSwipe implements IPostgresRepository.
func (u *PostgresRepository) GetSwipe(ctx context.Context, swiperID uint, date time.Time) (result []domain.Swipe, err error) {
	segment := logger.StartSegment(ctx, "PostgresRepository.GetSwipe")
	defer segment.End()

	if err := hystrix.DoC(ctx, constant.HystrixPostgres, func(ctx context.Context) error {
		if err := u.db.GORM(ctx).Model(&domain.Swipe{}).Where("swiper_id = ? AND date = ?", swiperID, date).Find(&result).Error; err != nil {
			if err.Error() == gorm.ErrRecordNotFound.Error() {
				err = constant.SwipeCanNotBeFound
			}

			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"swiper_id": swiperID,
				"date":      date,
			})
			return err
		}

		return nil
	}, nil); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"swiper_id": swiperID,
			"date":      date,
		})
		return result, err
	}

	return result, nil
}
