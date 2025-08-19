package main

import (
	"booking-service/config"
	"booking-service/handler"

	"booking-service/repository"
	"booking-service/service"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	bookingRepo := repository.NewBookingRepository(db)
	bookingService := service.NewBookingService(bookingRepo, cfg.EventServiceURL)
	bookingHandler := handler.NewBookingHandler(bookingService)

	app.Post("/bookings", bookingHandler.CreateBooking)
	app.Delete("/bookings/:userID/:eventID", bookingHandler.CancelBooking)

	fmt.Printf("Booking Service, Fiber ile %s portunda çalışıyor...\n", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		panic(err)
	}
}
