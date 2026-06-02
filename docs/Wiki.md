# Phone Tokens Wiki

Phone Tokens - система личных кабинетов для пользователей, агентов и администраторов. Продукт работает как split-приложение: Go backend, Vue/Vite frontend, PostgreSQL, Docker Compose и nginx.

## Навигация

| Раздел | Назначение |
| --- | --- |
| Installation | Локальный запуск backend, frontend и тестовой БД |
| API | Swagger, группы endpoint и авторизация |
| Architecture | Компоненты системы, роли и основные потоки |
| Database | Миграции, основные таблицы и эксплуатационные требования |
| Deployment | Production env, Docker Compose и nginx |
| Security | Сделанные security-исправления и проверки перед релизом |
| Release Plan | Release gates, артефакты, календарный план и handoff checklist |

## Текущий статус проверок

| Проверка | Статус |
| --- | --- |
| Frontend tests | Пройдено: `npm test` |
| Frontend build | Пройдено: `npm run build` |
| Backend targeted tests | Пройдено: `go test ./cmd ./internal/adapter/in ./internal/adapter/in/http` |
| Backend full tests | Блокер: ошибки компиляции тестов `internal/service/tokens`, `internal/service/users` |

---

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

---

# API

## Swagger

Swagger files:

- `docs/swagger.yaml`
- `docs/swagger.json`
- `docs/docs.go`

When backend is running:

```text
http://localhost:8080/swagger/
```

## Authentication

Protected endpoints use bearer JWT:

```http
Authorization: Bearer <jwt-token>
```

JWT claims are parsed by auth middleware. Role checks use `RequireRole`.

## Endpoint Groups

| Group | Endpoints | Access |
| --- | --- | --- |
| Auth | `POST /api/v1/register`, `POST /api/v1/login`, `GET /api/v1/me` | Public/authenticated |
| Tokens | `POST /api/v1/tokens`, `GET /api/v1/users/{userId}/tokens`, `PATCH /api/v1/tokens/{tokenID}`, `DELETE /api/v1/tokens/{tokenID}` | `user` |
| User Profile | `GET /api/v1/user-profile/filters`, `POST /api/v1/user-profile`, `GET /api/v1/user-profile/me` | Public/authenticated `user` |
| Agent CSR | `POST /api/v1/csr/upload`, `GET /api/v1/csr/signed` | `agent` |
| Agent SMS | `POST /api/v1/sms/send`, `POST /api/v1/sms/send_filtered`, `GET /api/v1/sms/agents/{agentId}` | `agent` |
| Admin CSR | `GET /api/v1/admin/csr`, `POST /api/v1/admin/csr/approve/{id}` | `admin` |
| Admin SMS | `GET /api/v1/sms/logs`, `GET /api/v1/sms/all`, `GET /api/v1/sms/status` | `admin` |
| Billing | `POST /api/v1/billing/balance`, `GET /api/v1/billing/{agent_id}/balance`, `GET /api/v1/agents/{agent_id}/transactions` | `agent` |
| SSO | `GET /api/v1/sso/authorize`, `POST /api/v1/sso/complete`, `GET /api/v1/sso/me` | Mixed |
| Dictionaries | `GET /api/v1/dictionary/countries`, `GET /api/v1/dictionary/regions`, `GET /api/v1/dictionary/cities` | Public |

---

# Architecture

## Компоненты системы

| Компонент | Путь | Ответственность |
| --- | --- | --- |
| Backend API | `phone-tokens-backend` | Go HTTP API, services, repositories, migrations |
| Frontend UI | `phone-tokens-frontend` | Vue/Vite кабинеты пользователя, агента и администратора |
| Database | PostgreSQL | Users, tokens, profiles, agents, SMS, billing, CSR data |
| Reverse Proxy | `nginx/conf.d/app.conf` | HTTPS, frontend static files, proxy `/api/` |
| External Providers | Stripe, SMS Aero | Billing checkout/webhooks и отправка SMS |

## Слои backend

| Слой | Путь | Назначение |
| --- | --- | --- |
| Entrypoint | `cmd/main.go` | Загружает config, собирает services, запускает HTTP server |
| App bootstrap | `internal/app` | Config, service wiring, migrations |
| HTTP adapters | `internal/adapter/in/http` | Routes, handlers, auth, CORS, role guards |
| Repositories | `internal/adapter/out/repository` | PostgreSQL persistence |
| Domain services | `internal/service` | Users, tokens, certificates, SMS, billing |
| Models/DTO | `internal/model`, `internal/adapter/dto` | Persistence models and API payloads |

## Потоки по ролям

### User

```text
register/login -> JWT -> profile -> tokens -> token operations
```

### Agent

```text
login -> agent dashboard -> CSR upload -> certificate -> SMS/billing flows
```

### Admin

```text
login -> admin CSR/SMS pages -> approve CSR -> inspect/sync SMS logs
```

## Request flow

```text
Browser -> nginx -> backend /api -> HTTP handler -> service -> repository -> PostgreSQL
```

В локальной разработке Vite проксирует `/api/*` напрямую на `http://localhost:8080`.

---

# Database

## Движок

Проект использует PostgreSQL. Локальный и production Docker Compose используют `postgres:16-alpine`.

## Миграции

Файлы миграций лежат в `database/`.

Backend применяет миграции при старте:

```text
internal/app/migrate.go -> goose.Up(db, "database")
```

Текущая последовательность миграций:

| Prefix | Назначение |
| --- | --- |
| `001_init.sql` | Initial users table |
| `002_tokens.sql` | User tokens |
| `003_token_metadata.sql` | Token metadata |
| `004_certificate_requests.sql` | CSR requests |
| `005_agent_info.sql` | Agent certificate info |
| `006_add_service_name.sql` | Service name fields |
| `007_sms_table.sql` | SMS table |
| `008_add_agent_id_token.sql` | Agent ID on tokens |
| `009_make_csr_id_unique.sql` | CSR uniqueness |
| `010_agents.sql` | Agents |
| `011_make_ext_id_string_sms.sql` | SMS external ID type |
| `012_add_token_delete_number_sms.sql` | SMS token field |
| `013_add_transactions.sql` | Billing transactions |
| `014_add_usage.sql` | Usage accounting |
| `015_add_packages.sql` | Packages |
| `016_add_agent_packages.sql` | Agent packages |
| `017_add_user_profile.sql` | User profile |
| `018_add_agent_id_certificate_requests.sql` | Agent ID on CSR |
| `019_add_package_duration.sql` | Package duration |

## Основные области данных

| Область | Tables/Models |
| --- | --- |
| Users | `users`, `internal/model/user.go` |
| Tokens | `user_tokens`, `internal/model/token.go` |
| Profiles | user profile migrations and `internal/model/user_profile.go` |
| Agents | `agents`, `agent_packages`, `internal/model/agent.go` |
| Certificates | `certificate_requests`, certificate info models |
| SMS | `sms`, `internal/model/sms.go` |
| Billing | `transactions`, `usage`, packages |

## Production notes

- PostgreSQL не должен быть доступен из публичной сети.
- DB credentials должны приходить из production secrets.
- Миграции выполняются при старте backend; rollback plan должен учитывать изменения схемы.
- Перед production deployment нужен backup.

---

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

---

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

---

# Release Plan

## 1. Цели

Релиз передает в эксплуатацию Phone Tokens: личный кабинет пользователя, агента и администратора для работы с телефонными токенами, сертификатами, SMS, биллингом и SSO.

В релиз входят:

- авторизация и регистрация с ролями `user`, `agent`, `admin`;
- пользовательский профиль и управление user tokens;
- агентские разделы: certificates, SMS logs, billing, tokens;
- административные разделы: CSR approval и SMS logs;
- SSO-сценарий привязки пользовательского токена к внешнему агенту;
- frontend static build, backend Docker image и Docker Compose конфигурация.

## 2. Release Gate

| Gate | Команда или проверка | Статус |
| --- | --- | --- |
| Frontend unit tests | `cd phone-tokens-frontend && npm test` | Пройдено: 9 файлов, 25 тестов |
| Frontend production build | `cd phone-tokens-frontend && npm run build` | Пройдено |
| Backend tests | `cd phone-tokens-backend && go test ./...` | Блокер: тесты не компилируются |
| Backend Docker image | `cd phone-tokens-backend && docker build -t phone-tokens-backend:<tag> .` | Пройдено локально для `phone-tokens-backend:security-check` |
| Full compose smoke | `cd phone-tokens-backend && docker compose up -d --build` | Требует проверки перед релизом |
| Security acceptance | Проверки из Security | Требует подтверждения перед релизом |

Backend-блокеры:

- `internal/service/tokens/service_test.go`: тестовый `memoryRepo` не реализует `GetTokensByUserIdAndAgentId`;
- `internal/service/users/service_test.go`: тесты ссылаются на неопределенный `ErrNotFound`.

## 3. Артефакты поставки

| Артефакт | Владелец | Требование |
| --- | --- | --- |
| Frontend static build | Frontend | `npm run build`, результат в `dist/` |
| Backend Docker image | Backend/DevOps | Image на базе `phone-tokens-backend/Dockerfile` |
| Docker Compose stack | DevOps | db, backend, frontend-builder, nginx |
| Production `.env` | DevOps | Хранится вне репозитория |
| GitHub release/tag | Бизнес/DevOps | Единый release tag со ссылками на frontend/backend commits |

Рекомендуемый формат тегов: `vYYYY.MM.DD`, например `v2026.06.10`.

## 4. Календарный план

| Дата | Работа | Владелец | Критерий готовности |
| --- | --- | --- | --- |
| 2026-06-04 | Зафиксировать scope релиза, участников и release document draft | Бизнес + DevOps | Документ опубликован в backend repository |
| 2026-06-05 | Подготовить production env checklist, Docker build checklist и список artifacts | DevOps | Env vars описаны, секреты вне repo |
| 2026-06-06 | Устранить backend test blocker и подтвердить `go test ./...` | Backend | Backend tests проходят |
| 2026-06-07 | Прогнать frontend tests, backend tests и Docker Compose smoke | QA + DevOps | Release gates зеленые |
| 2026-06-08 | Выполнить security acceptance | DevOps | Нет критичных security blockers |
| 2026-06-09 | Провести ручную бизнес-приемку | Бизнес + QA + DevOps | Приемка подтверждена |
| 2026-06-10 | Создать GitHub release/tag, опубликовать Docker images, выполнить production handoff | Бизнес + DevOps | Система передана в эксплуатацию |

## 5. Handoff checklist

- [ ] Участники релиза назначены поименно в release notes или release issue.
- [ ] Frontend tests проходят.
- [ ] Frontend production build проходит.
- [ ] Backend tests проходят.
- [ ] Backend Docker image собран.
- [ ] Full Docker Compose smoke test пройден.
- [ ] Production `.env` создан вне repository.
- [ ] Production secrets отличаются от `env.example`.
- [ ] PostgreSQL не опубликован наружу.
- [ ] CORS ограничен production domain.
- [ ] Stripe webhook secret проверен.
- [ ] HTTPS включен, TLS certificate валиден.
- [ ] Admin endpoints проверены на role-based access.
- [ ] Release tag и release notes опубликованы.
