# User Service Documentation

> Dokumentasi lengkap untuk microservice manajemen user, role, dan autentikasi  
> Module: `micro-warehouse/user-service`  
> Versi: 1.0 | Last Updated: 2026-02-14

---

## 📚 Table of Contents

- [1. Executive Summary](#1-executive-summary)
  - [Deskripsi Service](#deskripsi-service)
  - [Tech Stack](#tech-stack)
  - [Fitur Utama](#fitur-utama)
- [2. Quick Start Guide](#2-quick-start-guide)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Environment Setup](#environment-setup)
  - [Running the Service](#running-the-service)
- [3. Project Architecture](#3-project-architecture)
  - [Clean Architecture Layers](#clean-architecture-layers)
  - [Request Flow Diagram](#request-flow-diagram)
- [4. Application Setup](#4-application-setup)
  - [Fiber Configuration](#fiber-configuration)
  - [Middleware Stack](#middleware-stack)
  - [Graceful Shutdown](#graceful-shutdown)
- [5. Configuration Management](#5-configuration-management)
  - [Environment Variables](#environment-variables)
  - [Config Structs](#config-structs)
- [6. Database Layer](#6-database-layer)
  - [ER Diagram](#er-diagram)
  - [Entity Definitions](#entity-definitions)
  - [Migrations & Seeders](#migrations--seeders)
  - [Connection Pooling](#connection-pooling)
- [7. API Documentation](#7-api-documentation)
  - [Authentication](#authentication)
  - [User Management](#user-management)
  - [Role Management](#role-management)
  - [User-Role Assignment](#user-role-assignment)
  - [File Upload](#file-upload)
  - [Request DTOs](#request-dtos)
  - [Response DTOs](#response-dtos)
- [8. Business Logic Deep Dive](#8-business-logic-deep-dive)
  - [User Creation Flow](#user-creation-flow)
  - [Login Flow](#login-flow)
  - [Role Assignment Flow](#role-assignment-flow)
  - [File Upload Flow](#file-upload-flow)
- [9. External Integrations](#9-external-integrations)
  - [RabbitMQ Integration](#rabbitmq-integration)
  - [Supabase Storage](#supabase-storage)
- [10. Utility Packages](#10-utility-packages)
  - [pkg/conv](#pkgconv)
  - [pkg/pagination](#pkgpagination)
  - [pkg/validator](#pkgvalidator)
  - [pkg/storage](#pkgstorage)
- [11. Testing Guide](#11-testing-guide)
  - [Unit Testing](#unit-testing)
  - [Mocking Strategy](#mocking-strategy)
- [12. Deployment](#12-deployment)
  - [Build Command](#build-command)
  - [Environment-specific Considerations](#environment-specific-considerations)
- [13. Troubleshooting](#13-troubleshooting)
  - [Common Issues](#common-issues)
  - [Debug Tips](#debug-tips)

---

## 1. Executive Summary

### Deskripsi Service

**User Service** adalah microservice yang menangani seluruh aspek manajemen user dalam sistem micro-warehouse. Service ini bertanggung jawab untuk:

- **User Management**: CRUD operations untuk user
- **Role Management**: Manajemen role dan permission
- **User-Role Assignment**: Mapping user ke role tertentu
- **Authentication**: Login dengan email dan password
- **File Upload**: Upload foto profil ke cloud storage
- **Event Publishing**: Mengirim event ke RabbitMQ untuk komunikasi antar service

### Tech Stack

| Komponen | Teknologi | Versi | Kegunaan |
|----------|-----------|-------|----------|
| Language | Go | 1.24.7 | Primary language |
| Web Framework | Fiber | v2.52.9 | HTTP server & routing |
| ORM | GORM | v1.31.0 | Database operations |
| Database | PostgreSQL | - | Primary database |
| Message Broker | RabbitMQ | - | Event publishing |
| Storage | Supabase | - | File storage |
| CLI | Cobra | v1.10.1 | Command line interface |
| Config | Viper | v1.21.0 | Configuration management |
| Validation | go-playground/validator | v10.28.0 | Request validation |
| Logging | zerolog | v1.34.0 | Structured logging |
| Crypto | bcrypt | - | Password hashing |

### Fitur Utama

1. **CRUD User** dengan pagination dan search
2. **CRUD Role**
3. **Assign User to Role** (many-to-many relationship)
4. **Authentication** dengan bcrypt password hashing
5. **File Upload** dengan validasi (max 2MB, image only)
6. **Async Email Event** via RabbitMQ saat user created
7. **Graceful Shutdown**
8. **Context Cancellation** handling

---

## 2. Quick Start Guide

### Prerequisites

- Go 1.24.7 atau lebih tinggi
- PostgreSQL 12+
- RabbitMQ 3.8+
- Git

### Installation

```bash
# Clone repository
git clone <repository-url>
cd user-service

# Download dependencies
go mod download

# Verify installation
go mod verify
```

### Environment Setup

Buat file `.env` di root project:

```env
APP_ENV="development"
APP_PORT=8081

# Database Configuration
DATABASE_PORT=5432
DATABASE_HOST=localhost
DATABASE_USER=postgres
DATABASE_PASSWORD=your_password
DATABASE_NAME=warehouse_user_db
DATABASE_MAX_OPEN_CONNECTION=100
DATABASE_MAX_IDLE_CONNECTION=20

# RabbitMQ Configuration
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=guest
RABBITMQ_PASSWORD=guest

# Redis Configuration (reserved for future use)
REDIS_HOST=localhost
REDIS_PORT=6379

# Supabase Storage Configuration
SUPABASE_URL="https://your-project.supabase.co/storage/v1"
SUPABASE_KEY="your-supabase-key"
SUPABASE_BUCKET="your-bucket-name"
```

### Running the Service

```bash
# Using default .env file
go run main.go start

# Using custom config file
go run main.go start --config=/path/to/.env

# Build binary
go build -o user-service main.go

# Run binary
./user-service start
```

Service akan berjalan di `http://localhost:8081` (atau port yang dikonfigurasi).

---

## 3. Project Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────────┐
│                        EXTERNAL LAYER                           │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │  Client  │  │RabbitMQ  │  │PostgreSQL│  │ Supabase │        │
│  └────┬─────┘  └────▲─────┘  └────▲─────┘  └────▲─────┘        │
└───────┼─────────────┼─────────────┼─────────────┼───────────────┘
        │             │             │             │
        ▼             │             │             │
┌─────────────────────────────────────────────────────────────────┐
│                      INTERFACE LAYER                            │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  Controller (HTTP Handler, Request/Response DTO)          │ │
│  └───────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────────────────────┐
│                      USECASE LAYER                              │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  Usecase (Business Logic, Orchestration)                  │ │
│  └───────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────────────────────┐
│                    REPOSITORY LAYER                             │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  Repository (Data Access, GORM Operations)                │ │
│  └───────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────────────────────┐
│                      INFRASTRUCTURE                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │   DB     │  │ RabbitMQ │  │  Config  │  │ Storage  │        │
│  │ Connection│  │ Service  │  │   Manager│  │ Service  │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└─────────────────────────────────────────────────────────────────┘
```

### Request Flow Diagram

```mermaid
flowchart TB
    subgraph "Client Request"
        C[HTTP Client]
    end

    subgraph "Interface Layer"
        R[Router<br/>app/routes.go]
        CT[Controller<br/>Handler & DTO]
        VAL[Validator<br/>pkg/validator]
    end

    subgraph "Business Layer"
        UC[Usecase<br/>Business Logic]
    end

    subgraph "Data Layer"
        RP[Repository<br/>Data Access]
        DB[(PostgreSQL)]
    end

    subgraph "External Services"
        RMQ[RabbitMQ]
        SB[Supabase Storage]
    end

    C -->|HTTP Request| R
    R -->|Route| CT
    CT -->|Validate| VAL
    VAL -->|Parsed Request| CT
    CT -->|Call| UC
    UC -->|Query| RP
    RP -->|SQL| DB
    DB -->|Data| RP
    RP -->|Entities| UC
    UC -->|Async Publish| RMQ
    UC -->|Return| CT
    CT -->|HTTP Response| C
    
    CT -.->|File Upload| SB
```

### Dependency Injection Container

```mermaid
flowchart TB
    subgraph "Infrastructure"
        CFG[Config<br/>configs.Config]
        DB[(PostgreSQL<br/>Connection)]
        RMQ[RabbitMQ<br/>Connection]
        SB[Supabase<br/>Client]
    end

    subgraph "Repository Layer"
        UR[UserRepository]
        RR[RoleRepository]
    end

    subgraph "Usecase Layer"
        UU[UserUsecase]
        RU[RoleUsecase]
    end

    subgraph "Service Layer"
        FUH[FileUploadHelper]
    end

    subgraph "Controller Layer"
        UC[UserController]
        RC[RoleController]
        AC[AuthController]
        UC4[UploadController]
    end

    subgraph "Container"
        CONT[App Container<br/>app.Container]
    end

    CFG -->|Inject| DB
    CFG -->|Inject| RMQ
    CFG -->|Inject| SB
    
    DB -->|Inject| UR
    DB -->|Inject| RR
    RMQ -->|Inject| UU
    
    UR -->|Inject| UU
    RR -->|Inject| RU
    
    UU -->|Inject| UC
    UU -->|Inject| AC
    RU -->|Inject| RC
    
    SB -->|Inject| FUH
    FUH -->|Inject| UC4
    
    UC -->|Register| CONT
    RC -->|Register| CONT
    AC -->|Register| CONT
    UC4 -->|Register| CONT
```

---

## 4. Application Setup

### Fiber Configuration

```go
app := fiber.New(fiber.Config{
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        zerolog.Printf("Error: %v", err)
        return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
    },
})
```

**Konfigurasi:**
- **Error Handler**: Global error handler dengan zerolog logging
- **Port**: Dari config atau environment variable `APP_PORT`
- **Timeout**: Context timeout 5 detik saat shutdown

### Middleware Stack

```mermaid
flowchart LR
    A[Request] --> B[Recover Middleware]
    B --> C[CORS Middleware]
    C --> D[Logger Middleware]
    D --> E[Route Handler]
    E --> F[Response]
```

| Middleware | Fungsi |
|------------|--------|
| `Recover` | Menangkap panic dan mencegah crash |
| `CORS` | Cross-Origin Resource Sharing |
| `Logger` | Logging setiap request dengan format: `[${time} ${ip} ${status} - ${latency} ${method} ${path}]` |

**Implementasi:**
```go
app.Use(recover.New())
app.Use(cors.New())
app.Use(logger.New(logger.Config{
    Format: "[${time} ${ip} ${status} - ${latency} ${method} ${path}\n",
}))
```

### Graceful Shutdown

```mermaid
sequenceDiagram
    participant OS as OS Signal
    participant S as Server
    participant C as Context
    
    S->>S: Start HTTP Server
    OS->>S: SIGTERM/SIGINT
    S->>S: Receive quit signal
    S->>C: Create timeout context (5s)
    S->>S: ShutdownWithContext(ctx)
    S->>S: Server shutdown complete
```

**Flow:**
1. Server berjalan di goroutine terpisah
2. Menunggu sinyal `SIGTERM` atau `SIGINT`
3. Membuat context dengan timeout 5 detik
4. Memanggil `ShutdownWithContext()` untuk menyelesaikan request yang sedang berjalan
5. Log "Server shutdown complete"

---

## 5. Configuration Management

### Environment Variables

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| `APP_ENV` | development | Environment (development/staging/production) |
| `APP_PORT` | 8081 | Port HTTP server |
| `DATABASE_HOST` | localhost | PostgreSQL host |
| `DATABASE_PORT` | 5432 | PostgreSQL port |
| `DATABASE_USER` | - | PostgreSQL username |
| `DATABASE_PASSWORD` | - | PostgreSQL password |
| `DATABASE_NAME` | - | PostgreSQL database name |
| `DATABASE_MAX_OPEN_CONNECTION` | 100 | Maximum open connections |
| `DATABASE_MAX_IDLE_CONNECTION` | 20 | Maximum idle connections |
| `RABBITMQ_HOST` | localhost | RabbitMQ host |
| `RABBITMQ_PORT` | 5672 | RabbitMQ port |
| `RABBITMQ_USERNAME` | guest | RabbitMQ username |
| `RABBITMQ_PASSWORD` | guest | RabbitMQ password |
| `REDIS_HOST` | - | Redis host (reserved) |
| `REDIS_PORT` | 6379 | Redis port (reserved) |
| `SUPABASE_URL` | - | Supabase Storage URL |
| `SUPABASE_KEY` | - | Supabase API Key |
| `SUPABASE_BUCKET` | - | Supabase bucket name |

### Config Structs

```go
type Config struct {
    App      App      `json:"app"`
    SqlDB    SqlDB    `json:"sql_db"`
    Redis    Redis    `json:"redis"`
    RabbitMQ RabbitMQ `json:"rabbitmq"`
    Supabase Supabase `json:"supabase"`
}

type App struct {
    AppPort string `json:"app_port"`
    AppEnv  string `json:"app_env"`
}

type SqlDB struct {
    Host           string `json:"host"`
    Port           string `json:"port"`
    User           string `json:"user"`
    Password       string `json:"password"`
    DBName         string `json:"db_name"`
    DBMaxOpenConns int    `json:"db_max_open_conns"`
    DBMaxIdleConns int    `json:"db_max_idle_conns"`
}

type RabbitMQ struct {
    Host     string `json:"host"`
    Port     string `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type Supabase struct {
    Url    string `json:"url"`
    Key    string `json:"key"`
    Bucket string `json:"bucket"`
}
```

---

## 6. Database Layer

### ER Diagram

```mermaid
erDiagram
    USER ||--o{ USER_ROLE : has
    ROLE ||--o{ USER_ROLE : assigned_to
    
    USER {
        uint id PK
        string name
        string email
        string password
        string photo
        string phone
        datetime created_at
        datetime updated_at
    }
    
    ROLE {
        uint id PK
        string name
        datetime created_at
        datetime updated_at
    }
    
    USER_ROLE {
        uint id PK
        uint user_id FK
        uint role_id FK
        datetime created_at
        datetime updated_at
    }
```

### Entity Definitions

#### User Model
```go
type User struct {
    ID       uint      `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Password string    `json:"password"`
    Photo    string    `json:"photo"`
    Phone    string    `json:"phone"`
    Roles    []Role    `gorm:"many2many:user_role;"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### Role Model
```go
type Role struct {
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Users     []User    `json:"users" gorm:"many2many:user_role;"`
}
```

#### UserRole Model (Join Table)
```go
type UserRole struct {
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    UserID    uint      `json:"user_id" gorm:"not null"`
    RoleID    uint      `json:"role_id" gorm:"not null"`
    User      User      `gorm:"foreignKey:UserID"`
    Role      Role      `gorm:"foreignKey:RoleID"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Custom table name
func (UserRole) TableName() string {
    return "user_role"
}
```

### Migrations & Seeders

#### Auto Migration
```go
db.AutoMigrate(&model.User{}, &model.Role{}, &model.UserRole{})
```

#### Seeders

**Role Seeder** (`database/role_seeder.go`):
```go
func SeedRole(db *gorm.DB) {
    roles := []model.Role{
        {Name: "Manager"},
        {Name: "Keeper"},
    }
    
    for _, role := range roles {
        db.FirstOrCreate(&role, model.Role{Name: role.Name})
    }
}
```

**Manager Seeder** (`database/manager_seeder.go`):
```go
func SeedManager(db *gorm.DB) {
    // Buat user manager default
    admin := model.User{
        Name:     "manager",
        Email:    "manager@mail.com",
        Password: hashedPassword,
        Roles:    []model.Role{modelRole},
    }
    
    db.FirstOrCreate(&admin, model.User{Email: "manager@mail.com"})
}
```

### Connection Pooling

```go
sqlDB.SetMaxIdleConns(cfg.SqlDB.DBMaxIdleConns)      // 20
sqlDB.SetMaxOpenConns(cfg.SqlDB.DBMaxOpenConns)      // 100
```

---

## 7. API Documentation

### Authentication

#### POST `/api/v1/auth/login`
Login user dengan email dan password.

**Request:**
```json
{
  "email": "user@mail.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "message": "Login successful",
  "user": {
    "user_id": 1,
    "email": "user@mail.com",
    "role_name": ["Manager"]
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body atau validation error
- `404 Not Found`: User tidak ditemukan
- `401 Unauthorized`: Password salah

### User Management

#### POST `/api/v1/users`
Create new user.

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@mail.com",
  "password": "password123",
  "phone": "08123456789",
  "photo": "https://example.com/photo.jpg"
}
```

**Response (201 Created):**
```json
{
  "message": "User created successfully"
}
```

#### GET `/api/v1/users`
Get all users dengan pagination.

**Query Parameters:**
| Parameter | Type | Default | Deskripsi |
|-----------|------|---------|-----------|
| `page` | int | 1 | Nomor halaman |
| `limit` | int | 10 | Jumlah item per halaman (max 100) |
| `search` | string | - | Search by name atau email |
| `sort_by` | string | created_at | Sort field (id, name, email, created_at) |
| `sort_order` | string | desc | Sort direction (asc, desc) |

**Response (200 OK):**
```json
{
  "data": {
    "users": [
      {
        "id": 1,
        "name": "John Doe",
        "email": "john@mail.com",
        "phone": "08123456789",
        "photo": "https://example.com/photo.jpg",
        "role_name": "Manager"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_records": 50,
      "limit": 10,
      "has_next": true,
      "has_prev": false
    }
  },
  "message": "Users fetched successfully"
}
```

#### GET `/api/v1/users/:id`
Get user by ID.

**Response (200 OK):**
```json
{
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@mail.com",
    "phone": "08123456789",
    "photo": "https://example.com/photo.jpg",
    "role_name": "Manager"
  }
}
```

#### PUT `/api/v1/users/:id`
Update user.

**Request:**
```json
{
  "name": "John Updated",
  "email": "john.updated@mail.com",
  "password": "newpassword123",
  "phone": "08987654321",
  "photo": "https://example.com/new-photo.jpg"
}
```

**Note:** Password bersifat opsional (`omitempty`). Jika tidak diisi, password tidak berubah.

**Response (200 OK):**
```json
{
  "message": "User updated successfully"
}
```

#### DELETE `/api/v1/users/:id`
Delete user.

**Response (200 OK):**
```json
{
  "message": "User deleted successfully"
}
```

#### GET `/api/v1/users/role/:roleName`
Get users by role name.

**Example:** `GET /api/v1/users/role/Manager`

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@mail.com",
      "phone": "08123456789",
      "photo": "https://example.com/photo.jpg",
      "role_name": "Manager"
    }
  ],
  "message": "Users fetched successfully"
}
```

### Role Management

#### POST `/api/v1/roles`
Create new role.

**Request:**
```json
{
  "name": "Admin"
}
```

**Response (201 Created):**
```json
{
  "message": "Role created successfully"
}
```

#### GET `/api/v1/roles`
Get all roles.

**Response (200 OK):**
```json
{
  "message": "Roles fetched successfully",
  "data": [
    {
      "id": 1,
      "name": "Manager",
      "count_users": 5
    },
    {
      "id": 2,
      "name": "Keeper",
      "count_users": 10
    }
  ]
}
```

#### GET `/api/v1/roles/:id`
Get role by ID.

**Response (200 OK):**
```json
{
  "message": "Role fetched successfully",
  "data": {
    "id": 1,
    "name": "Manager",
    "users": [...],
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```

#### PUT `/api/v1/roles/:id`
Update role.

**Request:**
```json
{
  "name": "Administrator"
}
```

**Response (200 OK):**
```json
{
  "message": "Role updated successfully"
}
```

#### DELETE `/api/v1/roles/:id`
Delete role.

**Constraint:** Role tidak bisa dihapus jika masih memiliki users.

**Error Response (409 Conflict):**
```json
{
  "message": "role has users"
}
```

### User-Role Assignment

#### POST `/api/v1/assign-role`
Assign user ke role.

**Request:**
```json
{
  "user_id": 1,
  "role_id": 2
}
```

**Response (200 OK):**
```json
{
  "message": "User assigned to role successfully"
}
```

#### GET `/api/v1/assign-role`
Get all user-role assignments dengan pagination.

**Response (200 OK):**
```json
{
  "data": {
    "user_roles": [
      {
        "id": 1,
        "user_id": 1,
        "role_id": 2,
        "user": {
          "id": 1,
          "name": "John Doe",
          "email": "john@mail.com",
          "phone": "08123456789",
          "photo": "https://example.com/photo.jpg"
        },
        "role": {
          "id": 2,
          "name": "Keeper"
        }
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_records": 50,
      "limit": 10,
      "has_next": true,
      "has_prev": false
    }
  },
  "message": "User roles fetched successfully"
}
```

#### GET `/api/v1/assign-role/:id`
Get user-role assignment by ID.

#### PUT `/api/v1/assign-role/:id`
Edit user-role assignment.

**Request:**
```json
{
  "user_id": 2,
  "role_id": 3
}
```

**Response (200 OK):**
```json
{
  "message": "User role updated successfully"
}
```

### File Upload

#### POST `/api/v1/upload/photo`
Upload user photo.

**Content-Type:** `multipart/form-data`

**Form Field:**
- `image`: File image (jpg, jpeg, png, webp, svg)

**Constraints:**
- Max file size: 2MB
- Allowed extensions: .jpg, .jpeg, .png, .webp, .svg

**Response (200 OK):**
```json
{
  "message": "File uploaded successfully",
  "data": {
    "url": "https://ikowborhuhawxgpoxcxs.supabase.co/storage/v1/object/public/bwa-warehouse/users/image_1704067200.jpg",
    "path": "users/image_1704067200.jpg",
    "filename": "image_1704067200.jpg"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "message": "file size exceeds the maximum allowed size"
}
```

**Error Response (400 Bad Request):**
```json
{
  "message": "invalid file extension"
}
```

### Request DTOs

#### Auth Requests
```go
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}
```

#### User Requests
```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
    Phone    string `json:"phone" validate:"required"`
    Photo    string `json:"photo" validate:"required"`
}

type UpdateUserRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"omitempty,min=8"`
    Phone    string `json:"phone"`
    Photo    string `json:"photo"`
}

type GetAllUsersRequest struct {
    Page      int    `query:"page" validate:"omitempty,min=1"`
    Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
    Search    string `query:"search" validate:"omitempty"`
    SortBy    string `query:"sort_by" validate:"omitempty,oneof=id name email created_at"`
    SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}
```

#### Role Requests
```go
type CreateRoleRequest struct {
    Name string `json:"name" validate:"required"`
}
```

#### Assign Role Requests
```go
type AssignUserToRoleRequest struct {
    UserID uint `json:"user_id" validate:"required"`
    RoleID uint `json:"role_id" validate:"required"`
}
```

### Response DTOs

#### Auth Responses
```go
type LoginResponse struct {
    UserID uint     `json:"user_id"`
    Email  string   `json:"email"`
    Role   []string `json:"role_name"`
}
```

#### User Responses
```go
type UserResponse struct {
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Phone    string `json:"phone"`
    Photo    string `json:"photo"`
    RoleName string `json:"role_name"`
}

type GetAllUsersResponse struct {
    Users      []UserResponse                `json:"users"`
    Pagination pagination.PaginationResponse `json:"pagination"`
}
```

#### UserRole Responses
```go
type UserRoleResponse struct {
    ID     uint         `json:"id"`
    UserID uint         `json:"user_id"`
    RoleID uint         `json:"role_id"`
    User   UserResponse `json:"user"`
    Role   RoleResponse `json:"role"`
}

type GetAllUserRolesResponse struct {
    UserRoles  []UserRoleResponse            `json:"user_roles"`
    Pagination pagination.PaginationResponse `json:"pagination"`
}
```

#### Role Responses
```go
type RoleResponse struct {
    ID         uint   `json:"id"`
    Name       string `json:"name"`
    CountUsers int64  `json:"count_users"`
}
```

#### Upload Responses
```go
type UploadPhotoResponse struct {
    URL      string `json:"url"`
    Path     string `json:"path"`
    Filename string `json:"filename"`
}
```

---

## 8. Business Logic Deep Dive

### User Creation Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant CT as UserController
    participant VAL as Validator
    participant UC as UserUsecase
    ├── Conv as Password Hash
    participant RP as UserRepository
    participant DB as Database
    participant RMQ as RabbitMQ

    C->>CT: POST /api/v1/users
    CT->>CT: BodyParser(&req)
    CT->>VAL: Validate(req)
    VAL-->>CT: OK
    CT->>CT: Convert to model.User
    CT->>UC: CreateUser(ctx, user)
    UC->>UC: HashPassword(password)
    UC->>RP: CreateUser(ctx, user)
    RP->>DB: INSERT INTO users
    DB-->>RP: user dengan ID
    RP-->>UC: user
    
    par Async Email Event
        UC->>RMQ: PublishEmail(ctx, payload)
    end
    
    UC-->>CT: nil
    CT-->>C: 201 Created
```

**Key Points:**
1. Password di-hash menggunakan bcrypt (cost: 14) di usecase layer
2. Password asli disimpan sementara untuk dikirim via email
3. Event publishing ke RabbitMQ berjalan di goroutine (async)
4. Response 201 dikirim tanpa menunggu RabbitMQ

### Login Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant CT as AuthController
    participant VAL as Validator
    participant UC as UserUsecase
    participant RP as UserRepository
    participant DB as Database
    participant Conv as Password Check

    C->>CT: POST /api/v1/auth/login
    CT->>CT: BodyParser(&req)
    CT->>VAL: Validate(loginRequest)
    VAL-->>CT: OK
    CT->>UC: GetUserByEmail(ctx, email)
    UC->>RP: GetUserByEmail(ctx, email)
    RP->>DB: SELECT ... WHERE email = ?
    DB-->>RP: user dengan Roles
    RP-->>UC: user
    UC-->>CT: user
    CT->>Conv: CheckPasswordHash(password, user.Password)
    Conv-->>CT: true/false
    
    alt Password Valid
        CT->>CT: Extract role names
        CT-->>C: 200 OK + LoginResponse
    else Password Invalid
        CT-->>C: 401 Unauthorized
    end
```

**Key Points:**
1. User di-fetch dengan preload Roles (many-to-many)
2. Password diverifikasi menggunakan bcrypt
3. Response mengembalikan array role names
4. Tidak ada JWT generation (dilakukan di API Gateway)

### Role Assignment Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant CT as UserController
    participant VAL as Validator
    participant UC as UserUsecase
    participant RP as UserRepository
    participant DB as Database

    C->>CT: POST /api/v1/assign-role
    CT->>CT: BodyParser(&req)
    CT->>VAL: Validate(req)
    VAL-->>CT: OK
    CT->>UC: AssignUserToRole(ctx, userID, roleID)
    UC->>RP: AssignUserToRole(ctx, userID, roleID)
    RP->>RP: Check context cancellation
    RP->>DB: INSERT INTO user_role
    DB-->>RP: success
    RP-->>UC: nil
    UC-->>CT: nil
    CT-->>C: 200 OK
```

**Key Points:**
1. Check context cancellation sebelum query
2. Insert ke join table `user_role`
3. Foreign key constraints memastikan user dan role exist

### File Upload Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant CT as UploadController
    participant FUH as FileUploadHelper
    participant VAL as Validation
    participant SB as SupabaseStorage
    participant S3 as Supabase S3

    C->>CT: POST /api/v1/upload/photo (multipart/form-data)
    CT->>CT: FormFile("image")
    CT->>FUH: UploadPhoto(ctx, file)
    FUH->>VAL: validateImageFile(file, 2MB)
    
    alt Validation Failed
        VAL-->>FUH: Error
        FUH-->>CT: Error
        CT-->>C: 400 Bad Request
    else Validation Success
        VAL-->>FUH: OK
        FUH->>SB: UploadFile(ctx, file, "users")
        SB->>SB: Generate unique filename (timestamp)
        SB->>SB: Detect content type
        SB->>S3: Upload to bucket
        S3-->>SB: Success
        SB->>SB: GetPublicUrl
        SB-->>FUH: UploadResult
        FUH-->>CT: UploadResult
        CT-->>C: 200 OK + URL
    end
```

**Key Points:**
1. Validasi ukuran file (max 2MB)
2. Validasi ekstensi (.jpg, .jpeg, .png, .webp, .svg)
3. Filename digenerate dengan timestamp untuk uniqueness
4. Content type otomatis terdeteksi
5. File disimpan di folder `users/` dalam bucket

---

## 9. External Integrations

### RabbitMQ Integration

**Service File:** `service/rabbitmq_service.go`

```mermaid
flowchart TB
    subgraph "User Service"
        UC[UserUsecase<br/>CreateUser]
        RS[RabbitMQService]
    end
    
    subgraph "RabbitMQ Server"
        Q[Queue: email_queue]
    end
    
    subgraph "Consumer Services"
        NS[Notification Service]
        ES[Email Service]
    end
    
    UC -->|PublishEmail| RS
    RS -->|Declare Queue| Q
    RS -->|Publish| Q
    Q -->|Consume| NS
    Q -->|Consume| ES
```

**Email Payload Structure:**
```go
type EmailPayload struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Type     string `json:"type"`
    UserID   uint   `json:"user_id"`
    Name     string `json:"name"`
}
```

**Example Payload:**
```json
{
  "email": "john@mail.com",
  "password": "plaintextpassword",
  "type": "welcome_email",
  "user_id": 1,
  "name": "John Doe"
}
```

**Connection & Publishing:**
```go
// Connection
conn, err := amqp.Dial("amqp://username:password@host:port/")
ch, err := conn.Channel()

// Declare Queue (durable)
queue, err := ch.QueueDeclare(
    "email_queue", // name
    true,          // durable
    false,         // delete when unused
    false,         // exclusive
    false,         // no-wait
    nil,           // arguments
)

// Publish
err = ch.Publish(
    "",           // exchange (default)
    queue.Name,   // routing key
    false,        // mandatory
    false,        // immediate
    amqp.Publishing{
        ContentType: "application/json",
        Body:        body,
    },
)
```

**Async Pattern:**
```go
// Di UserUsecase.CreateUser
go func() {
    if err := u.rabbitMQService.PublishEmail(ctx, emailPayload); err != nil {
        log.Errorf("[UserUsecase] CreateUser - 3: %v", err)
    }
}()
```

**Catatan:** Publishing berjalan di goroutine terpisah agar tidak blocking response HTTP.

### Supabase Storage

**Service File:** `pkg/storage/supabase_storage.go`

```mermaid
flowchart LR
    subgraph "User Service"
        UC[UploadController]
        FUH[FileUploadHelper]
        SS[SupabaseStorage]
    end
    
    subgraph "Supabase"
        API[Storage API]
        Bucket[(bwa-warehouse<br/>Bucket)]
    end
    
    UC -->|UploadPhoto| FUH
    FUH -->|Validate| FUH
    FUH -->|UploadFile| SS
    SS -->|HTTP API| API
    API -->|Store| Bucket
```

**Upload Configuration:**
```go
const (
    MaxImageSize           = 2 * 1024 * 1024 // 2 MB
    AllowedImageExtensions = ".jpg,.jpeg,.png,.webp,.svg"
)
```

**Upload Process:**
1. Generate unique filename dengan timestamp
2. Detect content type berdasarkan extension
3. Create Supabase Storage client
4. Upload file ke bucket
5. Return public URL

**File Naming:**
```go
// Format: originalname_timestamp.ext
ext := filepath.Ext(file.Filename)                              // .jpg
timestamp := time.Now().Unix()                                  // 1704067200
filename := fmt.Sprintf("%s_%d%s", 
    strings.TrimSuffix(file.Filename, ext),  // image
    timestamp,                               // 1704067200
    ext)                                     // .jpg
// Result: image_1704067200.jpg
```

---

## 10. Utility Packages

### pkg/conv

**File:** `pkg/conv/conv.go`

Utility untuk password hashing dan type conversion.

**Password Hashing:**
```go
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

**Type Conversion:**
```go
func StringToUint(s string) uint {
    id, err := strconv.ParseUint(s, 10, 64)
    if err != nil {
        return 0
    }
    return uint(id)
}
```

### pkg/pagination

**File:** `pkg/pagination/pagination.go`

**PaginationResponse Struct:**
```go
type PaginationResponse struct {
    CurrentPage  int   `json:"current_page"`
    TotalPages   int   `json:"total_pages"`
    TotalRecords int64 `json:"total_records"`
    Limit        int   `json:"limit"`
    HasNext      bool  `json:"has_next"`
    HasPrev      bool  `json:"has_prev"`
}
```

**Calculation:**
```go
func CalculatePagination(page, limit, totalRecords int) PaginationResponse {
    totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))
    if totalPages == 0 {
        totalPages = 1
    }

    return PaginationResponse{
        CurrentPage:  page,
        TotalPages:   totalPages,
        TotalRecords: int64(totalRecords),
        Limit:        limit,
        HasNext:      page < totalPages,
        HasPrev:      page > 1,
    }
}
```

### pkg/validator

**File:** `pkg/validator/request_validator.go`

Menggunakan `go-playground/validator/v10` dengan custom error messages.

**Validation Rules:**
| Tag | Error Message |
|-----|---------------|
| `required` | `{Field} is required` |
| `email` | `{Field} is not a valid email` |
| `min` | `{Field} must be at least {Param} characters long` |
| `omitempty` | Skip validation if empty |
| `oneof` | Value must be one of allowed options |

**Usage:**
```go
if err := validator.Validate(req); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "message": err.Error(),
    })
}
```

**Example Error:**
```
Validation failed: Name is required, Email is not a valid email, Password must be at least 8 characters long
```

### pkg/storage

**File:** `pkg/storage/file_upload_helper.go`

**FileUploadHelper:**
```go
type FileUploadHelper struct {
    storage SupabaseInterface
    cfg     configs.Config
}

func (h *FileUploadHelper) UploadPhoto(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error)
```

**Validation Functions:**
```go
func validateFileSize(size int64, maxSize int64) bool
func validateFileExtension(extension string, allowedExtensions string) bool
func getFileExtension(filename string) string
```

---

## 11. Testing Guide

### Unit Testing

**Struktur Test:**
```
repository/
  user_repository.go
  user_repository_test.go
usecase/
  user_usecase.go
  user_usecase_test.go
controller/
  user_controller.go
  user_controller_test.go
```

**Mocking Strategy:**

1. **Repository Layer:** Mock GORM DB menggunakan `sqlmock`
2. **Usecase Layer:** Mock Repository menggunakan testify mock
3. **Controller Layer:** Mock Usecase dan menggunakan `httptest`

**Example Repository Test:**
```go
func TestUserRepository_CreateUser(t *testing.T) {
    // Setup sqlmock
    db, mock, err := sqlmock.New()
    gormDB, _ := gorm.Open(postgres.New(postgres.Config{
        Conn: db,
    }), &gorm.Config{})
    
    repo := NewUserRepository(gormDB)
    
    // Setup expectations
    mock.ExpectBegin()
    mock.ExpectQuery("INSERT INTO \"users\"").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
    mock.ExpectCommit()
    
    // Execute
    user := model.User{Name: "Test", Email: "test@mail.com"}
    result, err := repo.CreateUser(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, uint(1), result.ID)
}
```

**Example Usecase Test:**
```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
    args := m.Called(ctx, user)
    return args.Get(0).(*model.User), args.Error(1)
}

func TestUserUsecase_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockRabbit := new(MockRabbitMQService)
    
    usecase := NewUserUsecase(mockRepo, mockRabbit)
    
    user := model.User{Name: "Test", Password: "password123"}
    mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("model.User")).
        Return(&model.User{ID: 1}, nil)
    
    err := usecase.CreateUser(context.Background(), user)
    
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### Integration Testing

**Setup:**
1. Gunakan test database (PostgreSQL container)
2. Run migrations
3. Seed test data
4. Execute tests
5. Cleanup database

**Example:**
```go
func TestIntegration_CreateUser(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    
    // Build container
    container := BuildTestContainer(db)
    
    // Create test request
    app := fiber.New()
    SetupRoutes(app, container)
    
    req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(`{
        "name": "Test User",
        "email": "test@mail.com",
        "password": "password123",
        "phone": "08123456789",
        "photo": "http://example.com/photo.jpg"
    }`))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := app.Test(req)
    
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
}
```

---

## 12. Deployment

### Build Command

```bash
# Development build
go build -o user-service main.go

# Production build (optimized)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service main.go

# Build with version info
go build -ldflags "-X main.version=1.0.0 -X main.buildTime=$(date -u +%Y%m%d%H%M%S)" -o user-service main.go
```

### Environment-specific Considerations

**Development:**
```env
APP_ENV="development"
APP_PORT=8081
DATABASE_HOST=localhost
RABBITMQ_HOST=localhost
```

**Production:**
```env
APP_ENV="production"
APP_PORT=8080
DATABASE_HOST=prod-db-host
DATABASE_MAX_OPEN_CONNECTION=200
DATABASE_MAX_IDLE_CONNECTION=50
RABBITMQ_HOST=prod-rabbitmq-host
```

**Health Check:**
```bash
# Health check endpoint (tambahkan jika diperlukan)
curl http://localhost:8081/health
```

**Docker Example:**
```dockerfile
FROM golang:1.24.7-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/user-service .
CMD ["./user-service", "start"]
```

---

## 13. Troubleshooting

### Common Issues

#### 1. Database Connection Failed
**Error:**
```
[Postgres] ConnectionPostgres - 1: dial tcp: connect: connection refused
```

**Solution:**
- Cek apakah PostgreSQL running: `sudo systemctl status postgresql`
- Cek environment variables di `.env`
- Pastikan port tidak ter-block: `telnet localhost 5432`

#### 2. RabbitMQ Connection Failed
**Error:**
```
[RabbitMQService] NewRabbitMQService - 1: dial tcp: connect: connection refused
```

**Solution:**
- Cek RabbitMQ status: `sudo systemctl status rabbitmq-server`
- Verifikasi credentials di `.env`
- Cek RabbitMQ Management UI: `http://localhost:15672`

#### 3. Validation Error
**Error:**
```
Validation failed: Email is not a valid email
```

**Solution:**
- Cek format request body
- Pastikan semua field `required` terisi
- Cek validation tags di struct

#### 4. File Upload Error
**Error:**
```
file size exceeds the maximum allowed size
```

**Solution:**
- Pastikan file < 2MB
- Cek ekstensi file (hanya .jpg, .jpeg, .png, .webp, .svg)
- Verifikasi Supabase credentials

#### 5. Role Delete Failed
**Error:**
```
role has users
```

**Solution:**
- Hapus atau reassign users yang terkait dengan role terlebih dahulu
- Gunakan endpoint `GET /api/v1/roles/:id` untuk cek users

### Debug Tips

#### 1. Enable Debug Logging
```go
// di app/app.go
cfg.App.AppEnv = "development" // enable verbose logging
```

#### 2. Check Context Cancellation
```go
// Tambahkan log di repository
select {
case <-ctx.Done():
    log.Errorf("[Repository] Context cancelled: %v", ctx.Err())
    return ctx.Err()
default:
}
```

#### 3. Database Query Log
```go
// Enable GORM debug mode
db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

#### 4. RabbitMQ Queue Check
```bash
# List queues
rabbitmqctl list_queues

# Check queue messages
rabbitmqctl list_queues name messages_ready messages_unacknowledged
```

#### 5. Environment Validation
```bash
# Print all environment variables
go run main.go start 2>&1 | grep -i "config\|env"
```

### Log Pattern Reference

| Component | Format | Example |
|-----------|--------|---------|
| Controller | `[ControllerName] MethodName - step: error` | `[UserController] CreateUser - 1: invalid request` |
| Usecase | `[UsecaseName] MethodName - step: error` | `[UserUsecase] CreateUser - 2: db error` |
| Repository | `[RepositoryName] MethodName - step: error` | `[UserRepository] CreateUser - 1: context cancelled` |
| Service | `[ServiceName] MethodName - step: error` | `[RabbitMQService] PublishEmail - 1: connection failed` |

---

## 📎 Appendix

### API Endpoint Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Login user |
| POST | `/api/v1/users` | Create user |
| GET | `/api/v1/users` | Get all users (paginated) |
| GET | `/api/v1/users/:id` | Get user by ID |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Delete user |
| GET | `/api/v1/users/role/:roleName` | Get users by role |
| POST | `/api/v1/roles` | Create role |
| GET | `/api/v1/roles` | Get all roles |
| GET | `/api/v1/roles/:id` | Get role by ID |
| PUT | `/api/v1/roles/:id` | Update role |
| DELETE | `/api/v1/roles/:id` | Delete role |
| POST | `/api/v1/assign-role` | Assign user to role |
| GET | `/api/v1/assign-role` | Get all assignments |
| GET | `/api/v1/assign-role/:id` | Get assignment by ID |
| PUT | `/api/v1/assign-role/:id` | Update assignment |
| POST | `/api/v1/upload/photo` | Upload photo |

---

**End of Documentation**

*Dokumentasi ini dibuat untuk membantu tim baru memahami 100% cara kerja User Service. Jika ada pertanyaan atau perlu klarifikasi, silakan hubungi tim engineering.*
