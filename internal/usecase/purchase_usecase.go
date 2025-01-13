package usecase

import (
	"context"
	"time"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
)

type IPurchaseUsecase interface {
	Purchase(ctx context.Context, userID uint) error
}

// Purchase implements IUsecase.
func (u *Usecase) Purchase(ctx context.Context, userID uint) error {
	segment := logger.StartSegment(ctx, "Usecase.Purchase")
	defer segment.End()

	// make a purchase
	now := time.Now()
	if err := u.repository.CreatePurchase(ctx, &domain.Purchase{
		UserID:       userID,
		PurchaseDate: now,
		PackageType:  "premium",
	}); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"user_id": userID,
		})
		return err
	}

	// update the user is premium
	user, err := u.repository.GetUserDetail(ctx, &userID, nil)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"user_id": userID,
		})
		return err
	}
	user.IsPremium = true
	_, err = u.repository.UpsertUser(ctx, user)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"user_id": userID,
		})
		return err
	}

	return nil
}
