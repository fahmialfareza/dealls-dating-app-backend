package environment

import "time"

var (
	Environment string
	Port        string

	PostgreSQLDSN     string
	PostgresSetLimits bool

	RedisURL        string
	RedisExpireTime time.Duration
	RedisType       string

	JWTSecret string

	HystrixTimeout                int
	HystrixMaxConcurrentRequests  int
	HystrixErrorPercentThreshold  int
	HystrixRequestVolumeThreshold int
	HystrixSleepWindow            int

	NewRelicLicense string

	ImageKitPublicAPIKey  string
	ImageKitPrivateAPIKey string
	ImageKitURLEndpoint   string
)
