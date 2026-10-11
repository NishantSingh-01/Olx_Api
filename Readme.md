# OLX API

A modular, lightweight REST API backend for an OLX-style classified marketplace platform built with Go.

---

## 🚀 Features

- **Standard Library Routing**: Built with Go 1.22+ enhanced `http.ServeMux` pattern matching.
- **Fail-Fast Configuration**: Strict environment loading using `godotenv` ensuring essential variables (`PORT`, `ENV`) are validated on startup.
- **Hardened HTTP Server**: Production-ready timeout configurations (`ReadTimeout`, `WriteTimeout`, `IdleTimeout`) preventing slowloris and resource leaks.
- **Profiling Enabled**: Integrated `net/http/pprof` for real-time memory and CPU diagnostics.
- **Makefile Automation**: Convenient commands for building, running, and cleaning binaries.

---

## 📁 Project Structure

```text
olx_api/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point & server setup
├── internal/
│   ├── config/
│   │   └── config.go        # Environment and config loader
│   └── handlers/
│       └── healthz.go       # Health check handler
├── .env                     # Local environment variables (ignored by git)
├── .gitignore               # Ignored files (bin, .env)
├── go.mod                   # Go module definitions
├── go.sum                   # Checksums for dependencies
├── Makefile                 # Build & run automation
└── Readme.md                # Project documentation
```

---

## 🛠️ Prerequisites

- **Go**: Version `1.22` or higher (configured for `1.27.0` in `go.mod`)
- **Make** (optional, for running Makefile targets on Windows/Unix)
- **Git**

---

## ⚙️ Environment Variables

Create a `.env` file in the root directory (or copy from [`.env.example`](file:///c:/Users/NISHANT%20SINGH/OneDrive/Desktop/olx_api/.env.example)):

```bash
cp .env.example .env
```

```env
PORT=7000
ENV=development
```

| Variable | Description | Required | Example |
| :--- | :--- | :---: | :--- |
| `PORT` | The port the HTTP server listens on | Yes | `7000` |
| `ENV` | Environment mode (`development`, `production`, etc.) | Yes | `development` |

---

## 🏃 Getting Started

### 1. Clone the repository
```bash
git clone https://github.com/NishantSingh-01/Olx_Api.git
cd Olx_Api
```

### 2. Install dependencies
```bash
go mod tidy
```

### 3. Run the application

#### Using `go run`:
```bash
go run ./cmd/api/main.go
```

#### Using `Makefile`:
```bash
# Run directly
make run

# Build executable to bin/api/main.exe
make build

# Start the built executable
make start

# Clean build artifacts
make clean
```

---

## 📡 API Endpoints

### Health Check

Verify the API status and server availability.

- **URL:** `/health`
- **Method:** `GET`
- **Response Code:** `200 OK`
- **Content-Type:** `application/json`

**Sample Request:**
```bash
curl -X GET http://localhost:7000/health
```

**Sample Response:**
```json
{
  "status": "ok",
  "message": "server is Healthy"
}
```

---

## 🌐 Deployment

The service can be containerized or hosted on platforms like **Render**, **Fly.io**, or **Railway**.

- **Live URL:** [https://olx-api-5wfw.onrender.com/health](https://olx-api-5wfw.onrender.com/health)

---

## 🗺️ Roadmap

- [ ] **Authentication & Authorization**: JWT-based user sign-up, sign-in, and profile management.
- [ ] **Product Listings**: CRUD operations for classified items with categories, pricing, and location.
- [ ] **Media Uploads**: Cloud storage integration (AWS S3 / Cloudinary) for product photos.
- [ ] **Search & Filtering**: Query by category, price range, location, and keywords.
- [ ] **Real-time Chat**: Buyer-seller messaging via WebSockets.
- [ ] **Database Integration**: PostgreSQL / MongoDB with migrations.
