package main

import (
	"fmt"
	"log"
	"user-service/config"
	"user-service/handler"
	"user-service/repository"
	"user-service/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

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
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userService)

	app := fiber.New()

	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)

	
	fmt.Printf("User Service, Fiber ile %s portunda çalışıyor...\n", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		panic(err)
	}
}
