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

	// AutoMigrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.UserDetails{},
		&models.Vendor{},
		&models.Store{},
		&models.Product{},
		&models.Service{})
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// Setup full-text search
	setupFullTextSearch(db)

	// Add common indexes
	createIndexes(db)

	DB = DbInstance{Db: db}
}

func setupFullTextSearch(db *gorm.DB) {
	// Add tsvector columns and indexes
	statements := []string{
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS search_vector tsvector;`,
		`CREATE INDEX IF NOT EXISTS products_search_idx ON products USING gin(search_vector);`,
		`CREATE OR REPLACE FUNCTION products_trigger() RETURNS trigger AS $$
		BEGIN
			NEW.search_vector = to_tsvector('english', NEW.name || ' ' || COALESCE(NEW.description, ''));
			RETURN NEW;
		END
		$$ LANGUAGE plpgsql;`,
		`DROP TRIGGER IF EXISTS products_vector_update ON products;`,
		`CREATE TRIGGER products_vector_update BEFORE INSERT OR UPDATE ON products FOR EACH ROW EXECUTE FUNCTION products_trigger();`,

		`ALTER TABLE services ADD COLUMN IF NOT EXISTS search_vector tsvector;`,
		`CREATE INDEX IF NOT EXISTS services_search_idx ON services USING gin(search_vector);`,
		`CREATE OR REPLACE FUNCTION services_trigger() RETURNS trigger AS $$
		BEGIN
			NEW.search_vector = to_tsvector('english', NEW.name || ' ' || COALESCE(NEW.description, ''));
			RETURN NEW;
		END
		$$ LANGUAGE plpgsql;`,
		`DROP TRIGGER IF EXISTS services_vector_update ON services;`,
		`CREATE TRIGGER services_vector_update BEFORE INSERT OR UPDATE ON services FOR EACH ROW EXECUTE FUNCTION services_trigger();`,
	}

	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			log.Printf("Warning: Error executing search setup statement: %v", err)
		}
	}
}

func createIndexes(db *gorm.DB) {
	indexes := []string{
		// Orders
		`CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_orders_store_id ON orders(store_id);`,
		`CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);`,

		// Store
		`CREATE INDEX IF NOT EXISTS idx_stores_vendor_id ON stores(vendor_id);`,

		// Products
		`CREATE INDEX IF NOT EXISTS idx_products_store_id ON products(store_id);`,

		// Services
		`CREATE INDEX IF NOT EXISTS idx_services_store_id ON services(store_id);`,

		// Vendors
		`CREATE INDEX IF NOT EXISTS idx_vendors_user_id ON vendors(user_id);`,
	}

	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			log.Printf("Warning: Error creating index: %v", err)
		}
	}
}
