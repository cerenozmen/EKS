package handler

import (
	"booking-service/service"

	"encoding/json"

	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type BookingHandler struct {
	BookingService *service.BookingService
	UserServiceURL string
}

func NewBookingHandler(s *service.BookingService, userServiceURL string) *BookingHandler {
	return &BookingHandler{BookingService: s, UserServiceURL: userServiceURL}
}
func (h *BookingHandler) CheckUserMiddleware(c *fiber.Ctx) error {
	token := c.Get("X-Token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Yetkilendirme token'ı eksik"})
	}

	// User Service'e token ile istek at
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/me", h.UserServiceURL), nil) // User Service adresi ve /me endpoint'i
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "İstek oluşturma hatası"})
	}
	req.Header.Add("X-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User Service'e bağlantı hatası"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz yetkilendirme token'ı"})
	}

	// Yanıttan kullanıcı ID'sini al
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcı bilgileri okunamadı"})
	}

	var user map[string]interface{}
	if err := json.Unmarshal(body, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JSON çözümleme hatası"})
	}

	userID, ok := user["id"].(float64)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcı ID'si okunamadı"})
	}

	// Kullanıcı ID'sini Fiber context'ine kaydet
	c.Locals("userID", int(userID))

	return c.Next()
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	// Middleware'den kullanıcı ID'sini al
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı ID'si mevcut değil"})
	}

	var req struct {
		// user_id artık istek gövdesinde beklenmiyor
		EventID int `json:"event_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek gövdesi"})
	}

	eventName, err := h.BookingService.Register(userID, req.EventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Kayıt başarıyla oluşturuldu",
		"event_name": eventName,
	})
}

// CancelBooking handler fonksiyonu güncellenmiştir.
func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	// Middleware'den kullanıcı ID'sini al
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı ID'si mevcut değil"})
	}

	var req struct {
		// user_id artık istek gövdesinde beklenmiyor
		EventID int `json:"event_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek gövdesi"})
	}

	eventName, err := h.BookingService.CancelByIDs(userID, req.EventID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Kayıt başarıyla iptal edildi",
		"event_name": eventName,
	})
}
