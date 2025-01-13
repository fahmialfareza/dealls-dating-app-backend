package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/environment"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type IUserRedisRepository interface {
	// User
	GetUserDetail(ctx context.Context, id *uint, email *string) (domain.User, error)
	SetUserDetail(ctx context.Context, id *uint, email *string, data domain.User) error
	DeleteUserDetail(ctx context.Context, id *uint, email *string) error
}

// DeleteUserDetail implements IRedisRepository.
func (r *RedisRepository) DeleteUserDetail(ctx context.Context, id *uint, email *string) error {
	segment := logger.StartSegment(ctx, "RedisRepository.DeleteUserDetail")
	defer segment.End()

	if err := hystrix.DoC(ctx, constant.HystrixRedis, func(ctx context.Context) error {
		var key string
		if id != nil {
			key = fmt.Sprintf(userDetail, *id)
		} else if email != nil {
			key = fmt.Sprintf(userDetail, *email)
		}

		if err := r.redis.Del(ctx, key); err != nil {
			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"id":    id,
				"email": email,
			})
			return err
		}

		return nil
	}, nil); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"id":    id,
			"email": email,
		})
		return err
	}

	return nil
}

// GetUserDetail implements IRedisRepository.
func (r *RedisRepository) GetUserDetail(ctx context.Context, id *uint, email *string) (result domain.User, err error) {
	segment := logger.StartSegment(ctx, "RedisRepository.GetUserDetail")
	defer segment.End()

	if err := hystrix.DoC(ctx, constant.HystrixRedis, func(ctx context.Context) error {
		var key string
		if id != nil {
			key = fmt.Sprintf(userDetail, *id)
		} else if email != nil {
			key = fmt.Sprintf(userDetail, *email)
		}

		dataString, err := r.redis.Get(ctx, key)
		if err != nil {
			if err.Error() == redis.Nil.Error() {
				return nil
			}

			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"id":    id,
				"email": email,
			})
			return err
		}

		if err = json.Unmarshal([]byte(dataString), &result); err != nil {
			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"id":    id,
				"email": email,
			})
			return err
		}

		return nil
	}, nil); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"id":    id,
			"email": email,
		})
		return result, err
	}

	return result, nil
}

// SetUserDetail implements IRedisRepository.
func (r *RedisRepository) SetUserDetail(ctx context.Context, id *uint, email *string, data domain.User) error {
	segment := logger.StartSegment(ctx, "UserRedisRepo.SetUserDetail")
	defer segment.End()

	if err := hystrix.DoC(ctx, constant.HystrixRedis, func(ctx context.Context) error {
		var key string
		if id != nil {
			key = fmt.Sprintf(userDetail, *id)
		} else if email != nil {
			key = fmt.Sprintf(userDetail, *email)
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"id":    id,
				"email": email,
				"data":  data,
			})
			return err
		}

		dataString := string(dataBytes)
		if err := r.redis.Set(ctx, key, dataString, environment.RedisExpireTime); err != nil {
			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"id":    id,
				"email": email,
				"data":  data,
			})
			return err
		}

		return nil
	}, nil); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"id":    id,
			"email": email,
			"data":  data,
		})
		return err
	}

	return nil
}
