package api

import (
	"fmt"
	"strings"

	"github.com/fahmialfareza/deals-dating-app-backend/configs/cache/redis"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/database/postgres"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/environment"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/logger/newrelic"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/object_storage/imagekit"
	pgUtil "github.com/fahmialfareza/deals-dating-app-backend/pkg/database/postgres"
	imagekitUtil "github.com/fahmialfareza/deals-dating-app-backend/pkg/object_storage/imagekit"
)

type config struct {
	postgres postgres.IPostgres
	redis    redis.IRedis
	newRelic *newrelic.NewRelic
}

func loadConfig(service string) *config {
	// environment
	if err := environment.LoadEnvironment(service); err != nil {
		panic(err)
	}

	// database
	pgClient := postgres.NewPostgres(environment.PostgresSetLimits)

	// cache
	redisClient := redis.NewRedis(environment.RedisURL)

	// new relic
	newrelicAppName := fmt.Sprintf("%s - %s - Dating App Backend", strings.ToUpper(environment.Environment), strings.ToUpper(service))
	newRelic, err := newrelic.NewNewRelic(newrelicAppName, environment.NewRelicLicense)
	if err != nil {
		panic(err)
	}

	// image kit
	imageKit := imagekit.NewImageKit()

	// pkg
	pgUtil.NewPostgresUtil(pgClient)
	imagekitUtil.NewImageKit(imageKit)

	return &config{
		postgres: pgClient,
		redis:    redisClient,
		newRelic: newRelic,
	}
}
