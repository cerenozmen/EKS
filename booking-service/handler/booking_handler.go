package handler

import (
	"booking-service/service"
	"strconv"

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

	eventName, err := h.BookingService.Register(req.UserID, req.EventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Kayıt başarıyla oluşturuldu",
		"event_name": eventName,
	})
}
func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {

	userID, err := strconv.Atoi(c.Params("userID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz kullanıcı ID'si"})
	}

	eventID, err := strconv.Atoi(c.Params("eventID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz etkinlik ID'si"})
	}

	eventName, err := h.BookingService.CancelByIDs(userID, eventID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Kayıt başarıyla iptal edildi",
		"event_name": eventName,
	})
}