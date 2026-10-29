# LibraTrack — Theme 8 (Docker, Compose, K3s)

Library API (Gin + PostgreSQL) with HTMX catalog. Homework: containerize, run with Compose, deploy to Kubernetes.

## Docker Compose

```bash
docker compose up --build -d
curl http://localhost:8080/health
open http://localhost:8080/catalog
```

## Kubernetes / K3s

See [deploy/README.md](deploy/README.md): K3s demo apps 1–2, LibraTrack manifests under `deploy/k8s/`.

## Local dev

```bash
cp .env.example .env
docker compose up -d postgres
export DATABASE_URL=postgres://libratrack:libratrack@localhost:5432/libratrack?sslmode=disable
export JWT_SECRET=dev-secret
export MIGRATIONS_PATH=migrations
go run ./cmd/api
```
