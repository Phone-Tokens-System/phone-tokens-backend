# Security

## Сделанные меры

- CORS ограничен настроенным `FRONTEND_URL`; чужой preflight получает `403`.
- PostgreSQL не публикуется через `5432:5432` в Docker Compose.
- `cert.pem` и `key.pem` не копируются в backend image; они монтируются как read-only volumes.
- `.dockerignore` исключает `.env`, `*.env`, `*.pem`, `*.key`, `*.csr`.
- `env.example` использует placeholders вместо слабых секретов `postgres`, `babababa`, `key`.
- Backend image проверен: `/app/cert.pem` и `/app/key.pem` отсутствуют.

## Проверки перед релизом

- Production `.env` находится вне репозитория.
- `JWT_SECRET`, `POSTGRES_PASSWORD`, `STRIPE_KEY`, `WEBHOOK_SECRET`, `API_KEY` заданы production-grade секретами.
- `FRONTEND_URL` равен production origin.
- HTTPS включен, TLS certificate валиден.
- Admin endpoints доступны только роли `admin`.
- Agent endpoints доступны только роли `agent`.
- User endpoints доступны только роли `user`.

## Security smoke

Разрешенный CORS origin:

```bash
curl -i -X OPTIONS http://localhost:8080/api/v1/me \
  -H 'Origin: http://localhost:5173' \
  -H 'Access-Control-Request-Method: GET'
```

Заблокированный CORS origin:

```bash
curl -i -X OPTIONS http://localhost:8080/api/v1/me \
  -H 'Origin: https://evil.example.com' \
  -H 'Access-Control-Request-Method: GET'
```

Ожидаемый результат для заблокированного origin: `403`.

Проверка секретов в image:

```bash
docker run --rm --entrypoint sh phone-tokens-backend:<tag> \
  -c 'test ! -e /app/cert.pem && test ! -e /app/key.pem'
```
