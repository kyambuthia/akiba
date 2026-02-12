# Akiba - Fintech Banking Foundation

Monorepo foundation for a fintech/banking identity service.

## Why Expo for React Native
Expo is used for this foundation because it provides the fastest cross-platform scaffold for iOS, Android, and web/desktop-style targets with minimal native setup overhead.

## Structure
- `backend/`: Go API (clean architecture style)
- `client/`: Expo React Native TypeScript scaffold
- `docker-compose.yml`: local Mongo + backend runtime
- `Makefile`: common commands

## Backend Architecture
- `cmd/api`: app entrypoint
- `internal/domain`: core entities and domain validation helpers
- `internal/repository`: repository interfaces
- `internal/usecase`: business services
- `internal/infrastructure/mongo`: Mongo repo + idempotent index bootstrap
- `internal/transport/http`: handlers, middleware, router, error format
- `internal/auth`: JWT issue/verify helpers
- `internal/config`: env config loader
- `internal/observability`: structured logger

Domain code does not depend on Mongo or HTTP packages.

## Prerequisites
- Docker + Docker Compose
- Go 1.22+
- Node 18+
- GitHub CLI (`gh`) authenticated

## Environment
```bash
cp .env.example .env
```

Set a secure `JWT_SECRET` for non-dev environments.

### Backend env vars
- `ENV` (`development`)
- `PORT` (`8080`)
- `MONGO_URI` (`mongodb://mongo:27017`)
- `MONGO_DB_NAME` (`akiba`)
- `JWT_SECRET`
- `JWT_ISSUER` (`akiba-api`)
- `ACCESS_TOKEN_TTL` (`1h`)
- `DB_TIMEOUT` (`5s`)

## Run Local
```bash
make up
```

Stop:
```bash
make down
```

Run client scaffold:
```bash
cd client
npm install
npm run start
```

## Tests
```bash
make test
```

## API
- `POST /api/v1/auth/signup`
- `POST /api/v1/auth/login`
- `GET /api/v1/me` (Bearer token)
- `GET /health` (liveness)
- `GET /ready` (readiness)

### Curl examples
Signup:
```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","phone":"+14155552671","username":"user_1","password":"Password1"}'
```

Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"user_1","password":"Password1"}'
```

Me:
```bash
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

Error shape:
```json
{
  "error": {
    "code": "validation_error",
    "message": "invalid signup payload",
    "fields": { "phone": "must be valid E.164 format" }
  }
}
```

## Security Defaults
- bcrypt password hashing
- HS256 access-token JWT
- never logs passwords/tokens
- UTC timestamps
- request ID + recover + structured request logs
- strict JSON parsing (unknown fields rejected) + 1MB request body limit
- idempotent startup index creation on `users` collection (`emailLower`, `phoneE164`, `usernameLower`)
