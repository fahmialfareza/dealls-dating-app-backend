package environment

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnvironment(service string) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	// get the environment
	Environment = os.Getenv("ENVIRONMENT")
	Port = os.Getenv("PORT")
	PostgreSQLDSN = os.Getenv("POSTGRES_DSN")
	PostgresSetLimitsStr := os.Getenv("POSTGRES_SETLIMITS")
	RedisURL = os.Getenv("REDIS_URL")
	RedisExpireTimeStr := os.Getenv("REDIS_EXPIRE_TIME")
	RedisExpireTimeInt, err := strconv.Atoi(RedisExpireTimeStr)
	if err != nil {
		return err
	}
	RedisExpireTime = time.Minute * time.Duration(RedisExpireTimeInt)

	JWTSecret = os.Getenv("JWT_SECRET")

	if PostgresSetLimitsStr == "true" {
		PostgresSetLimits = true
	}

	// Hystrix
	HystrixTimeoutString := os.Getenv("HYSTRIX_TIMEOUT")
	HystrixTimeout, err = strconv.Atoi(HystrixTimeoutString)
	if err != nil {
		return err
	}
	HystrixMaxConcurrentRequestsString := os.Getenv("HYSTRIX_MAX_CONCURRENT_REQUESTS")
	HystrixMaxConcurrentRequests, err = strconv.Atoi(HystrixMaxConcurrentRequestsString)
	if err != nil {
		return err
	}
	HystrixErrorPercentThresholdString := os.Getenv("HYSTRIX_ERROR_PERCENT_THRESHOLD")
	HystrixErrorPercentThreshold, err = strconv.Atoi(HystrixErrorPercentThresholdString)
	if err != nil {
		return err
	}
	HystrixRequestVolumeThresholdString := os.Getenv("HYSTRIX_REQUEST_VOLUME_THRESHOLD")
	HystrixRequestVolumeThreshold, err = strconv.Atoi(HystrixRequestVolumeThresholdString)
	if err != nil {
		return err
	}
	HystrixSleepWindowString := os.Getenv("HYSTRIX_SLEEP_WINDOW")
	HystrixSleepWindow, err = strconv.Atoi(HystrixSleepWindowString)
	if err != nil {
		return err
	}

	// New Relic
	NewRelicLicense = os.Getenv("NEW_RELIC_LICENSE")

	// Imagekit
	ImageKitPublicAPIKey = os.Getenv("IMAGEKIT_PUBLIC_API_KEY")
	ImageKitPrivateAPIKey = os.Getenv("IMAGEKIT_PRIVATE_API_KEY")
	ImageKitURLEndpoint = os.Getenv("IMAGEKIT_URL_ENDPOINT")

	return nil
}
