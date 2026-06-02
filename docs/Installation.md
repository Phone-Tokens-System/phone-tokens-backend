# Installation

## Требования

- Go 1.25+.
- Node.js 18+ or 20+.
- Docker.
- Free local ports: `5173`, `8080`, optional PostgreSQL test port `55432`.

## Локальный запуск PostgreSQL

```bash
docker run --name phone-tokens-local-db \
  -e POSTGRES_USER=phone_tokens_app \
  -e POSTGRES_PASSWORD=phone_tokens_pass \
  -e POSTGRES_DB=phone_tokens \
  -p 55432:5432 \
  -d postgres:16-alpine
```

Проверить готовность:

```bash
docker exec phone-tokens-local-db pg_isready -U phone_tokens_app -d phone_tokens
```

## Запуск backend

```bash
cd phone-tokens-backend
DATABASE_URL='postgres://phone_tokens_app:phone_tokens_pass@localhost:55432/phone_tokens?sslmode=disable' \
HTTP_PORT=8080 \
JWT_SECRET='local-dev-secret-at-least-32-bytes' \
JWT_EXPIRES_IN_SEC=3600 \
FRONTEND_URL='http://localhost:5173' \
API_KEY='local-api-key' \
EMAIL='local@example.com' \
STRIPE_KEY='sk_test_local' \
WEBHOOK_SECRET='whsec_local' \
SERVER_URL='http://localhost:8080' \
go run ./cmd
```

Backend применяет SQL-миграции из `database/` и запускает API на `http://localhost:8080`.

## Запуск frontend

```bash
cd phone-tokens-frontend
npm ci
npm run dev
```

Frontend запускается на `http://localhost:5173`. Vite проксирует `/api/*` на `http://localhost:8080`.

## Smoke-проверка

1. Открыть `http://localhost:5173/register`.
2. Зарегистрировать пользователя с ролью `user`.
3. Войти этим пользователем.
4. Открыть `/dashboard/profile`.
5. Открыть `/dashboard/tokens`.
6. Создать token и убедиться, что таблица токенов обновилась.
