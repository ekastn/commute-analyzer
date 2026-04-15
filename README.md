# Commute Analyzer

Aplikasi web untuk menghitung dan mengelola estimasi **biaya tahunan** serta **waktu yang terbuang** untuk perjalanan pulang-pergi (commute) ke kantor.

## Tech Stack

| Layer | Teknologi |
|-------|-----------|
| Backend | Go (Gin) |
| Frontend | React 19, TypeScript, Vite, React-Leaflet, Tailwind CSS v4 |
| Database | PostgreSQL + PostGIS |
| Routing | OpenRouteService (ORS) API |
| ORM | sqlc |
| Migration | goose |
| Container | Docker Compose |

## Fitur

- **Buat Rute Baru** — Pilih lokasi rumah dan kantor langsung di peta, isi detail kendaraan dan harga BBM
- **Daftar Commute** — Lihat semua rute tersimpan di sidebar
- **Edit Rute** — Ubah nama, kendaraan, harga BBM, hari kerja, atau lokasi tanpa perlu hapus dan buat ulang
- **Hapus Rute** — Hapus rute yang tidak diperlukan
- **Auto-recalculate** — Biaya tahunan dan waktu dihitung ulang otomatis saat data berubah
- **Map Interaktif** — Leaflet map dengan marker rumah/kantor dan polyline rute

## Cara Kerja

### Perhitungan Biaya

```
Biaya harian  = (jarak_pulang_pergi × efisiensi / 100) × harga_BBM_per_liter
Biaya tahunan = biaya_harian × hari_kerja_per_minggu × 52.14

Efisiensi:
  - Mobil    : 10 km/liter
  - Motor   : 2.5 km/liter
```

### User Identification

Setiap browser otomatis dibuatkan UUID dan disimpan di `localStorage['device_id']`. UUID ini dikirim ke backend sebagai `device_id` untuk mengikat data ke perangkat. Tidak ada authentication yang diperlukan.

### Routing

Menggunakan [OpenRouteService](https://openrouteservice.org/) untuk menghitung rute jalan. Profile routing:
- Mobil → `driving-car`
- Motor → `cycling-regular`

## Struktur Proyek

```
commute-analyzer/
├── cmd/
│   ├── api/
│   │   └── main.go           # Entry point API server (Gin)
│   └── db/
│       ├── migrations/
│       │   └── 00001_initial_schema.sql  # goose migration
│       └── queries/
│           ├── commute.sql   # sqlc query definitions
│           └── user.sql
├── internal/
│   ├── config/       # Env config
│   ├── dto/         # Request/response DTOs
│   ├── env/         # Env helpers
│   ├── handler/     # HTTP handlers
│   ├── response/    # Standardized API responses
│   ├── service/    # Business logic + ORS client
│   └── store/       # Generated sqlc code + models
├── web/
│   └── src/
│       ├── components/  # React components
│       ├── contexts/    # AuthContext (device_id)
│       ├── hooks/       # useCommutes, useIsMobile
│       ├── lib/         # API client, types
│       └── services/    # API service layer
├── compose.yaml           # Docker Compose
├── Dockerfile.api         # Backend container
├── Dockerfile.web         # Frontend container
├── sqlc.yaml             # sqlc config
└── go.mod
```

## API Endpoints

Base URL: `GET/POST/PATCH/DELETE http://localhost:8080/api/v1`

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/commutes?device_id=xxx` | Daftar semua commute user |
| `POST` | `/commutes` | Buat commute baru |
| `PATCH` | `/commutes/:id` | Update commute |
| `DELETE` | `/commutes/:id` | Hapus commute |
| `GET` | `/health` | Health check |

### Create Commute

```json
POST /api/v1/commutes
{
  "device_id": "browser-uuid",
  "name": "Rute Pagi",
  "home_lat": -6.200000,
  "home_lng": 106.816666,
  "office_lat": -6.175000,
  "office_lng": 106.828333,
  "vehicle": "car",
  "fuel_price": 10000,
  "days_per_week": 5
}
```

### Update Commute

```json
PATCH /api/v1/commutes/:id
{
  "name": "Rute Baru",
  "vehicle": "motorcycle",
  "fuel_price": 15000,
  "days_per_week": 6,
  "home_lat": -6.190000,
  "home_lng": 106.820000,
  "office_lat": -6.170000,
  "office_lng": 106.835000
}
```

## Setup

### Prerequisites

- Go 1.24+
- Node.js 20+
- pnpm
- Docker & Docker Compose
- ORS API key

### 1. Clone & Environment

```bash
git clone https://github.com/ekastn/commute-analyzer
cd commute-analyzer

cp .env.example .env   # isi DATABASE_URL dan ORS_API_KEY
```

### 2. Jalankan dengan Docker Compose

```bash
docker compose up -d
```

- Web UI: http://localhost:3000
- API: http://localhost:8080
- Database: localhost:5432

### 3. Setup Database Migration

```bash
# Buat database
docker compose exec db psql -U postgres -c "CREATE DATABASE commute_analyzer;"

# Jalankan migration
docker compose exec db psql -U postgres -d commute_analyzer -c "CREATE EXTENSION IF NOT EXISTS postgis;"
docker compose exec db psql -U postgres -d commute_analyzer -f /docker-entrypoint-initdb.d/00001_initial_schema.sql
```

Atau dengan goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir cmd/db/migrations postgres "postgres://user:pass@localhost:5432/db?sslmode=disable" up
```

### 4. Development (Tanpa Docker)

```bash
# Backend
cp .env.example .env   # isi DATABASE_URL dan ORS_API_KEY
go run cmd/api/main.go

# Frontend
cd web
pnpm install
pnpm dev
```

## Environment Variables

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| `DATABASE_URL` | — | PostgreSQL connection string |
| `ORS_API_KEY` | — | OpenRouteService API key |
| `SRV_ADDR` | `:8080` | Server bind address |
| `SRV_ENV` | `dev` | `dev` = load `.env` file |

## Database Schema

```sql
users (user_id UUID, device_id TEXT UNIQUE, created_at)
commutes (
  id, user_id, name, home_point GEOMETRY(POINT,4326),
  office_point GEOMETRY(POINT,4326),
  route_geometry GEOMETRY(LINESTRING,4326),
  distance_km, duration_min, vehicle, fuel_price,
  days_per_week, annual_cost, annual_minutes,
  created_at, updated_at
)
```
