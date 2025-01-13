package http

import (
	"context"
	"strings"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
)

func (h *HTTPHandler) authMiddleware(ctx context.Context, authHeader string) (user domain.User, err error) {
	segment := logger.StartSegment(ctx, "HTTPHandler.authMiddleware")
	defer segment.End()

	if authHeader == "" {
		err := constant.Unauthorized
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"auth_header": authHeader,
		})
		return user, err
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		err := constant.Unauthorized
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"auth_header": authHeader,
		})
		return user, err
	}

	// Extract the token part after "Bearer "
	token := authHeader[len("Bearer "):]

	// validate the token
	userID, err := h.usecase.CheckToken(ctx, token)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"auth_header": authHeader,
		})
		return user, constant.TokenNotValid
	}

	// get user detail and validate the access
	user, err = h.usecase.GetProfile(ctx, userID)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"auth_header": authHeader,
		})
		return user, constant.TokenNotValid
	}

	return user, nil
}
