package postgres

import (
	"context"
	"strings"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
)

type IPurchasePostgresRepository interface {
	InsertPurchase(ctx context.Context, data *domain.Purchase) error
}

// InsertPurchase implements IPostgresRepository.
func (u *PostgresRepository) InsertPurchase(ctx context.Context, data *domain.Purchase) error {
	segment := logger.StartSegment(ctx, "PostgresRepository.InsertPurchase")
	defer segment.End()

	if err := hystrix.DoC(ctx, constant.HystrixPostgres, func(ctx context.Context) error {
		// create
		if err := u.db.GORM(ctx).Create(data).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				err = constant.PurchaseHasAlreadyExist
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
