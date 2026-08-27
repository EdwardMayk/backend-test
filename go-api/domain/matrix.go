package domain

import (
	"errors"
	"math"

	"gonum.org/v1/gonum/mat"
)

// Epsilon define el umbral de tolerancia para comparaciones de punto flotante.
const Epsilon = 1e-7

// Rotate90CW rota una matriz rectangular de dimensiones M x N 90 grados en sentido horario,
// resultando en una matriz de dimensiones N x M. Retorna un error si la matriz esta vacia
// o no es rectangular.
func Rotate90CW(matrix [][]float64) ([][]float64, error) {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return nil, errors.New("la matriz no puede estar vacia")
	}

	rows := len(matrix)
	cols := len(matrix[0])

	for _, row := range matrix {
		if len(row) != cols {
			return nil, errors.New("la matriz debe ser rectangular")
		}
	}

	rotated := make([][]float64, cols)
	for i := range rotated {
		rotated[i] = make([]float64, rows)
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			rotated[c][rows-1-r] = matrix[r][c]
		}
	}

	return rotated, nil
}

// QRDecomposition calcula la factorizacion QR (A = Q * R) para una matriz de M x N
// utilizando reflexiones de Householder a traves de Gonum. Requiere que M >= N.
func QRDecomposition(matrix [][]float64) ([][]float64, [][]float64, error) {
	rows := len(matrix)
	if rows == 0 {
		return nil, nil, errors.New("la matriz no puede estar vacia")
	}
	cols := len(matrix[0])
	if cols == 0 {
		return nil, nil, errors.New("la matriz no puede estar vacia")
	}

	if rows < cols {
		return nil, nil, errors.New("la descomposicion QR requiere que el numero de filas sea mayor o igual al numero de columnas (M >= N)")
	}

	flatData := make([]float64, 0, rows*cols)
	for _, row := range matrix {
		if len(row) != cols {
			return nil, nil, errors.New("la matriz debe ser rectangular")
		}
		flatData = append(flatData, row...)
	}

	A := mat.NewDense(rows, cols, flatData)

	var qr mat.QR
	qr.Factorize(A)

	var Q, R mat.Dense
	qr.QTo(&Q)
	qr.RTo(&R)

	qRows, qCols := Q.Dims()
	rRows, rCols := R.Dims()

	qResult := make([][]float64, qRows)
	for i := 0; i < qRows; i++ {
		qResult[i] = make([]float64, qCols)
		for j := 0; j < qCols; j++ {
			qResult[i][j] = Q.At(i, j)
		}
	}

	rResult := make([][]float64, rRows)
	for i := 0; i < rRows; i++ {
		rResult[i] = make([]float64, rCols)
		for j := 0; j < rCols; j++ {
			rResult[i][j] = R.At(i, j)
		}
	}

	return qResult, rResult, nil
}

// VerifyQR audita numericamente la estabilidad y validez de la descomposicion QR.
// Verifica que el producto Q * R aproxime A, y que Q^T * Q aproxime la matriz identidad,
// todo bajo la tolerancia definida por Epsilon.
func VerifyQR(A, Q, R [][]float64) bool {
	qRows := len(Q)
	if qRows == 0 {
		return false
	}
	qCols := len(Q[0])
	rRows := len(R)
	if rRows == 0 {
		return false
	}
	rCols := len(R[0])

	for i := 0; i < qRows; i++ {
		for j := 0; j < rCols; j++ {
			var sum float64
			for k := 0; k < qCols; k++ {
				sum += Q[i][k] * R[k][j]
			}
			if math.Abs(sum-A[i][j]) > Epsilon {
				return false
			}
		}
	}

	for i := 0; i < qCols; i++ {
		for j := 0; j < qCols; j++ {
			var sum float64
			for k := 0; k < qRows; k++ {
				sum += Q[k][i] * Q[k][j]
			}
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if math.Abs(sum-expected) > Epsilon {
				return false
			}
		}
	}

	return true
}
