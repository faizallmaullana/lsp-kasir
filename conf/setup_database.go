package conf

import (
	"fmt"
	"log"
	"os"
	"time"

	"faizalmaulana/lsp/models/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupDatabaseConnection creates a GORM DB connection and runs automigrations for entities.
func SetupDatabaseConnection(dbHost, dbPort, dbUser, dbPass, dbName string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)

	// Use a simple logger in non-production
	gormLogger := logger.Default
	if os.Getenv("GIN_MODE") == "release" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// run auto migrations for our entities
	if err := db.AutoMigrate(
		&entity.Users{},
		&entity.Sessions{},
		&entity.Profiles{},
		&entity.Items{},
		&entity.Transactions{},
		&entity.PivotItemsToTransaction{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	// configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db
}

// CloseDatabaseConnection closes the underlying sql.DB connection pool.
func CloseDatabaseConnection(db *gorm.DB) {
	if db == nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}
