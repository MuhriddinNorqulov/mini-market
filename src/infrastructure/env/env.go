package env

import (
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Env struct {
	HttpPort   string `env:"HTTP_PORT,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBSSLMode  string `env:"DB_SSLMODE,required"`

	DBMaxOpenConns        int   `env:"DB_MAX_OPEN_CONNS" envDefault:"20"`
	DBMaxIdleConns        int   `env:"DB_MAX_IDLE_CONNS" envDefault:"20"`
	DBConnMaxLifetimeMins int64 `env:"DB_CONN_MAX_LIFETIME_MINUTES" envDefault:"30"`
	DBConnMaxIdleMins     int64 `env:"DB_CONN_MAX_IDLE_MINUTES" envDefault:"5"`

	RedisAddress  string `env:"REDIS_ADDRESS,required"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`

	CorsAllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS,required"`
	CorsAllowOrigin      []string `env:"CORS_ALLOW_ORIGIN,required"`

	JwtSecret string `env:"JWT_SECRET,required"`

	AccessTimeExpireMinutes   int64 `env:"ACCESS_EXPIRE_MINUTES,required"`
	RefreshTokenExpireMinutes int64 `env:"REFRESH_EXPIRE_MINUTES,required"`

	DefaultAdminUsername     string `env:"DEFAULT_ADMIN_USERNAME,required"`
	DefaultAdminPassword     string `env:"DEFAULT_ADMIN_PASSWORD,required"`
	DefaultUserUsername      string `env:"DEFAULT_USER_USERNAME,required"`
	DefaultUserPassword      string `env:"DEFAULT_USER_PASSWORD,required"`
	DefaultDeveloperUsername string `env:"DEFAULT_DEVELOPER_USERNAME,required"`
	DefaultDeveloperPassword string `env:"DEFAULT_DEVELOPER_PASSWORD,required"`

	OtelEnabled               bool    `env:"OTEL_ENABLED" envDefault:"false"`
	OtelEndpoint              string  `env:"OTEL_OTLP_GRPC_ENDPOINT" envDefault:""`
	OtelIngestionKey          string  `env:"OTEL_INGESTION_KEY" envDefault:""`
	OtelInsecure              bool    `env:"OTEL_OTLP_GRPC_INSECURE" envDefault:"false"`
	OtelServiceNameHTTP       string  `env:"OTEL_SERVICE_NAME_HTTP" envDefault:"mini-market-http"`
	OtelServiceNameAsync      string  `env:"OTEL_SERVICE_NAME_ASYNC" envDefault:"mini-market-async"`
	OtelServiceVersion        string  `env:"OTEL_SERVICE_VERSION" envDefault:"dev"`
	OtelEnvironment           string  `env:"OTEL_ENVIRONMENT" envDefault:"development"`
	OtelTracesSampleRatio     float64 `env:"OTEL_TRACES_SAMPLER_ARG" envDefault:"1.0"`
	OtelMetricIntervalSeconds int64   `env:"OTEL_METRIC_INTERVAL_SECONDS" envDefault:"60"`

	SentryEnabled      bool   `env:"SENTRY_ENABLED" envDefault:"false"`
	SentryDSN          string `env:"SENTRY_DSN" envDefault:""`
	SentrySendBodies   bool   `env:"SENTRY_SEND_BODIES" envDefault:"false"`
	SentryFlushSeconds int64  `env:"SENTRY_FLUSH_SECONDS" envDefault:"5"`
}

// @inject
func NewEnv() *Env {
	return parseEnv()
}

func parseEnv() *Env {
	loadEnv()

	cfg := &Env{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	return cfg
}

func loadEnv() {
	if os.Getenv("CONTAINER_MODE") != "1" {
		_ = godotenv.Load(".env.local")
		_ = godotenv.Load(".env")
	}
}
