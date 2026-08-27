import { calculateStats, isDiagonalMatrix } from './stats';

describe('Calculos Estadisticos y de Matriz Diagonal', () => {
  test('isDiagonalMatrix identifica correctamente matrices diagonales', () => {
    const diagMatrix = [
      [5, 0, 0],
      [0, 12, 0],
      [0, 0, -3]
    ];
    expect(isDiagonalMatrix(diagMatrix)).toBe(true);

    const nonDiagMatrix = [
      [5, 1, 0],
      [0, 12, 0],
      [0, 0, -3]
    ];
    expect(isDiagonalMatrix(nonDiagMatrix)).toBe(false);

    const nearDiagMatrix = [
      [5, 1e-8, 0],
      [0, 12, 0],
      [0, 0, -3]
    ];
    expect(isDiagonalMatrix(nearDiagMatrix)).toBe(true);
  });

  test('calculateStats calcula estadisticas correctamente sobre multiples matrices', () => {
    const matrix1 = [
      [1, 2],
      [3, 4]
    ];
    const matrix2 = [
      [5, 0],
      [0, 6]
    ];

    const stats = calculateStats([matrix1, matrix2]);

    expect(stats.min).toBe(0);
    expect(stats.max).toBe(6);
    expect(stats.sum_total).toBe(21);
    expect(stats.average).toBe(21 / 8);
    expect(stats.is_diagonal).toBe(true);
  });
});
