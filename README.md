# 🏭 WarehouseX

Enterprise Inventory Management System built with Go, featuring concurrency-safe stock operations, RBAC-based approval workflows, and full audit trail.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.22 |
| Framework | Gin |
| Database | PostgreSQL 16 |
| Cache/Lock | Redis 7 |
| ORM | GORM |
| Auth | JWT + bcrypt |
| Logging | Zap |

## Architecture

Clean Architecture pattern:
```
├── cmd/api/          # Entry point
├── internal/
│   ├── config/       # Environment configuration
│   ├── domain/       # Entities & interfaces
│   ├── repository/   # GORM implementations
│   ├── service/      # Business logic
│   ├── handler/      # HTTP handlers (Gin)
│   ├── middleware/    # JWT + RBAC
│   └── infrastructure/ # Database + Redis
├── migrations/       # SQL migrations
```

## Quick Start

### Prerequisites
- Go 1.22+
- Docker & Docker Compose

### Run with Docker
```bash
docker-compose up -d
```

### Run locally
```bash
# Start dependencies
docker-compose up -d postgres redis

# Run migrations (manual via psql or auto-migrate)
psql -h localhost -U warehousex -d warehousex -f migrations/000001_init.up.sql

# Start API
cp .env.example .env
go run ./cmd/api
```

### Health Check
```bash
curl http://localhost:8080/health
```

## API Endpoints

### Auth (Public)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/register` | Register user |
| POST | `/api/v1/auth/login` | Login, get JWT |

### Inventory (Protected)
| Method | Path | Role | Description |
|--------|------|------|-------------|
| GET | `/api/v1/inventory` | Staff+ | List items |
| GET | `/api/v1/inventory/:id` | Staff+ | Get item |
| POST | `/api/v1/inventory` | Admin | Create item |
| PUT | `/api/v1/inventory/:id` | Admin | Update item |

### Requests (Protected)
| Method | Path | Role | Description |
|--------|------|------|-------------|
| POST | `/api/v1/requests/inbound` | Staff+ | Create inbound |
| POST | `/api/v1/requests/outbound` | Staff+ | Create outbound |
| GET | `/api/v1/requests` | Staff+ | List requests |
| GET | `/api/v1/requests/:id` | Staff+ | Get request |
| PUT | `/api/v1/requests/:id/approve` | Supervisor/Admin | Approve |
| PUT | `/api/v1/requests/:id/reject` | Supervisor/Admin | Reject |

### Audit Logs (Protected)
| Method | Path | Role | Description |
|--------|------|------|-------------|
| GET | `/api/v1/audit-logs` | Auditor+ | View logs |

## Key Features

- **Concurrency Safety**: Redis distributed lock + PostgreSQL `SELECT FOR UPDATE`
- **RBAC**: 4 roles (Staff, Supervisor, Admin, Auditor) with hierarchical permissions
- **Approval Workflow**: State machine (PENDING → APPROVED → COMPLETED / REJECTED)
- **Audit Trail**: Full JSONB before/after logging on all mutations
- **Stock Integrity**: `CHECK (quantity >= 0)` constraint, no negative stock