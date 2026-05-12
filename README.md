# KFC Sales Forecast

A fullstack service that generates daily sales predictions for KFC stores, helping cooks prepare products in advance and reduce waste.

## Architecture

```
Browser
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│                      Docker Network                     │
│                                                         │
│  ┌──────────────┐    ┌──────────────┐  ┌─────────────┐ │
│  │   Frontend   │    │   Backend    │  │  PostgreSQL  │ │
│  │  nginx :80   │──▶ │   Go :8080   │─▶│  postgres   │ │
│  │              │    │              │  │  :5432      │ │
│  │  React + TS  │    │  Gin router  │  │             │ │
│  │  Recharts    │    │  GORM        │  │  stores     │ │
│  │  Zustand     │    │  Cron job    │  │  products   │ │
│  └──────────────┘    └──────────────┘  │  hist_sales │ │
│   /api/* proxied                       │  forecasts  │ │
│   to backend                           └─────────────┘ │
└─────────────────────────────────────────────────────────┘
```

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.25, Gin, GORM, robfig/cron, Viper |
| Frontend | React 19, TypeScript, Vite, Tailwind CSS v4 |
| State | Zustand (persisted to localStorage) + TanStack Query |
| Charts | Recharts |
| Database | PostgreSQL 16 |
| Container | Docker + Docker Compose |

## Quick Start

**Prerequisites:** Docker Desktop

```bash
git clone <your-repo-url>
cd kfc-forecast
docker compose up --build
```

Open **http://localhost** — the app is ready.

On first boot the scheduler generates forecasts for tomorrow automatically. Click a store on the left, pick a date, and see hourly predictions per product.

## Configuration

All user-tunable settings live in `config.yaml` at the project root. Edit and restart to apply changes.

```yaml
server:
  port: 8080

database:
  host: postgres      # service name inside Docker
  port: 5432
  user: kfc
  password: kfc_secret
  name: kfc_forecast
  sslmode: disable

forecast:
  # Standard cron expression — when to run daily generation
  # Default: every day at 02:00 AM
  generation_cron: "0 2 * * *"

  # How many past days of sales to average over
  history_days: 30

  # How many days ahead to forecast (1 = tomorrow)
  days_ahead: 1
```

## API Reference

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Service health + DB ping |
| `GET` | `/api/stores` | List all stores |
| `GET` | `/api/forecasts?store_id=&date=` | Hourly forecasts for a store on a date |
| `POST` | `/api/forecasts/generate` | Manually trigger forecast generation |

### GET /api/forecasts

Query parameters:

| Param | Format | Example |
|---|---|---|
| `store_id` | integer | `1` |
| `date` | `YYYY-MM-DD` | `2026-05-13` |

Response shape:

```json
{
  "store": { "id": 1, "name": "KFC Tel Aviv Center", "location": "..." },
  "date": "2026-05-13",
  "forecasts": [
    {
      "product": { "id": 1, "name": "Zinger Burger", "category": "Burgers", "price": 39.90 },
      "hourly": [
        { "hour": 10, "predicted_quantity": 12.50 },
        { "hour": 11, "predicted_quantity": 9.30 }
      ]
    }
  ]
}
```

## Forecast Algorithm

```
For each store × product × hour (10–23):
  predicted_quantity = AVG(historical_sales.quantity)
  WHERE sale_date >= today - history_days
    AND sale_date <  today
  GROUP BY store_id, product_id, hour
```

Runs once per day on a configurable cron schedule. On startup, if no forecast exists for tomorrow, generation runs immediately. Each run is idempotent — existing rows for the target date are replaced atomically inside a transaction.

## Project Structure

```
kfc-forecast/
├── config.yaml              # All user-tunable settings
├── docker-compose.yml       # Orchestrates postgres + backend + frontend
├── docker/
│   └── init.sql             # Schema + seed data (runs on first DB boot)
├── backend/
│   ├── cmd/main.go          # Entry point
│   ├── internal/
│   │   ├── config/          # Viper config loader
│   │   ├── db/              # GORM connection with retry
│   │   ├── models/          # Store, Product, HistoricalSale, Forecast
│   │   ├── forecast/        # Avg algorithm + idempotent generation
│   │   ├── scheduler/       # robfig/cron wrapper + startup check
│   │   └── server/          # Gin router + handlers
│   └── Dockerfile
└── frontend/
    ├── src/
    │   ├── api/             # Axios client + typed query functions
    │   ├── store/           # Zustand store (persisted to localStorage)
    │   ├── components/      # Header, StoreList, DatePickerInput,
    │   │                    # ForecastCard, ForecastPanel
    │   └── pages/           # DashboardPage
    ├── nginx.conf           # SPA fallback + /api proxy
    └── Dockerfile
```

## Seed Data

The database is pre-loaded with:
- **3 stores** — KFC Tel Aviv Center, KFC Jerusalem Downtown, KFC Haifa Bay
- **8 products** — Zinger Burger, Crunchy Burger, Bucket 6/10pc, Fries Regular/Large, Coleslaw, Pepsi
- **30 days of historical sales** — hourly data (10:00–23:00) with realistic demand patterns and ±30% noise
