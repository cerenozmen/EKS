package main

import (
	"event-service/config"
	"event-service/handler"
	"event-service/repository"
	"event-service/service"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config yüklenirken hata oluştu: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DatabaseHost,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
		cfg.DatabasePort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Veritabanına bağlanılamadı: " + err.Error())
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: "",
		DB:       0,
	})

	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo, rdb)
	eventHandler := handler.NewEventHandler(eventService, cfg.UserServiceURL)

	app.Post("/events", eventHandler.CheckUserMiddleware, eventHandler.CreateEvent)
	app.Get("/events", eventHandler.GetEvents)
	app.Get("/events/:id", eventHandler.GetEventByID)

	fmt.Printf("Event Service Fiber ile %s portunda çalışıyor...\n", cfg.AppPort)
	logErr := app.Listen(":" + cfg.AppPort)
	if logErr != nil {
		panic(logErr)
	}
}
