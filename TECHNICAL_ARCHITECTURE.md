# 🧠 TECHNICAL ARCHITECTURE DOCUMENT
# WarehouseX

---

# 1. Clean Architecture Pattern

```
internal/
├── config/              # Configuration (env, JWT, DB, Redis)
├── model/               # Data models / entities (structs + business rules)
├── domain/
│   └── repository/      # Repository interfaces (contracts)
├── dto/                 # Data Transfer Objects (request/response)
├── controller/          # HTTP controllers (depend on service interfaces)
├── service/             # Service interfaces + business logic implementations
├── repository/          # Repository implementations (GORM)
├── infrastructure/      # External infrastructure (DB, Redis)
├── middleware/           # JWT Auth + RBAC middleware
└── router/              # Route definitions
```

**Dependency Flow:**
```
Controller → Service Interface → Repository Interface
                  ↓                      ↓
          Service (impl)         Repository (impl)
                  ↓                      ↓
             Model/DTO            Model + GORM/DB
```

---

# 2. Request Flow (Outbound)

1. Client request
2. JWT validated
3. RBAC checked
4. Redis lock acquired
5. DB transaction begin
6. SELECT FOR UPDATE
7. Validate stock
8. Update stock
9. Insert audit log
10. Commit
11. Release lock

---

# 3. Scalability Design

- Stateless API
- Redis externalized
- DB replication ready
- Horizontal scaling ready

---

# 4. Failure Scenario Handling

If Redis lock fails:
→ Return 409 Conflict

If DB transaction fails:
→ Rollback
→ Release lock

If crash during transaction:
→ PostgreSQL rollback automatically

---

# 5. Future Architecture Upgrade

- Event-driven with Kafka
- CQRS
- Microservices split:
  - Inventory Service
  - Approval Service
  - Audit Service

---