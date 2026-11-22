package config

import (
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mengambil sql.DB dari GORM: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Gagal set dialect goose: %v", err)
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatalf("Gagal menjalankan migrasi database: %v", err)
	}
	fmt.Println("Database migrations applied successfully...")
}
