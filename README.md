# Sistema de Procesamiento de Matrices

Sistema distribuido para el procesamiento de matrices matemáticas utilizando dos APIs independientes comunicadas por HTTP.

## Arquitectura
* **go-api (Fiber)**: API Orquestadora pública en Go. Procesa rotaciones de matrices, factorización QR, auto-verificación matemática y maneja la seguridad por JWT.
* **node-api (Express)**: Microservicio privado en Node.js. Calcula estadísticas consolidadas y valida propiedades de diagonalidad bajo tolerancia épsilon ($10^{-7}$).
* **frontend (Nginx)**: Cliente web estático para interactuar con la API de Go.

---

## Requisitos y Configuración

### 1. Variables de Entorno
Crea un archivo `.env` en la raíz del proyecto basándote en el archivo de ejemplo:
```bash
cp .env.example .env
```
Es obligatorio configurar estas variables en el archivo `.env`. Las APIs se apagarán al iniciar si las variables no están declaradas.

### 2. Ejecución con Docker
Para compilar y arrancar todo el entorno (Frontend, Go API y Node API):
```bash
docker compose up --build
```
* **Frontend**: http://localhost:80
* **API de Go**: http://localhost:8080

---

## Endpoints de la API (Go - Puerto 8080)

### Autenticación
* **`POST /api/v1/auth/login`**
  * Payload: `{"username": "admin", "password": "password"}` (configurable en `.env`)
  * Retorna: `{"token": "JWT_TOKEN"}`

### Procesamiento (Requieren Cabecera `Authorization: Bearer <JWT_TOKEN>`)
* **`POST /api/v1/matrix/rotate`**
  * Rota una matriz rectangular de $M \times N$ a $N \times M$ (90 grados en sentido horario).
  * Payload: `{"matrix": [[1, 2, 3], [4, 5, 6]]}`
* **`POST /api/v1/matrix/qr`**
  * Realiza la descomposición $A = QR$ para matrices con $M \ge N$ y verifica el resultado.
  * Payload: `{"matrix": [[12, -51, 4], [6, 167, -68], [-4, 24, -41]]}`

---

## Pruebas (Testing)

### API de Go (Unitarias e Integración HTTP)
```bash
cd go-api
go test -v ./...
```

### API de Node.js (Unitarias y Lógica)
```bash
cd node-api
npm install
npm test
```

---

## Despliegue en la Nube
El proyecto es compatible con plataformas como **Railway** o **Render**:
1. Conecta el repositorio de GitHub como un monorepo.
2. Crea tres servicios independientes apuntando a sus respectivos subdirectorios (`go-api`, `node-api`, `frontend`).
3. Configura las variables de entorno (`JWT_SECRET`, `NODE_API_URL`, `API_BASE_URL`, `AUTH_USERNAME`, `AUTH_PASSWORD`) en la interfaz de la plataforma.
