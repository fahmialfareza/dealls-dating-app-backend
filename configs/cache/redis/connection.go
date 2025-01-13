package redis

import (
	"context"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/redis/go-redis/v9"
)

func openRedis(redisURL string) (*redis.Client, error) {
	var rdb *redis.Client

	conn := func() error {
		options, err := redis.ParseURL(redisURL)
		if err != nil {
			return err
		}

		rdb = redis.NewClient(options)
		rdb.AddHook(nrredis.NewHook(options))

		if err := rdb.Ping(context.Background()).Err(); err != nil {
			return err
		}

		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 10 * time.Second

	err := backoff.Retry(conn, expBackoff)
	if err != nil {
		return rdb, err
	}

	return rdb, nil
}

func openRedisCluster(redisURL string) (*redis.ClusterClient, error) {
	var rdb *redis.ClusterClient

	conn := func() error {
		options, err := redis.ParseClusterURL(redisURL)
		if err != nil {
			return err
		}

		rdb = redis.NewClusterClient(options)
		rdb.AddHook(nrredis.NewHook(nil))

		if err := rdb.Ping(context.Background()).Err(); err != nil {
			return err
		}

		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 10 * time.Second

	err := backoff.Retry(conn, expBackoff)
	if err != nil {
		return rdb, err
	}

	return rdb, nil
}
