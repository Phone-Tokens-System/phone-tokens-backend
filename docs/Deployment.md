# Deployment

## Production environment

Production `.env` создается вне репозитория. Секреты нельзя коммитить.

```env
POSTGRES_USER=<production-db-user>
POSTGRES_PASSWORD=<strong-production-db-password>
POSTGRES_DB=<production-db-name>
DATABASE_URL=postgres://<production-db-user>:<strong-production-db-password>@db:5432/<production-db-name>?sslmode=disable

HTTP_PORT=8080
JWT_SECRET=<strong-random-secret>
JWT_EXPIRES_IN_SEC=<token-lifetime-seconds>
FRONTEND_URL=https://<production-domain>

STRIPE_KEY=<stripe-secret-key>
WEBHOOK_SECRET=<stripe-webhook-secret>
SERVER_URL=https://<production-domain>

API_KEY=<sms-provider-api-key>
EMAIL=<provider-or-service-email>

CA_CERT_PATH=/secure/path/cert.pem
CA_KEY_PATH=/secure/path/key.pem
```

## Сборка

Frontend:

```bash
cd phone-tokens-frontend
npm ci
npm test
npm run build
```

Backend:

```bash
cd phone-tokens-backend
go test ./...
docker build -t phone-tokens-backend:<release-tag> .
```

## Docker Compose

```bash
cd phone-tokens-backend
docker compose up -d --build
docker compose ps
```

Ожидаемое состояние:

- `db` здоров по healthcheck;
- backend доступен внутри Docker network на `backend:8080`;
- PostgreSQL не имеет публичного host port;
- CA certificate/private key подключены как read-only runtime volumes;
- nginx публикует HTTP/HTTPS;
- `/` отдает frontend;
- `/api/` проксируется в backend.

## Nginx

Config file: `nginx/conf.d/app.conf`.

Nginx отвечает за:

- HTTP -> HTTPS redirect;
- TLS certificates из `/etc/letsencrypt`;
- frontend static files из `/usr/share/nginx/html`;
- reverse proxy `/api/` на `backend:8080`.
