package repository

import (
	"github.com/fahmialfareza/deals-dating-app-backend/configs/cache/redis"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/database/postgres"
	postgresRepo "github.com/fahmialfareza/deals-dating-app-backend/internal/repository/postgres"
	redisRepo "github.com/fahmialfareza/deals-dating-app-backend/internal/repository/redis"
)

type IRepository interface {
	IUserRepository
	ISwiperRepository
	IPurchaseRepository
}

type Repository struct {
	postgres postgresRepo.IPostgresRepository
	redis    redisRepo.IRedisRepository
}

func NewRepository(postgres postgres.IPostgres, redis redis.IRedis) IRepository {
	return &Repository{
		postgres: postgresRepo.NewPostgresRepository(postgres),
		redis:    redisRepo.NewRedisRepository(redis),
	}
}
