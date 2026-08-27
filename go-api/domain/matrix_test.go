package domain

import (
	"testing"
)

func TestRotate90CW(t *testing.T) {
	input := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}

	expected := [][]float64{
		{4, 1},
		{5, 2},
		{6, 3},
	}

	result, err := Rotate90CW(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(expected) || len(result[0]) != len(expected[0]) {
		t.Fatalf("matrix dimensions mismatch: got %dx%d, want %dx%d", len(result), len(result[0]), len(expected), len(expected[0]))
	}

	for i := range expected {
		for j := range expected[i] {
			if result[i][j] != expected[i][j] {
				t.Errorf("mismatch at [%d][%d]: got %f, want %f", i, j, result[i][j], expected[i][j])
			}
		}
	}
}

func TestQRDecompositionAndVerification(t *testing.T) {
	A := [][]float64{
		{12, -51, 4},
		{6, 167, -68},
		{-4, 24, -41},
	}

	Q, R, err := QRDecomposition(A)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	verified := VerifyQR(A, Q, R)
	if !verified {
		t.Error("QR decomposition could not be verified within epsilon tolerance")
	}
}
