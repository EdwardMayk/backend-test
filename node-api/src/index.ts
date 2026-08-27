import express, { Request, Response, NextFunction } from 'express';
import cors from 'cors';
import { calculateStats } from './domain/stats';
import { verifyJWT } from './domain/auth';

if (!process.env.JWT_SECRET) {
  console.error("ERROR FATAL: la variable de entorno JWT_SECRET no esta configurada.");
  process.exit(1);
}

const app = express();
app.use(cors());
app.use(express.json());

const PORT = process.env.PORT || 3000;

const authMiddleware = (req: Request, res: Response, next: NextFunction) => {
  const authHeader = req.headers.authorization;
  if (!authHeader) {
    return res.status(401).json({ error: 'falta el token de autorizacion' });
  }

  const parts = authHeader.split(' ');
  if (parts.length !== 2 || parts[0].toLowerCase() !== 'bearer') {
    return res.status(401).json({ error: 'formato de cabecera de autorizacion invalido' });
  }

  try {
    const payload = verifyJWT(parts[1]);
    (req as any).user = payload;
    return next();
  } catch (error: any) {
    return res.status(401).json({ error: 'token invalido o expirado: ' + error.message });
  }
};

app.post('/api/v1/stats', authMiddleware, (req: Request, res: Response) => {
  const { matrices } = req.body;

  if (!matrices || !Array.isArray(matrices) || matrices.length === 0) {
    return res.status(400).json({ error: 'payload invalido: el campo "matrices" es requerido y debe ser un array' });
  }

  try {
    const stats = calculateStats(matrices);
    return res.status(200).json(stats);
  } catch (error: any) {
    return res.status(500).json({ error: error.message });
  }
});

app.listen(PORT, () => {
  console.log(`Node.js Express stats service is running on port ${PORT}`);
});
