# backend (PoC)

Go API skeleton for alive-check events.

## endpoints
- `POST /v1/events`
- `POST /v1/snoozes`
- `GET /v1/devices/{id}/latest-status`

## run
```bash
cd apps/backend
go run ./cmd/server
```

Note: local machine currently does not have Go toolchain installed.
