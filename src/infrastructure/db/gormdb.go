package db

import (
	"fmt"
	"mini-market/src/infrastructure/env"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

// @inject
func NewGormDB(env *env.Env) *gorm.DB {
	db, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
				env.DBHost, env.DBUser, env.DBPassword, env.DBName, env.DBPort, env.DBSSLMode,
			),
		),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		panic(fmt.Errorf("[NewGormDB] failed to connect database: %w", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("[NewGormDB] failed to get sql.DB: %w", err))
	}
	sqlDB.SetMaxOpenConns(env.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(env.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(env.DBConnMaxLifetimeMins) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(env.DBConnMaxIdleMins) * time.Minute)

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
	return db
}
