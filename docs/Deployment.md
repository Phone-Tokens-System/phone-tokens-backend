# Deployment

## Production environment

Production `.env` создается рядом с `docker-compose.yml` по примеру из `env.example`. Секреты нельзя коммитить.

```bash
cd phone-tokens-backend
cp env.example .env
```

После копирования заменить placeholder-значения на реальные production-секреты и настройки окружения.

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
