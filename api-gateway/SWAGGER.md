# Dokumentasi Implementasi Swagger

Dokumen ini menjelaskan cara Swagger diimplementasikan pada sistem Warehouse Microservice, mulai dari anotasi per-service hingga agregasi spec di API Gateway.

---

## Daftar Isi

- [Arsitektur](#arsitektur)
- [Setup Per-Service](#setup-per-service)
- [Anotasi Handler](#anotasi-handler)
- [Setup API Gateway](#setup-api-gateway)
- [Cara Akses](#cara-akses)
- [Troubleshooting](#troubleshooting)

---

## Arsitektur

Setiap service mengekspos endpoint `/swagger/doc.json` yang berisi spesifikasi Swagger 2.0 hasil generate `swag init`. API Gateway mengambil semua spec tersebut secara concurrent, menyatukannya menjadi satu spec gabungan, lalu menyajikannya melalui Swagger UI.

```
Browser
  └── GET /swagger
        └── Swagger UI membaca /swagger/aggregated.json
              └── API Gateway: AggregateSpecs()
                    ├── GET http://user-service:8081/swagger/doc.json
                    ├── GET http://product-service:8082/swagger/doc.json
                    ├── GET http://warehouse-service:8083/swagger/doc.json
                    ├── GET http://merchant-service:8084/swagger/doc.json
                    ├── GET http://transaction-service:8085/swagger/doc.json
                    └── GET http://notification-service:8086/swagger/doc.json
```

---

## Setup Per-Service

### 1. Install swaggo

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. Tambahkan dependency di go.mod

```bash
go get github.com/swaggo/swag
go get github.com/swaggo/fiber-swagger
```

### 3. Anotasi `main.go`

Tambahkan anotasi berikut di atas fungsi `main()`. Anotasi ini menjadi metadata utama spec yang di-generate.

```go
// @title User Service API
// @version 1.0
// @description This is the API documentation for the User Service.
// @host user-service:8081
// @BasePath /api/v1
// @SecurityDefinitions.apiKey Bearer
// @in header
// @name Authorization
func main() {
    cmd.Execute()
}
```

| Anotasi | Keterangan |
|---------|------------|
| `@title` | Judul API yang tampil di Swagger UI |
| `@version` | Versi API |
| `@description` | Deskripsi singkat service |
| `@host` | Host service (Docker hostname:port) |
| `@BasePath` | Base path semua endpoint |
| `@SecurityDefinitions.apiKey` | Nama skema keamanan |
| `@in header` | Token dikirim via HTTP header |
| `@name Authorization` | Nama header yang digunakan |

### 4. Generate docs

Jalankan perintah berikut di direktori root service setiap kali anotasi berubah:

```bash
swag init
```

Perintah ini akan membuat atau memperbarui direktori `docs/` yang berisi `docs.go`, `swagger.json`, dan `swagger.yaml`.

### 5. Daftarkan route Swagger di service

```go
import (
    fiberSwagger "github.com/swaggo/fiber-swagger"
    _ "micro-warehouse/user-service/docs"
)

// di dalam setup routes:
app.Get("/swagger/*", fiberSwagger.WrapHandler)
```

---

## Anotasi Handler

### POST / PUT / DELETE — Body Parameter

Untuk endpoint yang menerima request body, gunakan `in: body` dengan referensi ke struct request.

```go
// @Summary Create User
// @Description Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body request.CreateUserRequest true "Create User Request Body"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users [post]
// @Security Bearer
func (u *userController) CreateUser(c *fiber.Ctx) error {
```

### GET — Query Parameter

> **Penting:** Endpoint GET **tidak boleh** menggunakan `in: body`. Browser akan menolak request GET yang memiliki body dengan error `TypeError: Failed to execute 'fetch' on 'Window': Request with GET/HEAD method cannot have body`.

Gunakan `in: query` untuk setiap field parameter pada endpoint GET.

```go
// @Summary Get All Users
// @Description Get all users
// @Tags Users
// @Accept json
// @Produce json
// @Param page      query int    false "Page number"
// @Param limit     query int    false "Items per page"
// @Param search    query string false "Search keyword"
// @Param sortBy    query string false "Sort by field"
// @Param sortOrder query string false "Sort order (asc/desc)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users [get]
// @Security Bearer
func (u *userController) GetAllUsers(c *fiber.Ctx) error {
```

### Path Parameter

```go
// @Param id path string true "User ID"
// @Router /api/v1/users/{id} [get]
// @Security Bearer
func (u *userController) GetUserByID(c *fiber.Ctx) error {
```

### Referensi Anotasi

| Anotasi | Format | Keterangan |
|---------|--------|------------|
| `@Summary` | teks | Judul singkat endpoint |
| `@Description` | teks | Deskripsi panjang |
| `@Tags` | nama | Nama grup dropdown di Swagger UI |
| `@Accept` | mime | Format input (`json`, `multipart/form-data`) |
| `@Produce` | mime | Format output (`json`) |
| `@Param` | `nama lokasi tipe wajib "deskripsi"` | Definisi parameter |
| `@Success` | `kode {tipe} model` | Definisi response sukses |
| `@Failure` | `kode {tipe} model` | Definisi response error |
| `@Router` | `/path [method]` | Path dan HTTP method endpoint |
| `@Security` | nama skema | Menandai endpoint butuh autentikasi |

---

## Setup API Gateway

### Struktur File

```
api-gateway/
└── internal/
    └── swagger/
        └── aggregator.go   # Logic fetch & merge spec
└── main.go                 # Route /swagger dan /swagger/aggregated.json
```

### Route Swagger di `main.go`

API Gateway menyajikan dua endpoint swagger:

```go
func setupSwaggerRoutes(app *fiber.App) {
    // Swagger UI
    app.Get("/swagger", func(c *fiber.Ctx) error {
        html := `...
            SwaggerUIBundle({
                url: "/swagger/aggregated.json",
                dom_id: '#swagger-ui',
                tagsSorter: 'alpha'   // mengurutkan tag A-Z
            })
        ...`
        c.Set("Content-Type", "text/html")
        return c.SendString(html)
    })

    // Aggregated spec JSON
    app.Get("/swagger/aggregated.json", func(c *fiber.Ctx) error {
        spec, err := swaggerAgg.AggregateSpecs(c.BaseURL())
        if err != nil {
            log.Printf("[SetupSwaggerRoutes] AggregateSpecs: %v", err)
            return c.Status(fiber.StatusInternalServerError).JSON(...)
        }
        c.Set("Content-Type", "application/json")
        data, _ := json.Marshal(spec)
        return c.Send(data)
    })
}
```

### Cara Kerja `aggregator.go`

#### Daftar Service

```go
var Services = []ServiceSpec{
    {Name: "user-service",         URL: "http://user-service:8081/swagger/doc.json"},
    {Name: "product-service",      URL: "http://product-service:8082/swagger/doc.json"},
    {Name: "warehouse-service",    URL: "http://warehouse-service:8083/swagger/doc.json"},
    {Name: "merchant-service",     URL: "http://merchant-service:8084/swagger/doc.json"},
    {Name: "transaction-service",  URL: "http://transaction-service:8085/swagger/doc.json"},
    {Name: "notification-service", URL: "http://notification-service:8086/swagger/doc.json"},
}
```

URL menggunakan Docker hostname karena API Gateway dan semua service berjalan dalam satu Docker network.

#### `FetchSpec`

Mengambil `swagger/doc.json` dari satu service via HTTP GET dengan header internal:

```go
req.Header.Set("X-Gateway", "warehouse-api-gateway")
req.Header.Set("X-Internal-Request", "true")
```

#### `AggregateSpecs`

Struktur spec gabungan yang dihasilkan:

```go
merged := map[string]any{
    "swagger": "2.0",
    "info": map[string]any{
        "title":   "Warehouse Microservices API",
        "version": "1.0.0",
    },
    "host":                "localhost:8080",
    "basePath":            "/",
    "paths":               map[string]any{},   // diisi dari semua service
    "definitions":         map[string]any{},   // diisi dari semua service
    "securityDefinitions": map[string]any{},   // diisi dari semua service
}
```

Proses fetch dilakukan secara **concurrent** dengan goroutine. Jika ada service yang tidak bisa dijangkau, proses merge tetap berlanjut untuk service lainnya — hanya dicatat sebagai warning di log.

#### `rewriteRefs` — Mencegah Konflik Nama Schema

Karena semua service mungkin memiliki nama schema yang sama (misalnya `request.CreateUserRequest` dan `request.CreateProductRequest` sama-sama bernama `CreateRequest`), setiap schema diberi prefix nama service sebelum di-merge:

```
request.GetAllUsersRequest  →  user-service_request.GetAllUsersRequest
request.CreateProductRequest  →  product-service_request.CreateProductRequest
```

Semua `$ref` di dalam `paths` juga diperbarui agar menunjuk ke nama baru:

```
// Sebelum
"$ref": "#/definitions/request.GetAllUsersRequest"

// Setelah
"$ref": "#/definitions/user-service_request.GetAllUsersRequest"
```

---

## Cara Akses

### Membuka Swagger UI

```
http://localhost:8080/swagger
```

### Menggunakan Authorization

1. Lakukan login melalui endpoint `POST /api/v1/auth/login`
2. Salin nilai `token` dari response
3. Klik tombol **Authorize** di pojok kanan atas Swagger UI
4. Masukkan token dengan format:
   ```
   Bearer <token>
   ```
5. Klik **Authorize**, lalu tutup dialog
6. Semua endpoint yang ditandai `// @Security Bearer` sekarang akan menyertakan header `Authorization` secara otomatis

### Raw Spec JSON

Untuk mengakses spec gabungan secara langsung:

```
http://localhost:8080/swagger/aggregated.json
```

---

## Troubleshooting

### 1. Internal Server Error saat membuka `/swagger`

**Penyebab:** Panic karena `merged["paths"]` adalah `nil` ketika diakses dengan type assertion langsung.

**Solusi:** Pastikan field `"paths"` dan `"definitions"` diinisialisasi sebagai `map[string]any{}` di struct `merged` di dalam `AggregateSpecs`, bukan diletakkan di dalam field `"info"`.

---

### 2. Swagger UI tidak bisa diakses saat development lokal

**Penyebab:** URL di `Services` menggunakan Docker hostname (`http://user-service:8081`) yang hanya bisa di-resolve di dalam Docker network.

**Solusi:** Saat menjalankan API Gateway di luar Docker, ganti URL ke `localhost`:

```go
// Development lokal
{Name: "user-service", URL: "http://localhost:8081/swagger/doc.json"},

// Di dalam Docker (production/staging)
{Name: "user-service", URL: "http://user-service:8081/swagger/doc.json"},
```

Atau gunakan environment variable agar bisa dikonfigurasi tanpa mengubah kode.

---

### 3. `TypeError: Failed to execute 'fetch' ... GET/HEAD method cannot have body`

**Penyebab:** Anotasi `@Param` pada endpoint GET menggunakan `in: body`, padahal browser melarang GET request memiliki body.

**Solusi:** Ganti semua `body` menjadi `query` pada endpoint GET, lalu jalankan ulang `swag init`:

```go
// Salah
// @Param request body request.GetAllUsersRequest true "..."

// Benar
// @Param page      query int    false "Page number"
// @Param limit     query int    false "Items per page"
// @Param search    query string false "Search keyword"
```

---

### 4. Tombol Authorize tidak muncul di Swagger UI

**Penyebab:** Field `securityDefinitions` tidak di-merge ke spec gabungan, atau spec menggunakan format hybrid Swagger 2.0 + OpenAPI 3.0 (`components/schemas`) secara bersamaan.

**Solusi:**
- Gunakan `"definitions"` (bukan `"components/schemas"`) secara konsisten di seluruh spec gabungan karena spec menggunakan `"swagger": "2.0"`.
- Pastikan `securityDefinitions` di-merge dari tiap service di dalam `AggregateSpecs`.
- Pastikan `rewriteRefs` menulis ulang `$ref` ke `#/definitions/`, bukan `#/components/schemas/`.

---

### 5. Urutan tag dropdown tidak alfabetikal

**Penyebab:** Swagger UI mengurutkan tag berdasarkan urutan kemunculan pertama di `paths`, yang non-deterministic karena fetch dilakukan secara concurrent.

**Solusi:** Tambahkan opsi `tagsSorter: 'alpha'` pada konfigurasi Swagger UI:

```javascript
SwaggerUIBundle({
    url: "/swagger/aggregated.json",
    dom_id: '#swagger-ui',
    tagsSorter: 'alpha'
})
```
