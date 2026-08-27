package http

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"matrix-orchestrator/usecase"

	"github.com/gofiber/fiber/v2"
)

type mockMatrixUsecase struct{}

func (m *mockMatrixUsecase) ProcessRotation(token string, matrix [][]float64) ([][]float64, *usecase.StatsResponse, error) {
	return matrix, &usecase.StatsResponse{Max: 1, Min: 0, Average: 0.5, SumTotal: 1, IsDiagonal: false}, nil
}

func (m *mockMatrixUsecase) ProcessQR(token string, matrix [][]float64) ([][]float64, [][]float64, bool, *usecase.StatsResponse, error) {
	return matrix, matrix, true, &usecase.StatsResponse{Max: 1, Min: 0, Average: 0.5, SumTotal: 1, IsDiagonal: false}, nil
}

func TestHTTPIntegration(t *testing.T) {
	t.Setenv("JWT_SECRET", "clave-secreta-de-prueba-123456789")
	t.Setenv("AUTH_USERNAME", "usuario-de-prueba")
	t.Setenv("AUTH_PASSWORD", "contrasena-de-prueba")
	
	app := fiber.New()
	mockUC := &mockMatrixUsecase{}
	NewMatrixHandler(app, mockUC)

	t.Run("Login con credenciales invalidas debe fallar", func(t *testing.T) {
		payload, _ := json.Marshal(LoginRequest{Username: "usuario-de-prueba", Password: "contrasena-incorrecta"})
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error inesperado en test: %v", err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("se esperaba estado 401, se obtuvo %d", resp.StatusCode)
		}
	})

	var token string
	t.Run("Login con credenciales validas debe tener exito y retornar un token", func(t *testing.T) {
		payload, _ := json.Marshal(LoginRequest{Username: "usuario-de-prueba", Password: "contrasena-de-prueba"})
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error inesperado en test: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("se esperaba estado 200, se obtuvo %d", resp.StatusCode)
		}

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		token = body["token"]
		if token == "" {
			t.Error("se esperaba un token en la respuesta, se obtuvo vacio")
		}
	})

	t.Run("Rutas protegidas deben rechazar peticiones sin token", func(t *testing.T) {
		payload, _ := json.Marshal(MatrixRequest{Matrix: [][]float64{{1, 2}, {3, 4}}})
		req := httptest.NewRequest("POST", "/api/v1/matrix/rotate", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error inesperado en test: %v", err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("se esperaba estado 401, se obtuvo %d", resp.StatusCode)
		}
	})

	t.Run("Rutas protegidas deben aceptar peticiones con token valido", func(t *testing.T) {
		payload, _ := json.Marshal(MatrixRequest{Matrix: [][]float64{{1, 2}, {3, 4}}})
		req := httptest.NewRequest("POST", "/api/v1/matrix/rotate", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error inesperado en test: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("se esperaba estado 200, se obtuvo %d", resp.StatusCode)
		}
	})
}
