package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/environment"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func openPostgres(setLimits bool) (*gorm.DB, error) {
	var db *gorm.DB

	conn := func() (err error) {
		// GORM specific configuration
		gormConfig := &gorm.Config{
			// Optional: You can customize this based on your use case
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: func() logger.Interface {
				if environment.Environment == constant.StagingEnvironment {
					return logger.Default.LogMode(logger.Info)
				} else {
					return logger.Default.LogMode(logger.Silent)
				}
			}(), // Adjust log level as needed
		}

		// Connect to the database using GORM
		db, err = gorm.Open(postgres.Open(environment.PostgreSQLDSN), gormConfig)
		if err != nil {
			return err
		}

		// Get underlying SQL DB for setting connection limits
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}

		// Set connection limits if required
		if setLimits {
			fmt.Println("setting limits")
			sqlDB.SetMaxOpenConns(10)
			sqlDB.SetMaxIdleConns(10)
		}

		// Define 5 seconds timeout for connecting to the database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Ping the database to ensure connection
		err = sqlDB.PingContext(ctx)
		if err != nil {
			return err
		}

		// Migrate the table
		if err = db.AutoMigrate(&domain.User{}, &domain.Profile{}, &domain.Swipe{}, &domain.Purchase{}); err != nil {
			return err
		}

		return nil
	}

	// Exponential backoff for retrying the connection
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 10 * time.Second

	err := backoff.Retry(conn, expBackoff)
	if err != nil {
		return db, err
	}

	return db, nil
}
