package database

import (
	"aplikasieoq/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func StartDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error - Cannot Loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, dbPort)

	dsn := config
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error - Connection to Database : ", err)
	}

	fmt.Println("Connection Database Success")

	// Check if tables already exist
	if !tableExists(db, "barangs") {
		db.Debug().AutoMigrate(&models.Barang{})
	}
	if !tableExists(db, "eoqs") {
		db.Debug().AutoMigrate(&models.Eoq{})
	}
	if !tableExists(db, "pemesanans") {
		db.Debug().AutoMigrate(&models.Pemesanan{})
	}
	if !tableExists(db, "penjualans") {
		db.Debug().AutoMigrate(&models.Penjualan{})
	}
	if !tableExists(db, "penyimpanans") {
		db.Debug().AutoMigrate(&models.Penyimpanan{})
	}
	if !tableExists(db, "suppliers") {
		db.Debug().AutoMigrate(&models.Supplier{})
	}
	if !tableExists(db, "users") {
		db.Debug().AutoMigrate(&models.User{})
	}

}

// tableExists checks if a table exists in the database
func tableExists(db *gorm.DB, tableName string) bool {
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?", tableName).Scan(&count).Error
	if err != nil {
		log.Fatalf("failed to check table existence: %v", err)
	}
	return count > 0
}

func GetDB() *gorm.DB {
	return db
}
