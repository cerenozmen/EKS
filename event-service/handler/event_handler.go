package handler

import (
	"event-service/model"
	"event-service/service"

	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	Service *service.EventService
	UserServiceURL string
}

func NewEventHandler(s *service.EventService, userServiceURL string) *EventHandler {
	return &EventHandler{Service: s, UserServiceURL: userServiceURL}
}

func (h *EventHandler) CheckUserMiddleware(c *fiber.Ctx) error {
	token := c.Get("X-Token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Yetkilendirme token'ı eksik"})
	}

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

	c.Locals("userID", int(userID))

	return c.Next()
}

func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı ID'si mevcut değil"})
	}

	var e model.Event
	if err := c.BodyParser(&e); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek"})
	}

	e.UserId = userID

	createdEvent, err := h.Service.CreateEvent(e)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Etkinlik oluşturulamadı"})
	}

	return c.Status(fiber.StatusCreated).JSON(createdEvent)
}


func (h *EventHandler) GetEvents(c *fiber.Ctx) error {
	query := c.Query("isActive")
	var isActive *bool
	if query != "" {
		val := query == "true"
		isActive = &val
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	events, err := h.Service.GetEvents(isActive, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Etkinlikler alınamadı"})
	}

	return c.JSON(events)
}

func (h *EventHandler) GetEventByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz id"})
	}

	event, err := h.Service.GetEventByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Etkinlik bulunamadı"})
	}

	return c.JSON(event)
}
