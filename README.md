# Netflix Movies REST API

A production-ready REST API for Netflix Movies built with Go, Gin, GORM, and PostgreSQL.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.27+ |
| Web Framework | Gin |
| JSON | encoding/json/v2 |
| Database | PostgreSQL |
| ORM | GORM v2 |
| Migrations | Goose v3 |
| Configuration | Viper |
| Logging | Slog |
| Testing | Testify, Mockery |
| Containerization | Docker |
| Orchestration | Kubernetes |

## Configuration

The application uses **Viper** for configuration management. `godotenv` is **not** used.

### Precedence (Highest to Lowest)
1.  **Environment Variables** (e.g., `APP_SERVER_PORT=9090`)
2.  **Native .env file** (loaded by Viper)
3.  **Environment-specific Config** (`config/config.test.yaml` or `config/config.prod.yaml`)
4.  **Default Config** (`config/config.yaml` - optional)
5.  **Hardcoded Defaults**

### Local Development
For local development, you can create a `.env` file in the root directory. This file is git-ignored.
```bash
cp .env.example .env
```
Modify `.env` to set your local environment variables.

**Note:** Environment variables (e.g., set via shell export or Docker) will always override values in `.env` or YAML files.

## Setup

### Prerequisites

- Go 1.27+
- Docker & Docker Compose
- Make

### Local Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/amit/gin-crud-demo.git
   cd gin-crud-demo
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Start the database:
   ```bash
   make docker-up
   ```

4. Run migrations:
   ```bash
   make migrate-up
   ```

5. Run the application:
   ```bash
   make run
   ```

### Docker Setup

1. Build the Docker image:
   ```bash
   make docker-build
   ```

2. Run with Docker Compose:
   ```bash
   make docker-up
   ```

## Testing

Run all tests:
```bash
make test
```

Generate mocks:
```bash
make mocks
```

## API Endpoints

### Health Check
- `GET /api/v1/healthz`

### Movies
- `POST /api/v1/movies` - Create a movie
- `GET /api/v1/movies/:id` - Get a movie by ID
- `PUT /api/v1/movies/:id` - Update a movie
- `DELETE /api/v1/movies/:id` - Delete a movie
- `GET /api/v1/movies` - List movies (pagination: ?page=1&limit=10)

### TV Shows
- `POST /api/v1/tv-shows`
- `GET /api/v1/tv-shows/:id`
- `PUT /api/v1/tv-shows/:id`
- `DELETE /api/v1/tv-shows/:id`
- `GET /api/v1/tv-shows`

### Seasons
- `POST /api/v1/seasons`
- `GET /api/v1/seasons/:id`
- `PUT /api/v1/seasons/:id`
- `DELETE /api/v1/seasons/:id`
- `GET /api/v1/seasons`

### View Summaries
- `POST /api/v1/view-summaries`
- `GET /api/v1/view-summaries/:id`
- `PUT /api/v1/view-summaries/:id`
- `DELETE /api/v1/view-summaries/:id`
- `GET /api/v1/view-summaries`

## Curl Examples

### Create Movie
```bash
curl -X POST http://localhost:8080/api/v1/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Inception",
    "original_title": "Inception",
    "runtime": 148,
    "release_date": "2010-07-16",
    "available_globally": true,
    "locale": "en-US"
  }'
```

### Get Movie
```bash
curl http://localhost:8080/api/v1/movies/1
```

### List Movies
```bash
curl "http://localhost:8080/api/v1/movies?page=1&limit=5"
```

## Adding a New Model

1. Create the GORM model in `internal/models/`.
2. Create a migration file using `make migrate-create`.
3. Define the repository interface in `internal/repository/interfaces.go`.
4. Implement the repository in `internal/repository/postgres.go`.
5. Generate mocks using `make mocks`.
6. Create the service in `internal/service/` and add methods.
7. Create the handler in `internal/api/handler/`.
8. Register the routes in `internal/api/router.go`.
