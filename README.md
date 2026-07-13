# Warehouse Microservice System

A scalable microservice-based warehouse management system built with Go, designed to handle modern e-commerce operations with high performance and reliability.

## 🚀 Overview

This project implements a distributed microservice architecture for warehouse management, providing comprehensive solutions for inventory tracking, order fulfillment, merchant management, and user operations.

## 📋 Table of Contents

- [Architecture](#architecture)
- [Services](#services)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Development](#development)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)
- [License](#license)

## 🏗️ Architecture

The system follows a microservice architecture pattern where each service is independently deployable and scalable. All client requests are routed through the **API Gateway**, which handles authentication, rate limiting, and reverse proxying to downstream services.

```
                         ┌─────────────────────┐
                         │      Client          │
                         └──────────┬──────────┘
                                    │
                         ┌──────────▼──────────┐
                         │    API Gateway       │
                         │      :8080           │
                         │  (JWT Auth, Rate     │
                         │   Limit, Proxy)      │
                         └──────────┬──────────┘
                                    │
        ┌───────────┬───────────────┼───────────────┬───────────┐
        │           │               │               │           │
┌───────▼──────┐ ┌──▼───────┐ ┌────▼─────┐ ┌──────▼────┐ ┌────▼────────┐
│   User       │ │ Product  │ │ Merchant │ │ Warehouse │ │ Transaction │
│  Service     │ │ Service  │ │ Service  │ │  Service  │ │   Service   │
│   :8081      │ │  :8082   │ │  :8084   │ │   :8083   │ │    :8085    │
└──────────────┘ └──────────┘ └──────────┘ └───────────┘ └─────────────┘
                                    │
                         ┌──────────▼──────────┐
                         │  Notification        │
                         │    Service           │
                         │      :8086           │
                         └─────────────────────┘
```

## 🔧 Services

### Currently Implemented

| Service | Port | Description |
|---------|------|-------------|
| **API Gateway** | 8080 | Single entry point — JWT auth, rate limiting, reverse proxy, aggregated Swagger |
| **User Service** | 8081 | User authentication, authorization, role management, and profile |
| **Product Service** | 8082 | Product catalog, categories, barcode lookup, and image uploads |
| **Warehouse Service** | 8083 | Warehouse operations, stock management, and logistics |
| **Merchant Service** | 8084 | Merchant registration, management, and merchant product catalog |
| **Transaction Service** | 8085 | Orders, payments via Midtrans, and dashboard reporting |
| **Notification Service** | 8086 | Email notifications via RabbitMQ consumer and direct HTTP |

### Planned Services

- **Analytics Service**: Business intelligence and reporting
- **Shipping Service**: Integration with shipping providers
- **Payment Gateway Service**: Multiple payment method integration

## 💻 Tech Stack

### Core Technologies

- **Language**: Go 1.25+
- **Framework**: [Fiber](https://gofiber.io/) - Fast, minimalist web framework
- **ORM**: [GORM](https://gorm.io/) - Developer-friendly ORM for Go
- **Database**: PostgreSQL 15+
- **Storage**: [Supabase Storage](https://supabase.com/storage) - For image and file uploads
- **Containerization**: Docker & Docker Compose
- **Message Broker**: [RabbitMQ](https://www.rabbitmq.com/) - For async inter-service communication
- **Cache**: [Redis](https://redis.io/) - For caching and rate limiting
- **Payment Gateway**: [Midtrans](https://midtrans.com/) - For payment processing

### Additional Tools

- **API Documentation**: Swagger (swaggo/fiber-swagger) — per service + aggregated at API Gateway
- **Validation**: go-playground/validator
- **Configuration**: viper / cobra / godotenv
- **Logging**: fiber/v2/log / zerolog

## 📦 Prerequisites

Before running this project, ensure you have the following installed:

- Go 1.25 or higher
- Docker & Docker Compose
- PostgreSQL 15+ (if running locally without Docker)
- Git
- Make (optional, for Swagger generation commands)
- [swag CLI](https://github.com/swaggo/swag) (optional, for regenerating Swagger docs)

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/devbydenis/warehouse-microservice.git
cd warehouse-microservice
```

### 2. Environment Configuration

Each service has its own `.env.example` file. Copy and fill it for each service:

```bash
# Repeat for each service directory
cp user-service/.env.example user-service/.env
cp product-service/.env.example product-service/.env
cp warehouse-service/.env.example warehouse-service/.env
cp merchant-service/.env.example merchant-service/.env
cp transaction-service/.env.example transaction-service/.env
cp notification-service/.env.example notification-service/.env
cp api-gateway/.env.example api-gateway/.env
```

Environment variables vary per service. Below is a summary of key variables per service:

**user-service / product-service / warehouse-service**
```bash
APP_ENV=development
APP_PORT=808x

DATABASE_HOST=localhost
DATABASE_PORT=543x
DATABASE_USER=postgres
DATABASE_PASSWORD=
DATABASE_NAME=warehouse_xxx_db
DATABASE_MAX_OPEN_CONNECTION=100
DATABASE_MAX_IDLE_CONNECTION=20

RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=guest
RABBITMQ_PASSWORD=guest

REDIS_HOST=localhost
REDIS_PORT=6379

SUPABASE_URL=
SUPABASE_KEY=
SUPABASE_BUCKET=

JWT_SECRET_KEY=
JWT_ISSUER=
JWT_DURATION=

URL_API_GATEWAY=http://localhost:8080
```

**merchant-service** (additional inter-service URLs)
```bash
# ...same as above, plus:
URL_USER_SERVICE=http://localhost:8081
URL_PRODUCT_SERVICE=http://localhost:8082
URL_WAREHOUSE_SERVICE=http://localhost:8083
```

**transaction-service** (additional Midtrans config)
```bash
# ...base config, plus:
MIDTRANS_SERVER_KEY=
MIDTRANS_CLIENT_KEY=
MIDTRANS_MERCHANT_ID=
MIDTRANS_IS_PRODUCTION=false

JWT_SECRET_KEY=
JWT_ISSUER=
JWT_DURATION=

URL_API_GATEWAY=http://localhost:8080
```

**notification-service**
```bash
APP_ENV=development
APP_PORT=8086

DATABASE_HOST=localhost
DATABASE_PORT=5436
DATABASE_USER=postgres
DATABASE_PASSWORD=
DATABASE_NAME=warehouse_notification_db

RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=guest
RABBITMQ_PASSWORD=guest

EMAIL_HOST=
EMAIL_PORT=
EMAIL_USER=
EMAIL_PASSWORD=
EMAIL_FROM=
```

**api-gateway**
```bash
APP_ENV=development
APP_PORT=8080

USER_SERVICE_URL=http://localhost:8081
PRODUCT_SERVICE_URL=http://localhost:8082
WAREHOUSE_SERVICE_URL=http://localhost:8083
MERCHANT_SERVICE_URL=http://localhost:8084
TRANSACTION_SERVICE_URL=http://localhost:8085
NOTIFICATION_SERVICE_URL=http://localhost:8086

JWT_SECRET_KEY=
JWT_ISSUER=
JWT_DURATION=

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=
```

### 3. Run with Docker Compose

```bash
# Start all infrastructure and application services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### 4. Run Individual Service (Development)

```bash
cd user-service
go mod download
go run main.go start
```

### 5. Database Migration

Migration runs automatically via GORM AutoMigrate on service startup. No manual migration step is required.

## 📁 Project Structure

```
micro-warehouse/
├── docker-compose.yml
├── Makefile
├── README.md
├── AGENTS.md
│
├── api-gateway/               # API Gateway — entry point for all clients
│   ├── configs/               # JWT & Redis config loaders
│   ├── controller/            # Auth controller (login proxy)
│   ├── docs/                  # Swagger stub (aggregated from services)
│   ├── internal/swagger/      # Swagger aggregator logic
│   ├── middleware/            # JWT auth, Redis rate limiter
│   ├── main.go                # App entry point + route setup
│   ├── Dockerfile
│   └── .env.example
│
├── user-service/              # Example of standard service layout
│   ├── app/
│   │   ├── app.go
│   │   ├── container.go
│   │   └── routes.go
│   ├── cmd/
│   │   ├── root.go
│   │   └── start.go
│   ├── configs/
│   │   └── config.go
│   ├── controller/
│   │   ├── request/
│   │   ├── response/
│   │   ├── auth_controller.go
│   │   ├── role_controller.go
│   │   ├── upload_controller.go
│   │   └── user_controller.go
│   ├── database/
│   │   ├── postgres_database.go
│   │   ├── role_seeder.go
│   │   └── manager_seeder.go
│   ├── docs/                  # Swagger generated files (docs.go, swagger.json, swagger.yaml)
│   ├── middleware/            # gateway_auth.go
│   ├── model/
│   ├── pkg/
│   │   ├── conv/
│   │   ├── pagination/
│   │   ├── storage/
│   │   └── validator/
│   ├── repository/
│   ├── service/               # rabbitmq_service.go
│   ├── usecase/
│   ├── main.go
│   ├── Dockerfile
│   └── .env.example
│
├── product-service/           # Same structure as user-service
├── warehouse-service/         # Same structure as user-service
├── merchant-service/          # Same structure as user-service
├── transaction-service/       # Same structure as user-service
└── notification-service/      # Simplified — no DB model layer, RabbitMQ consumer
```

> All services follow the same layered architecture: `controller → usecase → repository → model`. See `AGENTS.md` for detailed architectural documentation.

## 🛠️ Development

### Generate Swagger Docs

```bash
# Generate for all services at once
make swag

# Or per service
make swag-user
make swag-product
make swag-merchant
make swag-warehouse
make swag-transaction
```

Requires [swag CLI](https://github.com/swaggo/swag):
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Running Tests

```bash
# Run tests for a specific service
cd user-service
go test ./...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
golangci-lint run
```

## 📚 API Documentation

### Aggregated Swagger (via API Gateway)

All service specs are aggregated and accessible at a single URL:

- **Aggregated Docs**: `http://localhost:8080/swagger`

### Per-Service Swagger

Each service also exposes its own Swagger UI:

| Service | Swagger URL |
|---------|-------------|
| API Gateway | `http://localhost:8080/swagger` |
| User Service | `http://localhost:8081/swagger` |
| Product Service | `http://localhost:8082/swagger` |
| Warehouse Service | `http://localhost:8083/swagger` |
| Merchant Service | `http://localhost:8084/swagger` |
| Transaction Service | `http://localhost:8085/swagger` |
| Notification Service | `http://localhost:8086/swagger` |

### Endpoints

All protected endpoints require a JWT Bearer token obtained from `POST /api/v1/auth/login`.

#### API Gateway
```
GET  /health                          # Gateway health check + service URLs
GET  /swagger                         # Aggregated Swagger UI
GET  /swagger/aggregated.json         # Aggregated OpenAPI spec
```

#### User Service
```
POST   /api/v1/auth/login

POST   /api/v1/users
GET    /api/v1/users
GET    /api/v1/users/:id
GET    /api/v1/users/role/:roleName
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id

POST   /api/v1/roles
GET    /api/v1/roles
GET    /api/v1/roles/:id
PUT    /api/v1/roles/:id
DELETE /api/v1/roles/:id

POST   /api/v1/assign-role
GET    /api/v1/assign-role
GET    /api/v1/assign-role/:userRoleID
PUT    /api/v1/assign-role/:userRoleID

POST   /api/v1/upload/photo
```

#### Product Service
```
POST   /api/v1/categories
GET    /api/v1/categories
GET    /api/v1/categories/:id
PUT    /api/v1/categories/:id
DELETE /api/v1/categories/:id

POST   /api/v1/products
GET    /api/v1/products
GET    /api/v1/products/:id
GET    /api/v1/products/barcode/:barcode
PUT    /api/v1/products/:id
DELETE /api/v1/products/:id
POST   /api/v1/products/:id/upload-image
```

#### Merchant Service
```
POST   /api/v1/merchants
GET    /api/v1/merchants
GET    /api/v1/merchants/:id
PUT    /api/v1/merchants/:id
DELETE /api/v1/merchants/:id

POST   /api/v1/merchant-products
GET    /api/v1/merchant-products
GET    /api/v1/merchant-products/:merchant_product_id
GET    /api/v1/merchant-products/barcode/:barcode
PUT    /api/v1/merchant-products/:merchant_product_id
DELETE /api/v1/merchant-products/:merchant_product_id
DELETE /api/v1/merchant-products/product/:product_id
GET    /api/v1/merchant-products/:product_id/total-stock

POST   /api/v1/upload-merchant
```

#### Warehouse Service
```
POST   /api/v1/warehouses
GET    /api/v1/warehouses
GET    /api/v1/warehouses/:id
PUT    /api/v1/warehouses/:id
DELETE /api/v1/warehouses/:id

POST   /api/v1/warehouse-products/:warehouse_id
GET    /api/v1/warehouse-products/:warehouse_id
GET    /api/v1/warehouse-products/:warehouse_id/detail/:product_id
PUT    /api/v1/warehouse-products/detail/:warehouse_product_id/:warehouse_id
DELETE /api/v1/warehouse-products/detail/:warehouse_product_id
DELETE /api/v1/warehouse-products/detail/products/:product_id
GET    /api/v1/warehouse-products/detail/products/:product_id/total-stock
GET    /api/v1/warehouse-products/detail/products/:product_id
GET    /api/v1/warehouse-products/detail/products/:warehouse_product_id/warehouses

POST   /api/v1/upload-warehouse
```

#### Transaction Service
```
POST   /api/v1/transactions
GET    /api/v1/transactions

GET    /api/v1/dashboard/manager
GET    /api/v1/dashboard/keeper/merchant/:merchant_id

POST   /api/v1/midtrans/callback     # Public — Midtrans webhook, no JWT required
```

#### Notification Service
```
POST   /send-email                    # Direct HTTP trigger
POST   /send-welcome-email            # Direct HTTP trigger
```

> Notification service is primarily event-driven via RabbitMQ. HTTP endpoints are available for direct triggers but are intended for internal use.

## 🔐 Security

- JWT-based authentication enforced at the API Gateway
- Redis-based rate limiting: global, per-auth-endpoint, and per-API-route
- Password hashing using bcrypt
- Input validation via go-playground/validator
- CORS configuration
- SQL injection prevention via GORM parameterized queries
- Internal requests validated via `X-Internal-Request` header (gateway auth middleware)

## 📈 Monitoring & Logging

- Structured logging with zerolog
- Request/Response logging middleware (Fiber logger)
- Error tracking with step-numbered log format: `[Controller] Method - step: error`
- Performance metrics (planned)

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Coding Standards

- Follow Go best practices and idioms
- Write unit tests for new features
- Update documentation as needed
- Use meaningful commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Denis Rahmadi** - *Backend Engineer* - [devbydenis](https://github.com/devbydenis)

## 🙏 Acknowledgments

- Fiber framework community
- GORM contributors
- Supabase team
- Go community

## 📞 Contact

For questions or support, please contact:
- Email: hello.denisrahmadi@gmail.com
- GitHub Issues: [Project Issues](https://github.com/devbydenis/warehouse-microservice/issues)

---

**Status**: 🚧 Under Active Development

Last Updated: July 2026
