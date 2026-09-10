# cursor-cli-to-api-api

OpenAI-compatible proxy that sits between clients and the Cursor CLI (`agent -p`).

## Prerequisites

1. Cursor Agent CLI installed and logged in: `agent login` (use `C:\Users\armin\AppData\Local\cursor-agent\agent.cmd` on this machine — PATH `agent` may be Grok).
2. Docker for Postgres.
3. Go 1.22+.

## Run

```powershell
docker compose up -d
copy .env.example .env
go run ./cmd/server
```

Listens on `:8201` by default.

## Auth

- Username: `armin`
- Password: `dopadopa123`
- HTTP Basic **or** `Authorization: Bearer dopadopa123`

## Endpoints

- `GET /health`
- `GET /v1/models`
- `POST /v1/chat/completions`
- `POST /v1/responses`
- `GET /v1/responses/{id}`
