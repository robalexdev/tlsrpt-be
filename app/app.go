package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func serve() {
	// Database setup
	dsn := fmt.Sprintf("host=db user=postgres password=%s dbname=tlsrpt port=5432 sslmode=disable TimeZone=UTC", os.Getenv("POSTGRES_PASSWORD"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}

	// Auto-migration
	db.AutoMigrate(&PolicyFailureDetail{})
	db.AutoMigrate(&Policy{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&User{})

	// Admin secret
	//adminSecret := os.Getenv("ADMIN_SECRET")

	// Gin setup
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	setupRoutes(r, db)
	r.Run()
}

func main() {
	serve()
}
