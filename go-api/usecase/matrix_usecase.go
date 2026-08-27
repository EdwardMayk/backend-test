package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"matrix-orchestrator/domain"
)

type MatrixUsecase interface {
	ProcessRotation(token string, matrix [][]float64) ([][]float64, *StatsResponse, error)
	ProcessQR(token string, matrix [][]float64) ([][]float64, [][]float64, bool, *StatsResponse, error)
}

type StatsRequest struct {
	Matrices [][][]float64 `json:"matrices"`
}

type StatsResponse struct {
	Max        float64 `json:"max"`
	Min        float64 `json:"min"`
	Average    float64 `json:"average"`
	SumTotal   float64 `json:"sum_total"`
	IsDiagonal bool    `json:"is_diagonal"`
}

type matrixUsecase struct {
	nodeAPIURL string
	httpClient *http.Client
}

func NewMatrixUsecase() MatrixUsecase {
	nodeURL := os.Getenv("NODE_API_URL")
	return &matrixUsecase{
		nodeAPIURL: nodeURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (u *matrixUsecase) fetchStatsFromNode(token string, matrices ...[][]float64) (*StatsResponse, error) {
	reqBody := StatsRequest{
		Matrices: matrices,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/stats", u.nodeAPIURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al comunicarse con el microservicio de estadisticas: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("el microservicio de estadisticas retorno un codigo de estado no exitoso: %d", resp.StatusCode)
	}

	var stats StatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (u *matrixUsecase) ProcessRotation(token string, matrix [][]float64) ([][]float64, *StatsResponse, error) {
	rotated, err := domain.Rotate90CW(matrix)
	if err != nil {
		return nil, nil, err
	}

	stats, err := u.fetchStatsFromNode(token, rotated)
	if err != nil {
		return nil, nil, err
	}

	return rotated, stats, nil
}

func (u *matrixUsecase) ProcessQR(token string, matrix [][]float64) ([][]float64, [][]float64, bool, *StatsResponse, error) {
	q, r, err := domain.QRDecomposition(matrix)
	if err != nil {
		return nil, nil, false, nil, err
	}

	verified := domain.VerifyQR(matrix, q, r)
	if !verified {
		return nil, nil, false, nil, fmt.Errorf("la verificacion matematica automatica de QR fallo")
	}

	stats, err := u.fetchStatsFromNode(token, q, r)
	if err != nil {
		return nil, nil, false, nil, err
	}

	return q, r, verified, stats, nil
}
