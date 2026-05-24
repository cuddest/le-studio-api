# Le Studio Contrology API

Production-oriented REST API for Le Studio Contrology using Go, Gin, GORM, and PostgreSQL.

## Architecture

```
Client -> Gin Handlers -> Services -> Repository Interfaces -> GORM Postgres
                         |-> Response Envelope
                         |-> Validator
                         |-> JWT + Refresh Token Store
```

## Discovery Notes

Public website confirms women-first branding, coaches roster, weekly planning, and FAQ-driven booking journey. Admin app login and route shell were accessible; authenticated screens are implemented from the provided specification as source of truth.

## Local Setup

1. Copy `.env.example` values into your environment.
2. Start database: `docker compose up -d postgres`
3. Install dependencies: `go mod tidy`
4. Run API: `make run`
5. Health check: `curl http://localhost:8080/healthz`

## Environment Variables

- `PORT`
- `ENV`
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `DB_SSLMODE`
- `JWT_SECRET`
- `ACCESS_TOKEN_DURATION`
- `REFRESH_TOKEN_DURATION`
- `CLOUDINARY_CLOUD_NAME`
- `CLOUDINARY_API_KEY`
- `CLOUDINARY_API_SECRET`
- `ALLOWED_ORIGINS`

## Sample cURL

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"first_name":"Sara","last_name":"N","email":"sara@example.com","password":"password123"}'

curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"sara@example.com","password":"password123"}'

curl http://localhost:8080/api/v1/pack-templates

curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Authorization: Bearer <access_token>" \
  -H 'Content-Type: application/json' \
  -d '{"slot_id":1,"user_pack_id":1}'
```
