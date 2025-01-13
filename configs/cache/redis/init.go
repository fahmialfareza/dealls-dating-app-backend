package redis

import (
	"context"
	"time"

	"github.com/fahmialfareza/deals-dating-app-backend/configs/environment"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	goRedis "github.com/redis/go-redis/v9"
)

type IRedis interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Scan(ctx context.Context, pattern string) (result []string, err error)
	Del(ctx context.Context, key string) error

	XAdd(ctx context.Context, stream string, values interface{}) (entryID string, err error)
	XGroupCreateMkStream(ctx context.Context, stream string, group string, start string) error
	XReadGroup(ctx context.Context, a *goRedis.XReadGroupArgs) ([]goRedis.XStream, error)
	XAck(ctx context.Context, stream string, group string, ids ...string) (int64, error)
	XDel(ctx context.Context, stream string, ids ...string) error
	XPending(ctx context.Context, stream string, group string) (*goRedis.XPending, error)
	XPendingExt(ctx context.Context, a *goRedis.XPendingExtArgs) ([]goRedis.XPendingExt, error)
	XClaim(ctx context.Context, a *goRedis.XClaimArgs) ([]goRedis.XMessage, error)
}

type RedisClient struct {
	client *goRedis.Client
}

type RedisClusterClient struct {
	client *goRedis.ClusterClient
}

func NewRedis(redisURL string) IRedis {
	if environment.Environment == constant.ProductionEnvironment {
		redisClusterClient, err := openRedisCluster(redisURL)
		if err != nil {
			panic(err)
		}

		return &RedisClusterClient{
			client: redisClusterClient,
		}
	}

	redisClient, err := openRedis(redisURL)
	if err != nil {
		panic(err)
	}

	return &RedisClient{
		client: redisClient,
	}
}
