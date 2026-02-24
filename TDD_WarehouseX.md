# 🏗 TECHNICAL DESIGN DOCUMENT (TDD)
# WarehouseX

Version: 1.0  
Language: Go 1.22  

---

# 1. System Architecture

Client
  ↓
API (Gin)
  ↓
Middleware (JWT + RBAC)
  ↓
Service Layer
  ↓
Redis (Distributed Lock)
  ↓
PostgreSQL (ACID Transaction)

---

# 2. Technology Stack

| Layer | Technology |
|-------|------------|
| Language | Go |
| Framework | Gin |
| DB | PostgreSQL 16 |
| Cache | Redis 7 |
| ORM | SQLX / GORM |
| Auth | JWT |
| Logging | Zap |
| Migration | golang-migrate |
| Test | Testify |

---

# 3. Database Schema

## users

- id UUID
- name
- email
- role
- password_hash
- created_at

## inventory

- id UUID
- item_name
- quantity INT
- version INT
- created_at
- updated_at

## requests

- id UUID
- type (INBOUND/OUTBOUND)
- status
- item_id
- quantity
- created_by
- approved_by
- created_at

## audit_logs

- id UUID
- entity
- entity_id
- action
- user_id
- before_value JSONB
- after_value JSONB
- created_at

---

# 4. Concurrency Strategy

## 4.1 Redis Distributed Lock

SET lock:item:<id> value NX EX 10

Purpose:
- Prevent concurrent stock mutation
- Cross-instance safe

---

## 4.2 PostgreSQL ACID Transaction

BEGIN;

SELECT quantity FROM inventory WHERE id=$1 FOR UPDATE;

UPDATE inventory SET quantity = quantity - $1 WHERE id=$2;

INSERT INTO audit_logs ...

COMMIT;

---

# 5. RBAC Middleware

JWT contains:
- user_id
- role

Middleware:
- Validate token
- Validate role-permission mapping

---

# 6. Error Handling Strategy

| Error | Response |
|-------|----------|
| Insufficient stock | 400 |
| Unauthorized | 403 |
| Lock failed | 409 |
| Internal error | 500 |

---

# 7. Logging Strategy

Structured logging:
- request_id
- user_id
- action
- latency
- error

---

# 8. Testing Strategy

## Unit Test
- Service layer
- Mock repository

## Integration Test
- Real Postgres
- Real Redis
- Test outbound transaction

## Load Test
- Simulate concurrent outbound

---

# 9. Deployment Strategy

Docker containerized  
Docker Compose for local  
Kubernetes ready  

---

# 10. Security Consideration

- Password hashed (bcrypt)
- HTTPS only
- Rate limiting
- JWT expiration
- RBAC enforcement

---