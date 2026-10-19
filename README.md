# URL Shortener Service

A high-performance URL shortener built with Go, Redis, and AWS Lambda support — part of the Enterpret Backend SDE Intern portfolio.

## 🏗️ Architecture

```
Client → Go HTTP Server → Redis (URL storage + click tracking)
       → Custom alias support
       → Rate limiting
       → 301 redirects
```

## 🚀 Quick Start

### Docker (Recommended)
```bash
docker-compose up --build
```
Shortener available at `http://localhost:8080`

### Local Go
```bash
# Start Redis
docker run -d -p 6379:6379 redis:7-alpine

# Run the service
go run ./cmd/api
```

## 📖 API Endpoints

### Create Short URL
```http
POST /shorten
Content-Type: application/json

{
  "url": "https://example.com/very-long-product-page",
  "custom_alias": "my-product",
  "expires_in": 86400
}
```

**Response:**
```json
{
  "short_url": "http://localhost:8080/my-product",
  "alias": "my-product",
  "expires_at": ""
}
```

### Redirect
```http
GET /my-product
```
Returns 301 redirect to original URL.

## ☁️ Deploy to AWS Lambda

### Build for Lambda (provided.al2023 runtime)
```bash
GOOS=linux GOARCH=arm64 go build -o bootstrap ./cmd/api
zip function.zip bootstrap
```

### Deploy with SAM
```bash
sam build
sam deploy --guided
```

### With ElastiCache Redis
Use AWS ElastiCache for Redis as the backend store.

## 🔧 Make Commands
| Command | Description |
|---------|-------------|
| `make build` | Build Go binary |
| `make run` | Run locally with Redis |
| `make docker-up` | Start with Docker Compose |
| `make test` | Run Go tests |
| `make tidy` | Clean up go.mod |

## 📊 Features
- ✅ Custom alias support
- ✅ Optional TTL expiration
- ✅ Click tracking in Redis
- ✅ Rate limiting (100 req/min per IP)
- ✅ CORS enabled
- ✅ Docker + AWS Lambda ready

---
**Built for Enterpret Backend SDE Intern application.**