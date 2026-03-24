# Go WX

A weather forecast application with a Go (Gin) API backend and a Next.js/React frontend. Click anywhere on the map to drop a pin and see the current forecast for that location, powered by the [National Weather Service API](https://www.weather.gov/documentation/services-web-api).

## Project Structure

```
go-wx/
├── backend/           # Go API server
│   ├── cmd/server/    # Entry point
│   ├── internal/
│   │   ├── api/       # Router and middleware
│   │   ├── handlers/  # HTTP handlers
│   │   ├── models/    # Request/response structs and NWS API models
│   │   └── services/  # Weather service (NWS API client)
│   ├── Dockerfile
│   └── go.mod
├── frontend/          # Next.js React app
│   ├── src/
│   │   ├── app/       # Next.js app router pages
│   │   ├── components/# Map component (Leaflet)
│   │   └── lib/       # Shared utilities
│   ├── Dockerfile
│   └── package.json
└── docker-compose.yml
```

## Prerequisites

### For Docker (recommended)

- [Docker](https://docs.docker.com/get-docker/) (v20+)
- [Docker Compose](https://docs.docker.com/compose/install/) (v2+)

### For running locally without Docker

- [Go](https://go.dev/dl/) (v1.23+)
- [Node.js](https://nodejs.org/) (v22+)
- npm (included with Node.js)

## Getting Started

### Option 1: Docker (recommended)

This is the quickest way to get both services running. No Go or Node.js installation required.

```bash
docker compose up --build
```

Once both containers are running:

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

To stop the containers:

```bash
docker compose down
```

### Option 2: Run locally without Docker

#### 1. Set up the backend

```bash
cd backend
go mod tidy
go run ./cmd/server
```

The API server will start on http://localhost:8080.

#### 2. Set up the frontend (in a separate terminal)

```bash
cd frontend
npm install
npm run dev
```

The frontend dev server will start on http://localhost:3000.

## Building for Production (without Docker)

### Backend

```bash
cd backend
go mod tidy
go build -o server ./cmd/server
./server
```

### Frontend

```bash
cd frontend
npm install
npm run build
npm start
```

## Running Tests

### Backend

```bash
cd backend
go test ./...
```

### Frontend

```bash
cd frontend
npm install
npm test
```

## API Endpoints

| Method | Path            | Description                       |
|--------|-----------------|-----------------------------------|
| GET    | /api/health     | Health check                      |
| GET    | /api/forecast   | Get forecast for a location       |

### GET /api/forecast

Query parameters:

| Parameter | Type  | Required | Description       |
|-----------|-------|----------|-------------------|
| lat       | float | yes      | Latitude          |
| lon       | float | yes      | Longitude         |

Example:

```bash
curl "http://localhost:8080/api/forecast?lat=39.7456&lon=-97.0892"
```

Response:

```json
{
  "shortForecast": "Partly Cloudy",
  "temperature": 68,
  "temperatureUnit": "F",
  "temperatureCharacterization": "moderate"
}
```

Temperature characterization values:

- **hot**: above 75°F
- **moderate**: 50°F to 75°F
- **cold**: below 50°F

## Environment Variables

| Variable              | Default                  | File                  | Description                          |
|-----------------------|--------------------------|-----------------------|--------------------------------------|
| PORT                  | 8080                     | backend/.env          | Backend server port                  |
| NEXT_PUBLIC_API_URL   | http://localhost:8080     | frontend/.env.local   | Backend API URL (used by frontend)   |

Default `.env` files are included in the repository. To override values locally, edit:

- `backend/.env` — loaded by the Go server via [godotenv](https://github.com/joho/godotenv)
- `frontend/.env.local` — loaded automatically by Next.js (not committed to git by convention)
