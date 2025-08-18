package handler

import (
	"event-service/model"
	"event-service/service"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	Service *service.EventService
}

func NewEventHandler(s *service.EventService) *EventHandler {
	return &EventHandler{Service: s}
}

func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	var e model.Event
	if err := c.BodyParser(&e); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek"})
	}

	createdEvent, err := h.Service.CreateEvent(e)
	if err != nil {
		fmt.Println("CreateEvent hata:", err)
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

	events, err := h.Service.GetEvents(isActive)
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
