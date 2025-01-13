package usecase

import (
	"context"
	"time"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/generator"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
)

type ISwipeUsecase interface {
	GetSwipe(ctx context.Context, swiperID uint) (domain.User, error)
	Swipe(ctx context.Context, swiperID uint, swipedID uint, swipeType string) error
}

// GetSwipe implements IUsecase.
func (u *Usecase) GetSwipe(ctx context.Context, swiperID uint) (result domain.User, err error) {
	segment := logger.StartSegment(ctx, "Usecase.GetSwipe")
	defer segment.End()

	now := time.Now()
	swipes, err := u.validateSwipe(ctx, swiperID, now)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"swiper_id": swiperID,
		})
		return result, err
	}

	// collect the swipes to exclude the swipes that are already swiped
	swipedIDs := []uint{swiperID}
	for _, swipe := range swipes {
		swipedIDs = append(swipedIDs, swipe.SwipedID)
	}

	// get the profile not in swiped
	profiles, err := u.repository.GetUserProfile(ctx, swipedIDs)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"swiper_id": swiperID,
		})
		return result, err
	}

	// collect the profilesIDs to be selected into
	var profileIDs []uint
	for _, profile := range profiles {
		profileIDs = append(profileIDs, profile.ID)
	}
	randomizeProfile := generator.RandomizeUint(ctx, profileIDs)

	// get the random profile
	result, err = u.repository.GetUserDetail(ctx, &randomizeProfile, nil)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"swiper_id": swiperID,
		})
		return result, err
	}

	return result, nil
}

// Swipe implements IUsecase.
func (u *Usecase) Swipe(ctx context.Context, swiperID uint, swipedID uint, swipeType string) error {
	segment := logger.StartSegment(ctx, "Usecase.Swipe")
	defer segment.End()

	if swipedID == swiperID {
		err := constant.CanNotSwipeYourSelf
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"swiper_id":   swiperID,
			"swiped_id":   swipedID,
			"swiped_type": swipeType,
		})
		return err
	}

	now := time.Now()
	if _, err := u.validateSwipe(ctx, swiperID, now); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"swiper_id":   swiperID,
			"swiped_id":   swipedID,
			"swiped_type": swipeType,
		})
		return err
	}

	if err := u.repository.CreateSwipe(ctx, &domain.Swipe{
		SwiperID: swiperID,
		SwipedID: swipedID,
		Type:     swipeType,
		Date:     now,
	}); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"swiper_id":   swiperID,
			"swiped_id":   swipedID,
			"swiped_type": swipeType,
		})
		return err
	}

	return nil
}

func (u *Usecase) validateSwipe(ctx context.Context, swiperID uint, date time.Time) (swipes []domain.Swipe, err error) {
	segment := logger.StartSegment(ctx, "Usecase.validateSwipe")
	defer segment.End()

	// get swiper information
	swiper, err := u.repository.GetUserDetail(ctx, &swiperID, nil)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"swiper_id": swiperID,
		})
		return swipes, err
	}

	// get swipes in the same date
	swipes, err = u.repository.GetSwipe(ctx, swiperID, date)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"swiper_id": swiperID,
		})
		return swipes, err
	}
	if len(swipes) >= 10 && swiper.IsPremium == false {
		err := constant.YouNeedToBePremiumMember
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"swiper_id": swiperID,
		})
		return swipes, err
	}

	return swipes, nil
}
