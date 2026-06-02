# Phone Tokens

Phone Tokens - система личных кабинетов для пользователей, агентов и администраторов.

Что делает продукт:

- пользователи регистрируются, заполняют профиль и выпускают phone tokens;
- агенты работают с сертификатами, SMS logs, billing и токенами пользователей;
- администраторы подтверждают CSR-заявки и просматривают SMS logs.

## Запуск через Docker Compose

1. Подготовить `.env` рядом с `docker-compose.yml`.

Минимальный пример:

```env
POSTGRES_USER=phone_tokens_app
POSTGRES_PASSWORD=change-me
POSTGRES_DB=phone_tokens
DATABASE_URL=postgres://phone_tokens_app:change-me@db:5432/phone_tokens?sslmode=disable

HTTP_PORT=8080
JWT_SECRET=change-me-generate-random-secret
JWT_EXPIRES_IN_SEC=3600
FRONTEND_URL=http://localhost

STRIPE_KEY=sk_test_change_me
WEBHOOK_SECRET=whsec_change_me
SERVER_URL=http://localhost

API_KEY=change-me
EMAIL=service@example.com

CA_CERT_PATH=./cert.pem
CA_KEY_PATH=./key.pem
```

2. Запустить стек:

```bash
docker compose up -d --build
```

3. Проверить контейнеры:

```bash
docker compose ps
```

После запуска:

- frontend доступен через nginx на `http://localhost`;
- backend API проксируется через `http://localhost/api/`;
- Swagger доступен по `/swagger/`, если backend поднялся корректно.

Остановить:

```bash
docker compose down
```
