package converter

import (
	"context"
	"encoding/json"

	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func AutoMap(ctx context.Context, from interface{}, to interface{}) error {
	segment := logger.StartSegment(ctx, "converter.AutoMap")
	defer segment.End()

	dataByte, err := json.Marshal(from)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"from": from,
			"to":   to,
		})
		return err
	}

	if err := json.Unmarshal(dataByte, to); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"from": from,
			"to":   to,
		})
		return err
	}

	return nil
}

func HashPassword(ctx context.Context, password string) (string, error) {
	segment := logger.StartSegment(ctx, "converter.HashPassword")
	defer segment.End()

	// generate the password with bcrypt
	hashedPasswordByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"password": password,
		})
		return "", err
	}
	hashedPassword := string(hashedPasswordByte)

	return hashedPassword, nil
}
