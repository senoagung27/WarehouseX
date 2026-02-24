# 📌 PRODUCT REQUIREMENT DOCUMENT (PRD)
# WarehouseX – Enterprise Inventory Management System

Version: 1.0  
Author: Engineering Team  
Date: 2026  

---

# 1. Executive Summary

WarehouseX adalah sistem inventory enterprise yang dirancang untuk:

- Mengelola inbound & outbound warehouse
- Mengimplementasikan approval workflow berbasis RBAC
- Menjamin data integrity menggunakan ACID transaction
- Mencegah race condition menggunakan Redis Distributed Lock
- Menyediakan audit trail untuk compliance

Target sistem: Enterprise warehouse, distribution center, supply-chain company.

---

# 2. Problem Statement

Dalam sistem inventory tradisional:

- Terjadi race condition saat multiple outbound
- Stock menjadi negatif
- Tidak ada approval workflow
- Tidak ada audit trail
- Tidak scalable

WarehouseX menyelesaikan masalah tersebut dengan concurrency-safe architecture.

---

# 3. Product Goals

| Goal | Description |
|------|------------|
| Stock Integrity | Tidak ada negative stock |
| Concurrency Safety | Zero race condition |
| Compliance | Full audit trail |
| Scalability | Horizontal-ready |
| Security | RBAC + JWT |

---

# 4. User Persona

## 4.1 Warehouse Staff
- Create inbound/outbound request
- View stock

## 4.2 Supervisor
- Approve / Reject request

## 4.3 Admin
- Full access
- Override approval

## 4.4 Auditor
- Read-only access
- Export report

---

# 5. Functional Requirements

## 5.1 Inventory Management

- Create item
- Update item
- View stock
- Stock history

---

## 5.2 Inbound Workflow

1. Staff create inbound request
2. Status = PENDING
3. Supervisor approve
4. Stock incremented

---

## 5.3 Outbound Workflow

1. Staff create outbound request
2. Validate stock
3. Acquire Redis distributed lock
4. Deduct stock in DB transaction
5. Insert audit log
6. Release lock

---

## 5.4 Approval State Machine

PENDING → APPROVED → COMPLETED  
PENDING → REJECTED  

Constraints:
- Only Supervisor/Admin can approve
- Cannot approve twice
- Approval logged

---

# 6. Non-Functional Requirements

| Requirement | Target |
|------------|--------|
| API Response Time | < 200ms |
| Availability | 99.9% |
| Concurrency | 1000 concurrent users |
| Consistency | ACID guaranteed |
| Security | JWT-based auth |

---

# 7. Risk & Mitigation

| Risk | Mitigation |
|------|------------|
| Race condition | Redis distributed lock |
| DB deadlock | Row-level locking (FOR UPDATE) |
| Double submission | Idempotency key |
| Privilege abuse | Strict RBAC |

---

# 8. Success Metrics

- No negative stock
- No double approval
- No lost transaction
- Successful load test with 1000 concurrent requests

---

# 9. Roadmap

## v1
- Inventory
- Inbound/Outbound
- Approval
- Redis Lock
- Audit Log

## v2
- Multi-warehouse
- Multi-tenant
- Report export

## v3
- Event-driven architecture
- Kafka integration
- Real-time dashboard

---