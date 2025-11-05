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

The system follows a microservice architecture pattern where each service is independently deployable and scalable. Services communicate through REST APIs and are containerized using Docker for consistent deployment across environments.

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   User      │     │  Merchant   │     │  Product    │
│  Service    │     │  Service    │     │  Service    │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                    │
       └───────────────────┼────────────────────┘
                           │
                    ┌──────┴──────┐
                    │   API       │
                    │   Gateway   │
                    └──────┬──────┘
                           │
       ┌───────────────────┼────────────────────┐
       │                   │                    │
┌──────┴──────┐     ┌──────┴──────┐     ┌──────┴──────┐
│Transaction  │     │ Notification│     │  Warehouse  │
│  Service    │     │  Service    │     │  Service    │
└─────────────┘     └─────────────┘     └─────────────┘
```

## 🔧 Services

### Currently Implemented

- **User Service**: Manages user authentication, authorization, and profile management
- **Merchant Service**: Handles merchant registration, verification, and management
- **Product Service**: Manages product catalog, inventory, and product information
- **Transaction Service**: Processes orders, payments, and transaction history
- **Warehouse Service**: Manages warehouse operations, stock movements, and logistics
- **Notification Service**: Handles email, SMS, and push notifications across the system

### Planned Services

- **Analytics Service**: Business intelligence and reporting
- **Shipping Service**: Integration with shipping providers
- **Payment Gateway Service**: Multiple payment method integration

## 💻 Tech Stack

### Core Technologies

- **Language**: Go 1.21+
- **Framework**: [Fiber](https://gofiber.io/) - Fast, minimalist web framework
- **ORM**: [GORM](https://gorm.io/) - Developer-friendly ORM for Go
- **Database**: PostgreSQL 15+
- **Storage**: [Supabase Storage](https://supabase.com/storage) - For image and file uploads
- **Containerization**: Docker & Docker Compose

### Additional Tools

- **Migration**: GORM Auto Migration / golang-migrate
- **Validation**: go-playground/validator
- **Configuration**: viper / cobra / godotenv
- **Logging**: fiber/v2/log / zerolog /

## 📦 Prerequisites

Before running this project, ensure you have the following installed:

- Go 1.21 or higher
- Docker & Docker Compose
- PostgreSQL 15+ (if running locally without Docker)
- Git
- Make (optional, for Makefile commands)

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/devbydenis/warehouse-microservice.git
cd warehouse-microservice
```

### 2. Environment Configuration

Create `.env` file in each service directory:

```bash
APP_ENV="development"
APP_PORT=8081

DATABASE_PORT=5432
DATABASE_HOST=localhost
DATABASE_USER=postgres
DATABASE_PASSWORD=
DATABASE_NAME=warehouse_user_db
DATABASE_MAX_OPEN_CONNECTION=100
DATABASE_MAX_IDLE_CONNECTION=20

RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=guest
RABBITMQ_PASSWORD=guest

REDIS_HOST=warehouse_redis
REDIS_PORT=6379

SUPABASE_URL=
SUPABASE_KEY=
SUPABASE_BUCKET=
```

### 3. Run with Docker Compose

```bash
docker-compose up -d
```

### 4. Run Individual Service (Development)

```bash
cd services/user-service
go mod download
go run main.go
```

### 5. Database Migration

```bash
# Run migrations for all services
make migrate-up

# Or manually for each service
cd services/user-service
go run main.go migrate
```

## 📁 Project Structure

```
.
├── api-gateway
├── merchant-service
├── notification-service
├── product-service
├── transaction-service
├── user-service
│   ├── app
│   │   ├── app.go
│   │   ├── container.go
│   │   └── routes.go
│   ├── cmd
│   │   ├── root.go
│   │   └── start.go
│   ├── configs
│   │   └── config.go
│   ├── controller
│   │   ├── auth_controller.go
│   │   ├── request
│   │   │   ├── auth_request.go
│   │   │   ├── role_request.go
│   │   │   └── user_request.go
│   │   ├── response
│   │   │   ├── auth_response.go
│   │   │   ├── role_response.go
│   │   │   ├── upload_response.go
│   │   │   └── user_response.go
│   │   ├── role_controller.go
│   │   ├── upload_controller.go
│   │   └── user_controller.go
│   ├── database
│   │   ├── manager_seeder.go
│   │   ├── postgres_database.go
│   │   └── role_seeder.go
│   ├── model
│   │   ├── role_model.go
│   │   ├── user_model.go
│   │   └── user_role_model.go
│   ├── pkg
│   │   ├── conv
│   │   │   └── conv.go
│   │   ├── pagination
│   │   │   └── pagination.go
│   │   ├── storage
│   │   │   ├── file_upload_helper.go
│   │   │   └── supabase_storage.go
│   │   └── validator
│   │       └── request_validator.go
│   ├── repository
│   │   ├── role_repository.go
│   │   └── user_repository.go
│   ├── service
│   │   └── rabbitmq_service.go
│   └── usecase
│       ├── role_usecase.go
│       └── user_usecase.go
│   ├── .env
│   ├── .gitignore
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
└── warehouse-service
├── docker-compose.yml
```

## 🛠️ Development

### Running Tests

```bash
# Run all tests
make test

# Run tests for specific service
cd services/user-service
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Code Formatting

```bash
# Format all Go files
make fmt

# Or manually
go fmt ./...
```

### Linting

```bash
# Run linter
make lint

# Or using golangci-lint
golangci-lint run
```

## 📚 API Documentation

API documentation is available for each service:

- User Service: `http://localhost:8001/swagger`
- Merchant Service: `http://localhost:8002/swagger`
- Product Service: `http://localhost:8003/swagger`
- Transaction Service: `http://localhost:8004/swagger`
- Warehouse Service: `http://localhost:8005/swagger`
- Notification Service: `http://localhost:8006/swagger`

### Example Endpoints

#### User Service
```
POST      /api/v1/auth/login

POST      /api/v1/users
GET       /api/v1/users
GET       /api/v1/users/:id
PUT       /api/v1/users/:id
DELETE    /api/v1/users/:id

POST      /api/v1/roles
GET       /api/v1/roles
GET       /api/v1/roles/:id
PUT       /api/v1/roles/:id
DELETE    /api/v1/roles/:id

POST      /api/v1/assign-role
GET       /api/v1/assign-role
GET       /api/v1/assign-role/:id
PUT       /api/v1/assign-role/:id
GET       /api/v1/assign-role/role/:roleName

POST      /api/v1/upload/photo
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

## 🔐 Security

- JWT-based authentication
- Password hashing using bcrypt
- Input validation and sanitization
- Rate limiting on API endpoints
- CORS configuration
- SQL injection prevention via GORM

## 📈 Monitoring & Logging

- Structured logging with fiber/log / zerolog
- Request/Response logging middleware
- Error tracking and reporting
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
- Email: hello.denisrahmadi@example.com
- GitHub Issues: [Project Issues](https://github.com/devbydenis/warehouse-microservice/issues)

---

**Status**: 🚧 Under Active Development

Last Updated: October 2025