package jwt

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/environment"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/converter"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
)

func SetToken(ctx context.Context, userID uint) (string, error) {
	segment := logger.StartSegment(ctx, "jwt.SetToken")
	defer segment.End()

	currentTime := time.Now().UTC()
	expiredTime := currentTime.Add(24 * time.Hour)

	data := JWTDataToken{
		ID: userID,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"iat":  currentTime.Unix(),
		"exp":  expiredTime.Unix(),
	}).SignedString([]byte(environment.JWTSecret))
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"user_id": userID,
		})
		return token, err
	}

	return token, nil
}

func ExtractJWT(ctx context.Context, jwtToken string) (uint, error) {
	segment := logger.StartSegment(ctx, "jwt.ExtractJWT")
	defer segment.End()

	var data JWTDataToken

	// Parse the JWT token
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// For simplicity, we're using a constant secret key here,
		// but in production, you should load this from a secure location.
		return []byte(environment.JWTSecret), nil
	})

	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"jwt_token": jwtToken,
		})

		return 0, err
	}

	// Check if the token is valid
	if !token.Valid {
		err = constant.TokenNotValid
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"jwt_token": jwtToken,
		})
		return 0, err
	}

	// Extract claims (payload) from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = constant.TokenNotValid
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"jwt_token": jwtToken,
		})
		return 0, err
	}

	if err = converter.AutoMap(ctx, claims["data"], &data); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"jwt_token": jwtToken,
		})
		return 0, err
	}

	return data.ID, nil
}
