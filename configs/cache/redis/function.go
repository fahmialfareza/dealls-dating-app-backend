package redis

import (
	"context"
	"time"

	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
	goRedis "github.com/redis/go-redis/v9"
	redis "github.com/redis/go-redis/v9"
)

/* Redis Client */
func (rdc *RedisClient) Ping(ctx context.Context) error {
	segment := logger.StartSegment(ctx, "RedisClient.Ping")
	defer segment.End()

	if err := rdc.client.Ping(ctx).Err(); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		return err
	}

	return nil
}

func (rdc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	segment := logger.StartSegment(ctx, "RedisClient.Set")
	defer segment.End()

	if err := rdc.client.Set(ctx, key, value, expiration).Err(); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"key":        key,
			"value":      value,
			"expiration": expiration,
		})
		return err
	}

	return nil
}

func (rdc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	segment := logger.StartSegment(ctx, "RedisClient.Get")
	defer segment.End()

	result, err := rdc.client.Get(ctx, key).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"key": key,
		})
		return result, err
	}

	return result, nil
}

func (rdc *RedisClient) Del(ctx context.Context, key string) error {
	segment := logger.StartSegment(ctx, "RedisClient.Del")
	defer segment.End()

	if err := rdc.client.Del(ctx, key).Err(); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"key": key,
		})
		return err
	}

	return nil
}

func (rdc *RedisClient) Scan(ctx context.Context, pattern string) (result []string, err error) {
	segment := logger.StartSegment(ctx, "RedisClient.Scan")
	defer segment.End()

	iter := rdc.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		result = append(result, iter.Val())
	}
	if err := iter.Err(); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"pattern": pattern,
		})
		return result, err
	}

	return result, nil
}

// XAdd implements IRedis.
func (rdc *RedisClient) XAdd(ctx context.Context, stream string, values interface{}) (entryID string, err error) {
	segment := logger.StartSegment(ctx, "RedisClient.XAdd")
	defer segment.End()

	entryID, err = rdc.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream, // Name of the stream
		Values: values,
	}).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"stream": stream,
			"values": values,
		})
		return entryID, err
	}

	return entryID, nil
}

func (rdc *RedisClient) XGroupCreateMkStream(ctx context.Context, stream string, group string, start string) error {
	segment := logger.StartSegment(ctx, "RedisClient.XGroupCreateMkStream")
	defer segment.End()

	if err := rdc.client.XGroupCreateMkStream(ctx, stream, group, start).Err(); err != nil && err != redis.Nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"stream": stream,
				"group":  group,
				"start":  start,
			})
			return err
		}
	}

	return nil
}

// XReadGroup implements IRedis.
func (rdc *RedisClient) XReadGroup(ctx context.Context, a *goRedis.XReadGroupArgs) ([]goRedis.XStream, error) {
	segment := logger.StartSegment(ctx, "RedisClient.XReadGroup")
	defer segment.End()

	streams, err := rdc.client.XReadGroup(ctx, a).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"a": a,
		})
		return streams, err
	}

	return streams, nil
}

// XAck implements IRedis.
func (rdc *RedisClient) XAck(ctx context.Context, stream string, group string, ids ...string) (int64, error) {
	segment := logger.StartSegment(ctx, "RedisClient.XAck")
	defer segment.End()

	result, err := rdc.client.XAck(ctx, stream, group, ids...).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"stream": stream,
			"group":  group,
			"ids":    ids,
		})
		return result, err
	}

	return result, nil
}

// XDel implements IRedis.
func (rdc *RedisClient) XDel(ctx context.Context, stream string, ids ...string) error {
	segment := logger.StartSegment(ctx, "RedisClient.XDel")
	defer segment.End()

	if err := rdc.client.XDel(ctx, stream, ids...).Err(); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"stream": stream,
			"ids":    ids,
		})
		return err
	}

	return nil
}

// XPending implements IRedis.
func (rdc *RedisClient) XPending(ctx context.Context, stream string, group string) (*goRedis.XPending, error) {
	segment := logger.StartSegment(ctx, "RedisClient.XPending")
	defer segment.End()

	result, err := rdc.client.XPending(ctx, stream, group).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"stream": stream,
			"group":  group,
		})
		return result, err
	}

	return result, nil
}

// XClaim implements IRedis.
func (rdc *RedisClient) XClaim(ctx context.Context, a *goRedis.XClaimArgs) ([]goRedis.XMessage, error) {
	segment := logger.StartSegment(ctx, "RedisClient.XClaim")
	defer segment.End()

	result, err := rdc.client.XClaim(ctx, a).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"a": a,
		})
		return result, err
	}

	return result, nil
}

// XPendingExt implements IRedis.
func (rdc *RedisClient) XPendingExt(ctx context.Context, a *goRedis.XPendingExtArgs) ([]goRedis.XPendingExt, error) {
	segment := logger.StartSegment(ctx, "RedisClient.XPendingExt")
	defer segment.End()

	result, err := rdc.client.XPendingExt(ctx, a).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"a": a,
		})
		return result, err
	}

	return result, nil
}

/* Redis Cluster Client */
func (rdc *RedisClusterClient) Ping(ctx context.Context) error {
	segment := logger.StartSegment(ctx, "RedisClusterClient.Ping")
	defer segment.End()

	if err := rdc.client.Ping(ctx).Err(); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		return err
	}

	return nil
}

func (rdc *RedisClusterClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	segment := logger.StartSegment(ctx, "RedisClusterClient.Set")
	defer segment.End()

	if err := rdc.client.Set(ctx, key, value, expiration).Err(); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"key":        key,
			"value":      value,
			"expiration": expiration,
		})
		return err
	}

	return nil
}

func (rdc *RedisClusterClient) Get(ctx context.Context, key string) (string, error) {
	segment := logger.StartSegment(ctx, "RedisClusterClient.Get")
	defer segment.End()

	result, err := rdc.client.Get(ctx, key).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"key": key,
		})
		return result, err
	}

	return result, nil
}

func (rdc *RedisClusterClient) Del(ctx context.Context, key string) error {
	segment := logger.StartSegment(ctx, "RedisClusterClient.Del")
	defer segment.End()

	if err := rdc.client.Del(ctx, key).Err(); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"key": key,
		})
		return err
	}

	return nil
}

func (rdc *RedisClusterClient) Scan(ctx context.Context, pattern string) (result []string, err error) {
	segment := logger.StartSegment(ctx, "RedisClusterClient.Scan")
	defer segment.End()

	iter := rdc.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		result = append(result, iter.Val())
	}
	if err := iter.Err(); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"pattern": pattern,
		})
		return result, err
	}

	return result, nil
}

// XAdd implements IRedis.
func (rdc *RedisClusterClient) XAdd(ctx context.Context, stream string, values interface{}) (entryID string, err error) {
	segment := logger.StartSegment(ctx, "RedisClusterClient.XAdd")
	defer segment.End()

	entryID, err = rdc.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream, // Name of the stream
		Values: values,
	}).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"stream": stream,
			"values": values,
		})
		return entryID, err
	}

	return entryID, nil
}

// XGroupCreateMkStream implements IRedis.
func (rdc *RedisClusterClient) XGroupCreateMkStream(ctx context.Context, stream string, group string, start string) error {
	segment := logger.StartSegment(ctx, "RedisClusterClient.XGroupCreateMkStream")
	defer segment.End()

	if err := rdc.client.XGroupCreateMkStream(ctx, stream, group, start).Err(); err != nil && err != redis.Nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"stream": stream,
				"group":  group,
				"start":  start,
			})
			return err
		}
	}

	return nil
}

// XReadGroup implements IRedis.
func (rdc *RedisClusterClient) XReadGroup(ctx context.Context, a *goRedis.XReadGroupArgs) ([]goRedis.XStream, error) {
	segment := logger.StartSegment(ctx, "RedisClusterClient.XReadGroup")
	defer segment.End()

	streams, err := rdc.client.XReadGroup(ctx, a).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"a": a,
		})
		return streams, err
	}

	return streams, nil
}

// XAck implements IRedis.
func (rdc *RedisClusterClient) XAck(ctx context.Context, stream string, group string, ids ...string) (int64, error) {
	segment := logger.StartSegment(ctx, "RedisClusterClient.XAck")
	defer segment.End()

	result, err := rdc.client.XAck(ctx, stream, group, ids...).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"stream": stream,
			"group":  group,
			"ids":    ids,
		})
		return result, err
	}

	return result, nil
}

// XDel implements IRedis.
func (rdc *RedisClusterClient) XDel(ctx context.Context, stream string, ids ...string) error {
	segment := logger.StartSegment(ctx, "RedisClusterClient.XDel")
	defer segment.End()

	if err := rdc.client.XDel(ctx, stream, ids...).Err(); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"stream": stream,
			"ids":    ids,
		})
		return err
	}

	return nil
}

// XPending implements IRedis.
func (rdc *RedisClusterClient) XPending(ctx context.Context, stream string, group string) (*goRedis.XPending, error) {
	segment := logger.StartSegment(ctx, "RedisClusterClient.XPending")
	defer segment.End()

	result, err := rdc.client.XPending(ctx, stream, group).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"stream": stream,
			"group":  group,
		})
		return result, err
	}

	return result, nil
}

// XClaim implements IRedis.
func (rdc *RedisClusterClient) XClaim(ctx context.Context, a *goRedis.XClaimArgs) ([]goRedis.XMessage, error) {
	segment := logger.StartSegment(ctx, "RedisClusterClient.XClaim")
	defer segment.End()

	result, err := rdc.client.XClaim(ctx, a).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"a": a,
		})
		return result, err
	}

	return result, nil
}

// XPendingExt implements IRedis.
func (rdc *RedisClusterClient) XPendingExt(ctx context.Context, a *goRedis.XPendingExtArgs) ([]goRedis.XPendingExt, error) {
	segment := logger.StartSegment(ctx, "RedisClusterClient.XPendingExt")
	defer segment.End()

	result, err := rdc.client.XPendingExt(ctx, a).Result()
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"a": a,
		})
		return result, err
	}

	return result, nil
}
