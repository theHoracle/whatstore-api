package database

import (
	"log"
	"os"

	"github.com/theHoracle/whatstore-api/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var DB DbInstance

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("COULD NOT CONNECT TO DB\n " + err.Error())
		os.Exit(2)
	}

	log.Println("Connected to DB successfully")

	db.Logger = logger.Default.LogMode(logger.Info)

	// TP
	log.Println("Running migrations")

	migrationPath := os.Getenv("MIGRATION_PATH")
	if migrationPath == "" {
		migrationPath = "./db/migrations" // Default path if env var is not set
	}

	err = db.AutoMigrate(&models.User{}, &models.Vendor{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	DB = DbInstance{Db: db}
}
