// config/database.go

package configs

import (
	"log"
	"time"

	"github.com/skanderphilipp/sisyApi/internal/domain/artist"
	"github.com/skanderphilipp/sisyApi/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDatabaseConnection initializes and returns a GORM database connection
func SetupDatabaseConnection() *gorm.DB {
	dsn := "host=localhost user=your_user password=your_password dbname=your_dbname port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	// Configure connection pool (optional)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database from GORM: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Perform the migration
	err = db.AutoMigrate(&artist.Artist{}, &artist.SocialMedia{}, &model.Stage{}, &model.TimetableEntry{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// Ensure you call this function from somewhere in your application, typically main.go
func CloseDatabaseConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error on getting DB from GORM: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Error on closing db connection: %v", err)
	}
}

/**package configs

import (
	"github.com/skanderphilipp/sisyApi/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=your_user password=your_password dbname=your_dbname port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	DB.AutoMigrate(&models.Artist{})
} */
