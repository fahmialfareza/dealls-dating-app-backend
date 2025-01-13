package postgres

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/database/postgres"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/environment"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
)

type IPostgresRepository interface {
	IUserPostgresRepository
	ISwipePostgresRepository
	IPurchasePostgresRepository
}

type PostgresRepository struct {
	db postgres.IPostgres
}

func NewPostgresRepository(db postgres.IPostgres) IPostgresRepository {
	hystrix.ConfigureCommand(constant.HystrixPostgres, hystrix.CommandConfig{
		Timeout:                environment.HystrixTimeout,
		MaxConcurrentRequests:  environment.HystrixMaxConcurrentRequests,
		ErrorPercentThreshold:  environment.HystrixErrorPercentThreshold,
		RequestVolumeThreshold: environment.HystrixRequestVolumeThreshold,
		SleepWindow:            environment.HystrixSleepWindow,
	})

	return &PostgresRepository{
		db: db,
	}
}
