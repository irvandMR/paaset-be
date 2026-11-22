package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabaseConfig(env *Environment) *gorm.DB {
	databaseURL := env.DatabaseUrl

	if databaseURL == "" {
		log.Fatalf("Database URL is empty")
	}

	logMode := logger.Warn
	if !env.IsProduction() {
		logMode = logger.Info
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		log.Fatalf("Gagal inisialisasi GORM: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mengambil instance sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Gagal ping database: %v", err)
	}

	fmt.Println("Database Connected successfully via GORM...")
	return db
}
