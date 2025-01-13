package dto

import (
	"context"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
)

func ValidateAuthRegister(ctx context.Context, data domain.RegisterRequest) error {
	segment := logger.StartSegment(ctx, "domain.ValidateAuthRegister")
	defer segment.End()

	if data.Name == "" {
		err := constant.UserNameIsRequired
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return err
	}
	if data.Email == "" {
		err := constant.UserEmailIsRequired
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return err
	}
	if data.Password == "" {
		err := constant.UserPasswordIsRequired
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return err
	}
	if data.ImageData == "" || data.ImageFileName == "" {
		err := constant.UserImageIsRequired
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return err
	}
	if data.Bio == "" {
		err := constant.UserBioIsRequired
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return err
	}

	return nil
}

func ValidateAuthLogin(ctx context.Context, data domain.LoginRequest) error {
	segment := logger.StartSegment(ctx, "domain.ValidateAuthLogin")
	defer segment.End()

	if data.Email == "" {
		err := constant.UserEmailIsRequired
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return err
	}
	if data.Password == "" {
		err := constant.UserPasswordIsRequired
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return err
	}

	return nil
}
