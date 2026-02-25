# Akiba
Akiba is a monorepo for a fintech-ready authentication foundation: a Go API for signup/login/JWT auth plus an Expo React Native client scaffold for integrating the flows.

## Stack
- Backend: Go 1.22+, Chi, MongoDB, bcrypt, JWT (HS256)
- Client: Expo React Native + TypeScript (scaffold only)
- Infra: Docker, Docker Compose, Makefile
- CI: GitHub Actions (`go vet`, `go test`, `golangci-lint`, Docker build)

## Monorepo Layout
- `backend/` Go API (clean architecture)
- `client/` Expo app scaffold (`Signup`, `Login` placeholders)
- `docker-compose.yml` local Mongo + backend runtime
- `Makefile` common developer commands

## Quickstart
### Prerequisites
- Docker + Docker Compose
- Go 1.22+
- Node 18+
- GitHub CLI (`gh`) authenticated

### Environment
```bash
cp .env.example .env
```

Core backend vars:
- `ENV` (default `development`)
- `PORT` (default `8080`)
- `MONGO_URI` (default `mongodb://mongo:27017`)
- `MONGO_DB_NAME` (default `akiba`)
- `JWT_SECRET` (set secure value outside local dev)
- `JWT_ISSUER` (default `akiba-api`)
- `ACCESS_TOKEN_TTL` (default `1h`)
- `DB_TIMEOUT` (default `5s`)

### Run
```bash
make up
```

```bash
make down
```

### Test
```bash
make test
```

### Client Scaffold
```bash
cd client
npm install
npm run start
```

## API
Base path: `/api/v1`

- `POST /auth/signup`
- `POST /auth/login`
- `GET /me` (Bearer token)
- `GET /health` (liveness)
- `GET /ready` (readiness; Mongo ping)

### Auth Payloads
`POST /auth/signup`
```json
{
  "email": "user@example.com",
  "phone": "+14155552671",
  "username": "user_1",
  "password": "Password1"
}
```

`POST /auth/login`
```json
{
  "login": "user_1",
  "password": "Password1"
}
```

### Validation Rules
- `email`: valid format, normalized lowercase
- `phone`: E.164
- `username`: `^[a-zA-Z0-9_]{3,20}$`, normalized lowercase
- `password`: min 8, at least 1 letter and 1 number

### Error Contract
```json
{
  "error": {
    "code": "validation_error",
    "message": "invalid signup payload",
    "fields": {
      "phone": "must be valid E.164 format"
    }
  }
}
```

### Curl
```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","phone":"+14155552671","username":"user_1","password":"Password1"}'
```

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"user_1","password":"Password1"}'
```

```bash
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

## Architecture (Backend)
- `cmd/api` process bootstrap
- `internal/domain` core entities + validation primitives
- `internal/repository` repository interfaces
- `internal/usecase` business logic
- `internal/infrastructure/mongo` Mongo repository + idempotent index setup
- `internal/transport/http` handlers, middleware, router, response contract
- `internal/auth` JWT issue/verify
- `internal/config` env loader
- `internal/observability` structured logging

Design rule: domain layer has no HTTP or Mongo dependencies.

## Security and Runtime Defaults
- Password hashing with bcrypt
- JWT access tokens (HS256)
- UTC timestamps
- Request ID + panic recovery + structured request logs
- Strict JSON decoding (`DisallowUnknownFields`)
- Request body size limit: 1MB
- Idempotent startup indexes on `users`:
- `emailLower` unique
- `phoneE164` unique
- `usernameLower` unique

## OpenAPI
Skeleton spec: `backend/openapi.yaml`
