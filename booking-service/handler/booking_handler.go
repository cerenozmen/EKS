package handler

import (
	"booking-service/service"
	"github.com/gofiber/fiber/v2"
	
)

type BookingHandler struct {
	BookingService *service.BookingService
}


func NewBookingHandler(s *service.BookingService) *BookingHandler {
	return &BookingHandler{BookingService: s}
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	var req struct {
		UserID  int `json:"user_id"`
		EventID int `json:"event_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek gövdesi"})
	}

	eventName, userID, err := h.BookingService.Register(req.UserID, req.EventID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message":       "Kayıt başarıyla oluşturuldu",
        "user_id":       userID,
        "event_name":    eventName,
    })
}
