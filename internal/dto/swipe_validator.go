package dto

import (
	"context"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
)

func ValidateSwipe(ctx context.Context, data domain.SwipeRequest) error {
	segment := logger.StartSegment(ctx, "domain.ValidateSwipe")
	defer segment.End()

	if data.SwipedID == 0 {
		err := constant.SwipedIDIsRequired
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return err
	}
	if data.SwipeType == "" {
		err := constant.SwipedTypeIsRequired
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return err
	}

	return nil
}
