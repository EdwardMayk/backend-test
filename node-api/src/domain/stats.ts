/**
 * Tolerancia epsilon para mitigar micro-desviaciones por representacion 
 * de punto flotante de precision doble bajo la norma IEEE 754.
 */
const EPSILON = 1e-7;

export interface MatrixStats {
  max: number;
  min: number;
  average: number;
  sum_total: number;
  is_diagonal: boolean;
}

/**
 * Valida si una matriz cumple las condiciones estructurales de una matriz diagonal.
 * Una matriz es diagonal si todos sus elementos fuera de la diagonal principal
 * son menores al valor umbral EPSILON (aproximacion a cero).
 * 
 * @param matrix Matriz de numeros a evaluar.
 * @returns boolean true si la matriz es diagonal; false en caso contrario.
 */
export function isDiagonalMatrix(matrix: number[][]): boolean {
  if (matrix.length === 0 || matrix[0].length === 0) {
    return false;
  }

  const rows = matrix.length;
  const cols = matrix[0].length;

  for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
      if (r !== c) {
        if (Math.abs(matrix[r][c]) > EPSILON) {
          return false;
        }
      }
    }
  }
  return true;
}

/**
 * Calcula metricas estadisticas agregadas consolidadas sobre un conjunto de matrices.
 * Determina el valor maximo, el valor minimo, la suma total, el promedio
 * y comprueba si al menos una de las matrices recibidas es diagonal.
 * 
 * @param matrices Coleccion de matrices numericas tridimensionales.
 * @returns MatrixStats Estructura de datos que consolida las estadisticas de las matrices.
 */
export function calculateStats(matrices: number[][][]): MatrixStats {
  if (matrices.length === 0) {
    throw new Error("No se proporcionaron matrices para el calculo de estadisticas");
  }

  let max = -Infinity;
  let min = Infinity;
  let sum_total = 0;
  let elementCount = 0;
  let hasDiagonal = false;

  for (const matrix of matrices) {
    if (matrix.length === 0 || matrix[0].length === 0) {
      continue;
    }

    if (isDiagonalMatrix(matrix)) {
      hasDiagonal = true;
    }

    for (const row of matrix) {
      for (const val of row) {
        if (val > max) max = val;
        if (val < min) min = val;
        sum_total += val;
        elementCount++;
      }
    }
  }

  if (elementCount === 0) {
    return { max: 0, min: 0, average: 0, sum_total: 0, is_diagonal: false };
  }

  return {
    max,
    min,
    average: sum_total / elementCount,
    sum_total,
    is_diagonal: hasDiagonal,
  };
}
