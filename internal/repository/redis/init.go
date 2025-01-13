package postgres

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/cache/redis"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/environment"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
)

type IRedisRepository interface {
	IUserRedisRepository
}

type RedisRepository struct {
	redis redis.IRedis
}

func NewRedisRepository(redis redis.IRedis) IRedisRepository {
	hystrix.ConfigureCommand(constant.HystrixRedis, hystrix.CommandConfig{
		Timeout:                environment.HystrixTimeout,
		MaxConcurrentRequests:  environment.HystrixMaxConcurrentRequests,
		ErrorPercentThreshold:  environment.HystrixErrorPercentThreshold,
		RequestVolumeThreshold: environment.HystrixRequestVolumeThreshold,
		SleepWindow:            environment.HystrixSleepWindow,
	})

	return &RedisRepository{
		redis: redis,
	}
}
