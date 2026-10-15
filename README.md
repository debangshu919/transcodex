# Transcodex

A media transcoding API service written in Go. Upload video files to S3, request format conversions (mp4 / mkv / webm), and let workers process them asynchronously. Includes API-key auth and a credit-based billing model.

## Status

Work in progress. The working pieces today:

- HTTPS API server with middleware (rate limiting, CORS, compression, security headers, HPP protection)
- S3-compatible file storage (upload, list, download via presigned URL, delete)
- ffmpeg/ffprobe converter library with NVIDIA NVENC hardware acceleration and live progress reporting

Still stubbed out:

- Async job pipeline (`cmd/worker`) — not implemented
- Auth endpoints and conversion endpoints — return "Not Implemented"
- Credit/billing logic — schema exists, no code yet

## Requirements

- Go 1.26+
- PostgreSQL (or Docker)
- ffmpeg and ffprobe on `PATH`
- An NVIDIA GPU with CUDA/NVENC support for actual conversions
- An S3-compatible bucket (AWS S3, Floci, MinIO, etc.)

## Getting Started

1. Start PostgreSQL:

   ```sh
   docker compose up -d
   ```

2. Create your env files. `.env.development` for local work:

   ```env
   PORT=3000
   ENV=development
   DATABASE_URL=postgresql://myuser:mysecretpassword@localhost:5432/transcodex?sslmode=disable
   LOG_LEVEL=info
   S3_BUCKET=your-bucket
   ```

   `.env` is used for production. Never commit real credentials.

3. Run migrations:

   ```sh
   make migrate-up
   ```

   Revert with `make migrate-down` or `make migrate-down STEPS=1`.

4. Start the server in dev mode:

   ```sh
   make dev
   ```

   The API listens on `https://localhost:3000` using the committed self-signed `cert.pem` / `key.pem`.

## Make Targets

| Target | Description |
| --- | --- |
| `make build` | Build the API server to `bin/api` |
| `make run` | Build and run in production mode |
| `make dev` | Build and run in development mode |
| `make test` | Run all Go tests |
| `make clean` | Remove the `bin/` directory |
| `make migrate-up [STEPS=n]` | Apply migrations (optionally n steps) |
| `make migrate-down [STEPS=n]` | Revert migrations (optionally n steps) |

## API Endpoints

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/healthz` | Health check |
| `GET` | `/formats` | List supported target formats |
| `GET` | `/formats/{format}` | Get details for one format |
| `POST` | `/files/upload` | Upload a file to S3 |
| `GET` | `/files` | List files |
| `GET` | `/files/{id}` | Redirect to a presigned download URL |
| `DELETE` | `/files/{id}` | Delete a file |
| `*` | `/auth/api-keys...` | Not implemented yet |
| `*` | `/conversions...` | Not implemented yet |

## Project Structure

```
cmd/
  api/        API server entry point
  converter/  CLI demo of the ffmpeg conversion pipeline
  migrate/    Database migration runner
  test/       S3 smoke-test tool
  worker/     Job worker (stub)
internal/
  api/        Router, handlers, middlewares
  config/     Env loading and validation
  converter/  ffprobe probing, format matrix, ffmpeg execution
  db/         PostgreSQL connection
  storage/    S3 client
migrations/   SQL migrations (users, api_keys, files, conversions, jobs, credits)
pkg/
  logger/     slog fan-out logger (stdout, app/error/debug log files)
  utils/      Middleware application helpers
```

## Notes

- The converter library targets NVENC (`-hwaccel cuda`, `av1_nvenc`) and will not work on machines without a compatible NVIDIA GPU without changes.
- Logs are written to `logs/app.log`, `logs/error.log`, and optionally `logs/debug.log` depending on `LOG_LEVEL`.
- The API uses a self-signed certificate for local development; browsers will warn on first visit.
