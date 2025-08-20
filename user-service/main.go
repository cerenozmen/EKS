package main

import (
	"user-service/config"
	"user-service/handler"

	"user-service/repository"
	"user-service/service"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

)

func main() {
	app := fiber.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Konfigürasyon yüklenemedi: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DatabaseHost, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseName, cfg.DatabasePort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Veritabanı bağlantısı kurulamadı: " + err.Error())
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()
	r.POST("/register", userHandler.Register)

	
	fmt.Printf("User Service, Fiber ile %s portunda çalışıyor...\n", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		panic(err)
	}
}
