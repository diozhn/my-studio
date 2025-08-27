package database

import (
	"fmt"
	"log"
	"os"

	"my-studio/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL both environment variable not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}

	DB = db

	// Rodar automigração dos models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Artwork{},
	)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	fmt.Println("Database connected and migrated successfully")
}
