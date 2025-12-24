# Auth Service

## Phase 1 – Auth Service (Backend, Functional)

### Features

- User registration
- User login
- Password hashing
- JWT issuance (RS256)
- Transactionally correct logic
- Clean separation of concerns

### Unit tests

- Component-level API tests

### Endpoints

- POST /register
- POST /login
- GET /health
- POST /logout
- POST /refresh
- (optional) GET /metrics

### Access Tokens

#### Stateless JWT

Short-lived (e.g. 10–15 minutes)
Contains:
- iss, aud
- sub (user id)
- family_id
- role
- exp, iat

### Refresh Tokens

- Opaque tokens (not JWT)
- Stored server-side (database)
- Long-lived (e.g. 7–30 days)
- One refresh token per session/device

### Token Storage

- refresh_tokens table
id
user_id
token_hash
expires_at
revoked_at
created_at

Refresh tokens are hashed at rest (same principle as passwords).

### Refresh Flow

- Client calls POST /refresh
- Sends refresh token
- Auth service:
validates token exists
checks expiration
checks revocation
issues new access token
rotates refresh token

Old refresh token is revoked. This prevents replay attacks.

### Logout

- POST /logout
- Refresh token is revoked server-side

### Forced Revocation

- Password change
- Account deactivation
- Admin action (future)

### Revocation Strategy
- Access tokens remain stateless
- Refresh tokens are checked against DB
- Revocation is immediate for refresh tokens
- Access tokens naturally expire

### Security (Backend)

- Password hashing (bcrypt / argon2)
- JWT hardening
issuer (iss), audience (aud), short expiration, strong key management

- Rate limiting (basic, in-service)
- Secure HTTP headers
- Audit logs (auth events only)
- Refresh token hashing
- Token rotation on refresh
- Token revocation tracking
- Audit logs
    - login success
    - login failure
    - token refresh
    - logout
    - revocation events
Audit logs are append-only and never expose secrets.

## Phase 2 – Minimal UI (Auth Only)
### UI Features

- Login page
- Registration page
- Logged-in confirmation screen
- Basic error handling
- Logout button

### Token Handling (Frontend)

- Access token: memory
- Refresh token: HTTP-only cookie (recommended)
- Silent refresh via /refresh

## Phase 3 – Local Production-like Setup
### Features

- Dockerfile for auth-service
- Dockerfile for UI
- Docker Compose setup
- Environment-based configuration
- SQLite (file-based)

## Phase 4 – Kubernetes Deployment (Security-First)
###  Kubernetes Features
- Pod Security
- PodSecurity (restricted profile)
- Non-root containers
- Read-only root filesystem
- Networking
- NetworkPolicies
- Ingress (HTTP routing)
- TLS termination
- Configuration & Secrets
- Secrets (JWT private key, DB config)
- ConfigMaps (non-sensitive config)
- Identity & Access
- RBAC (least privilege)
- ServiceAccounts per workload

## Phase 5 – CI/CD (Security-Oriented)
### CI (GitHub Actions)

- Build backend image
- Run unit tests
- Run component API tests
- Image scanning (Trivy / similar)
- Fail pipeline on critical vulnerabilities

### CD (GitOps)

- ArgoCD
- Declarative manifests
- Automated deploy on merge
- Environment separation (dev / prod-style)

## Phase 6 – Gateway Integration (Later)
### Scope

Introduce a Gateway after Auth is complete.

### Responsibilities

- JWT verification
- TLS
- Routing
- Claim extraction
- Auth service remains unchanged.

Refresh tokens

OAuth / OIDC

MFA

Token revocation

Social login

Complex UI styling