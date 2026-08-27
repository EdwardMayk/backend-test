package http

import (
	"os"
	"strings"

	"matrix-orchestrator/domain"
	"matrix-orchestrator/usecase"

	"github.com/gofiber/fiber/v2"
)

type MatrixRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MatrixHandler struct {
	usecase usecase.MatrixUsecase
}

func NewMatrixHandler(app *fiber.App, uc usecase.MatrixUsecase) {
	handler := &MatrixHandler{usecase: uc}

	api := app.Group("/api/v1")

	api.Post("/auth/login", handler.Login)

	matrixGroup := api.Group("/matrix", JWTMiddleware())
	matrixGroup.Post("/rotate", handler.Rotate)
	matrixGroup.Post("/qr", handler.QR)
}

func (h *MatrixHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payload json invalido",
		})
	}

	authUsername := os.Getenv("AUTH_USERNAME")
	authPassword := os.Getenv("AUTH_PASSWORD")

	if req.Username == authUsername && req.Password == authPassword {
		token, err := domain.GenerateToken(req.Username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error al generar el token",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "credenciales invalidas",
	})
}

func (h *MatrixHandler) Rotate(c *fiber.Ctx) error {
	var req MatrixRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payload json invalido",
		})
	}

	if len(req.Matrix) == 0 || len(req.Matrix[0]) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "la matriz no puede estar vacia",
		})
	}

	authHeader := c.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimPrefix(token, "bearer ")

	rotated, stats, err := h.usecase.ProcessRotation(token, req.Matrix)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"rotated_matrix": rotated,
		"statistics":     stats,
	})
}

func (h *MatrixHandler) QR(c *fiber.Ctx) error {
	var req MatrixRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payload json invalido",
		})
	}

	if len(req.Matrix) == 0 || len(req.Matrix[0]) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "la matriz no puede estar vacia",
		})
	}

	authHeader := c.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimPrefix(token, "bearer ")

	q, r, verified, stats, err := h.usecase.ProcessQR(token, req.Matrix)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"q_matrix":                q,
		"r_matrix":                r,
		"mathematically_verified": verified,
		"statistics":              stats,
	})
}
