# Infra Activity Log

Internal web app to replace the team's Excel activity/change log sheet. Single Go binary (API + embedded frontend), PostgreSQL backend, image attachments for Old Value / New Value fields.

## Stack

- Go 1.22, chi router, pgx (PostgreSQL driver)
- Frontend: vanilla HTML/JS embedded into the binary via `go:embed` — no separate frontend build/deploy
- Postgres 16 (Docker)
- No authentication (internal network only — do not expose port 8080 publicly)

## Local setup

```bash
cp .env.example .env
# edit .env — set a real POSTGRES_PASSWORD

docker compose up --build -d
```

First run: Postgres auto-applies `migrations/0001_init.sql` via `docker-entrypoint-initdb.d` (only runs on a fresh, empty volume — if you change the schema later, add a new migration file and apply manually).

App available at: `http://<server-ip>:8080`

## API reference

| Method | Path | Description |
|---|---|---|
| GET | `/api/logs?date_from=&date_to=&pic=&status=&category=&search=` | List logs with filters |
| GET | `/api/logs/{id}` | Get single log |
| POST | `/api/logs` | Create (multipart form) |
| PUT | `/api/logs/{id}` | Update (multipart form) |
| DELETE | `/api/logs/{id}` | Delete |
| GET | `/uploads/{filename}` | Serve uploaded image |
| GET | `/healthz` | Health check |

### Example: create a log entry

```bash
curl -X POST http://localhost:8080/api/logs \
  -F "tanggal=2026-08-09" \
  -F "job_title=Network Engineer" \
  -F "pic=NETWORK" \
  -F "application=Core Switch" \
  -F "label=Perapian wlan & ap group" \
  -F "old_value_text=AP group tidak sinkron" \
  -F "new_value_text=AP group sudah sinkron" \
  -F "status=Done" \
  -F "category=Daily" \
  -F "old_value_image=@/path/to/before.png" \
  -F "new_value_image=@/path/to/after.png"
```

### Example: filter

```bash
curl "http://localhost:8080/api/logs?pic=DCO&status=Done&category=Daily"
```

## Valid enum values

- `pic`: `DBA`, `DCO`, `NETWORK`
- `status`: `Open`, `Process`, `Hold`, `Done`
- `category`: `Change`, `Daily`, `Incident`

## CI/CD

`.github/workflows/ci.yml`:
- On every push/PR: `go vet`, build check, `go test`
- On push to `main`: builds and pushes Docker image to `ghcr.io/<owner>/<repo>:latest` and `:<sha>`

**To auto-deploy to this server after image push**, add a deploy job to the workflow using SSH (e.g. `appleboy/ssh-action`) that runs `docker compose pull && docker compose up -d` on this host. Not included by default — set up a deploy key/secret first if you want this.

## Backup

Data lives in two Docker named volumes: `pgdata` (Postgres) and `uploads` (images). Back these up regularly, e.g.:

```bash
docker run --rm -v infra-activity-log_pgdata:/data -v $(pwd):/backup alpine \
  tar czf /backup/pgdata-backup.tar.gz -C /data .
```

## Notes

- Max image upload size: 5MB, jpg/png only.
- No auth — restrict network access (firewall/VPN) since this is internal-only by design.
