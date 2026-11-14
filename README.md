# Family Space

A secure, cloud-native family identity and collaboration platform built with a microservices architecture, featuring a production-ready authentication service as the foundation.

> ⚠️ **Project Status: Very Early Stage** — This project is in its initial phases of development. APIs, architecture, and documentation are subject to significant changes. 

## 📦 Repository Structure

This repository contains the **backend** of the Family Space project. For the frontend, see: **[family-space-frontend](https://github.com/Tata-Matata/family-space-frontend)**

- **Backend (this repo)**: Microservices, authentication, APIs, database schemas
- **Frontend**: Web UI for login, registration, and user-facing features

## Project Overview

Family Space is a platform that enables family members to securely authenticate, collaborate, and access shared resources. The architecture prioritizes security, scalability, and clear separation of concerns.

### Design Philosophy

- **Zero-trust security**: Verification at system boundaries (gateway), not within
- **Stateless verification**: Use of short-lived, stateless tokens for API authentication
- **Clean separation**: Auth concerns isolated in dedicated service; authorization delegated to gateway
- **Incremental delivery**: Phased approach, each phase complete and deployable
- **Cloud-native**: Kubernetes-ready from the start with security-first practices

## Project Status

### Phase 1: Auth Service ✅ (In Progress)
Functional authentication backend with user registration, login, and token management.

**Current Features:**
- User registration and login
- Password hashing (bcrypt)
- RS256-signed JWT access tokens (short-lived, stateless)
- Opaque refresh tokens (long-lived, server-side, revocable)
- Refresh token rotation and replay attack prevention
- Logout with refresh token revocation
- Transaction management abstraction for clean database handling
- Support for both SQLite and PostgreSQL

**Testing:**
- Unit tests for core logic (password hashing, JWT signing, services)
- Component-level API tests (request/response validation with fakes)
- Integration tests planned for transactional behavior

**Deployment:**
- Docker support
- Local setup with environment variables
- Kubernetes-ready (architecture prepared for Phase 4)

### Phase 2: Minimal UI (Planned)
Login and registration frontend with secure token handling.

### Phase 3: Production-like Local Setup (Planned)
Docker Compose orchestration with SQLite backend.

### Phase 4: Kubernetes Deployment (Planned)
Production deployment with pod security, network policies, RBAC, and secrets management.

### Phase 5: CI/CD (Planned)
GitHub Actions for testing and scanning; ArgoCD for GitOps deployments.

### Phase 6: Gateway & Microservices (Planned)
API Gateway for JWT verification and routing; additional microservices for features.

## Project Structure

```
family-space/
├── apps/                          # Microservices
│   ├── auth-service/              # Authentication service (Go)
│   │   ├── internal/
│   │   │   ├── auth/              # Auth logic (HTTP handlers, JWT, password, refresh tokens)
│   │   │   ├── domain/            # Domain models (User, Family, Membership)
│   │   │   ├── storage/           # Storage abstraction and implementations
│   │   │   └── errors/            # Error definitions
│   │   ├── main/                  # Application entry point
│   │   ├── migrations/            # Database migrations
│   │   ├── keys/                  # RSA key storage
│   │   ├── go.mod, go.sum         # Go dependencies
│   │   ├── Dockerfile             # Container image
│   │   └── README.md              # Service-specific documentation
│   ├── content-service/           # (Planned) Content management
│   ├── gateway/                   # (Planned) API Gateway
│   └── review-service/            # (Planned) Review/approval system
├── clusters/                      # Infrastructure
│   └── local/                     # Local development cluster (Kustomize)
├── docs/                          # Project documentation
│   ├── architecture.md            # High-level architecture and phases
│   └── trajectory.md              # Detailed feature roadmap (see Phase 1 section)
└── README.md                      # This file
```

## Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (for containerized setup)
- PostgreSQL or SQLite (for database)

### Running the Auth Service Locally

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Tata-Matata/family-space.git
   cd family-space
   ```

2. **Navigate to the auth service:**
   ```bash
   cd apps/auth-service
   ```

3. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   # Required: DB_DRIVER, DB_PATH or DATABASE_URL, JWT_PRIVATE_KEY_PATH, REFRESH_TOKEN_HMAC_KEY
   ```

4. **Generate RSA keys (if needed):**
   ```bash
   # Private key (keep secure)
   openssl genrsa -out keys/jwt_private.pem 2048
   
   # Public key (for gateway)
   openssl rsa -in keys/jwt_private.pem -pubout -out keys/jwt_public.pem
   ```

5. **Run the service:**
   ```bash
   go run ./main/main.go
   ```

   The auth service will start on `http://localhost:8080`

6. **Test the endpoints:**
   ```bash
   # Register a new user
   curl -X POST http://localhost:8080/register \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"secure_password"}'
   
   # Login
   curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"secure_password"}'
   
   # Health check
   curl http://localhost:8080/health
   ```

### Running Tests

```bash
cd apps/auth-service

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestLoginHandler ./internal/auth/http
```

### Running with Docker

```bash
cd apps/auth-service

# Build the image
docker build -t family-space-auth .

# Run the container
docker run -p 8080:8080 \
  -e DB_DRIVER=sqlite \
  -e DB_PATH=/app/auth.db \
  -e JWT_PRIVATE_KEY_PATH=/app/keys/jwt_private.pem \
  -e REFRESH_TOKEN_HMAC_KEY=your_secret_key \
  family-space-auth
```

## Authentication Flow

### Login

1. Client sends credentials to `POST /login`
2. Auth Service validates credentials and loads user membership
3. Auth Service issues:
   - **Access token** (short-lived JWT, ~15 min) → stored in memory
   - **Refresh token** (long-lived opaque token, ~30 days) → stored as HTTP-only cookie
4. Client uses access token for API requests

### Token Refresh

1. Client receives `401 Unauthorized` (access token expired)
2. Client sends refresh token to `POST /refresh`
3. Auth Service validates, rotates, and issues new tokens
4. Client retries original request with new access token

### Logout

1. Client calls `POST /logout` with refresh token
2. Auth Service revokes the refresh token
3. Client clears both tokens

## API Endpoints

| Method | Endpoint | Purpose | Request | Response |
|--------|----------|---------|---------|----------|
| POST | `/register` | Create new user account | `{email, password}` | `{user_id, family_id}` |
| POST | `/login` | Authenticate and get tokens | `{email, password}` | `{access_token, refresh_token}` |
| POST | `/refresh` | Refresh expired access token | `{refresh_token}` | `{access_token, refresh_token}` |
| POST | `/logout` | Revoke refresh token | `{refresh_token}` | `{success}` |
| GET | `/health` | Health check | — | `{status: "ok"}` |

For detailed request/response schemas, see [Auth Service README](./apps/auth-service/README.md).

## Technology Stack

### Auth Service
- **Language**: Go 1.25
- **Cryptography**: golang.org/x/crypto (bcrypt), golang-jwt/jwt (RS256)
- **Databases**: SQLite (dev), PostgreSQL (production-ready)
- **Testing**: testify, httptest
- **Containerization**: Docker

### Planned
- **Frontend**: React or similar (Phase 2)
- **Gateway**: Kong, Traefik, or custom (Phase 4)
- **Orchestration**: Kubernetes with Kustomize (Phase 4)
- **CI/CD**: GitHub Actions + ArgoCD (Phase 5)

## Security Considerations

### Currently Implemented

- ✅ Password hashing with bcrypt
- ✅ RS256 asymmetric JWT signing
- ✅ Short-lived stateless access tokens (~15 min)
- ✅ Long-lived revocable refresh tokens (~30 days)
- ✅ Refresh token rotation on every refresh
- ✅ Replay attack prevention via token rotation
- ✅ Transaction-level consistency for atomic operations
- ✅ Refresh token hashing at rest
- ✅ Database abstraction supporting multiple backends

### Planned

- 🔄 Refresh token reuse detection and anomaly alerting
- 🔄 Rate limiting (auth endpoints)
- 🔄 Secure HTTPS headers (CORS, CSP, etc.)
- 🔄 Audit logging (all auth events)
- 🔄 MFA support
- 🔄 OAuth/OIDC integration
- 🔄 Pod security policies and network segmentation (Kubernetes)
- 🔄 Image scanning (Trivy) in CI/CD


### Code Structure

- **handlers** (`internal/auth/http/`): HTTP request/response mapping
- **services** (`internal/auth/service/`): Business logic
- **storage** (`internal/storage/`): Data access abstractions and implementations
- **domain** (`internal/domain/`): Core domain models
- **jwt/password/refresh**: Cryptographic operations

## Documentation

- [**Architecture & Phases**](./docs/architecture.md) — High-level overview and development roadmap
- [**Detailed Trajectory**](./docs/trajectory.md) — In-depth feature specifications for each phase
- [**Auth Service README**](./apps/auth-service/README.md) — Detailed auth service design, token model, and testing strategy
- [**TODO List**](./apps/auth-service/TODO.md) — Current tasks and improvements

