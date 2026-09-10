# cursor-cli-to-api

OpenAI-compatible API + WebUI that proxies requests to the **Cursor CLI** (`agent -p`).

## Layout

| Path | Role |
|------|------|
| `cursor-cli-to-api-api/` | Gin + PostgreSQL OpenAI-shaped proxy |
| `cursor-cli-to-api-webui/` | Vue playground + docs (PWA) |

## Quick start

1. Log in to Cursor Agent on this machine (`C:\Users\armin\AppData\Local\cursor-agent\agent.cmd login`). Do **not** rely on PATH `agent` if it points at another tool.
2. API:

```powershell
cd cursor-cli-to-api-api
docker compose up -d
copy .env.example .env
go run ./cmd/server
```

3. WebUI:

```powershell
cd cursor-cli-to-api-webui
npm install
npm run dev
```

- API: `http://127.0.0.1:8201`
- WebUI: `http://localhost:5201`
- Auth: `armin` / `dopadopa123` (Basic or Bearer password)

## Client usage

Point any OpenAI-compatible client at `http://127.0.0.1:8201/v1` with API key `dopadopa123`.

Endpoints:

- `POST /v1/chat/completions`
- `POST /v1/responses`
- `GET /v1/responses/{id}`
- `GET /v1/models`
