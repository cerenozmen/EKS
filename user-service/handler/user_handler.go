package handler

import (
	"user-service/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	// JSON parse et
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Service katmanına gönder
	user, err := h.userService.Register(input.Username, input.Password, input.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        user.ID,
		"username":  user.Username,
		"name":      user.Name,
		"createdAt": user.CreatedAt,
	})
}
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, _, err := h.userService.Login(input.Username, input.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid username or password"})
	}

	// Sadece token dön
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}
func (h *UserHandler) Me(c *fiber.Ctx) error {
	tokenString := c.Get("X-Token")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing X-Token header",
		})
	}

	// Token doğrulama
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.userService.GetJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	// Claims'den kullanıcı bilgilerini al
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot read claims",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":       claims["id"],
		"username": claims["username"],
		"name":     claims["name"],
	})
}
