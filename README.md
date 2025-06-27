# 🏦 Amartha Loan Service

A RESTful API to simulate a **Loan Engine** system with multi-state transitions (`proposed → approved → invested → disbursed`), complete with state validation, data relationships, clean architecture, Swagger documentation, and PostgreSQL integration via Docker.

---

## ✅ Features

- Create new loan (initial status: `proposed`)
- Approve loan with proof of photo and field officer
- Investors can contribute partially until loan is fully funded
- Status automatically changes to `invested` when fully funded
- Simulated investor email is sent once loan is fully funded
- Disburse loan with agreement letter and field officer
- Auto-generated Swagger documentation (`/swagger/index.html`)
- Clean code & architecture structure
- Integration tests directly against PostgreSQL
- Fully Dockerized for local and CI setup

---

## 🧱 Tech Stack

- **Golang** 1.23.x
- **PostgreSQL** 13
- **Gin** (HTTP framework)
- **Swaggo** (Swagger generator)
- **Docker & Docker Compose**
- **Clean Architecture** pattern

---

## 📦 Project Structure

loan-service/
├── cmd/ # Main entry point
├── configs/ # Configuration loader
├── delivery/
│ └── http/ # HTTP handlers & routes
├── domain/ # Entities and enums
├── repository/
│ ├── interface.go # Interface definitions
│ └── postgres/ # PostgreSQL implementations
├── usecase/ # Business logic
├── utils/ # Utilities (e.g., dummy email)
├── docs/ # Auto-generated Swagger files
├── migrations/ # SQL schema setup
├── tests/ # Integration tests (direct DB)
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md


---

## 📥 API Endpoints

| Method | Endpoint                   | Description                 |
|--------|----------------------------|-----------------------------|
| POST   | `/v1/loans`                | Create a new loan           |
| POST   | `/v1/loans/{id}/approve`   | Approve a loan              |
| POST   | `/v1/loans/{id}/invest`    | Add investment to a loan    |
| POST   | `/v1/loans/{id}/disburse`  | Disburse an approved loan   |
| GET    | `/v1/loans/{id}`           | Retrieve loan details       |

Swagger UI: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 🚀 Running Locally

```bash
# 1. Build and start the service
docker-compose up --build

# 2. Apply database schema
docker exec -i loan-service-db-1 psql -U postgres -d loan_db < migrations/init.sql

---

## 🛠 Makefile Commands

The following handy shortcuts are available via the provided `Makefile`:

| Command           | Description                                      |
|-------------------|--------------------------------------------------|
| `make docker-rebuild`   | Rebuild containers from scratch and start up       |
| `make swagger`          | Regenerate Swagger docs into `/docs` folder        |
| `make test`             | Run all tests in `/tests` directory                |
| `make test-integration` | Run integration tests against local PostgreSQL     |

> ℹ️ Integration tests assume your local DB is accessible with env:
>
> `DB_HOST=localhost`, `DB_NAME=loan_db`

---

## 🧪 Example Usage

```bash
# Clean rebuild + run
make docker-rebuild

# Regenerate Swagger docs
make swagger

# Run all tests
make test

# Run integration tests directly on DB
make test-integration
