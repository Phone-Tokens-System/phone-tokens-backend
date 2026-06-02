# План релиза Phone Tokens

## 1. Цели релиза

Релиз передает в эксплуатацию Phone Tokens: личный кабинет пользователя, агента и администратора для работы с телефонными токенами, сертификатами, SMS, биллингом и SSO.

В релиз входят:

- авторизация и регистрация с ролями `user`, `agent`, `admin`;
- пользовательский профиль и управление user tokens;
- агентские разделы: certificates, SMS logs, billing, tokens;
- административные разделы: CSR approval и SMS logs;
- SSO-сценарий привязки пользовательского токена к внешнему агенту;
- frontend static build, backend Docker image и Docker Compose конфигурация.

Текущий статус релиза: `No-Go`. Передача в эксплуатацию запрещена до закрытия backend test blocker.

## 2. Release Gate

| Gate | Команда или проверка | Текущий статус |
| --- | --- | --- |
| Frontend unit tests | `cd phone-tokens-frontend && npm test` | Пройдено: 9 файлов, 25 тестов |
| Frontend production build | `cd phone-tokens-frontend && npm run build` | Пройдено |
| Backend tests | `cd phone-tokens-backend && go test ./...` | Блокер: тесты не компилируются |
| Backend Docker image | `cd phone-tokens-backend && docker build -t phone-tokens-backend:<tag> .` | Пройдено локально для `phone-tokens-backend:security-check` |
| Full compose smoke | `cd phone-tokens-backend && docker compose up -d --build` | Требует проверки перед релизом |
| Security acceptance | Проверки из раздела 5 | Требует подтверждения перед релизом |

Текущий backend-блокер:

- `internal/service/tokens/service_test.go`: тестовый `memoryRepo` не реализует `GetTokensByUserIdAndAgentId`;
- `internal/service/users/service_test.go`: тесты ссылаются на неопределенный `ErrNotFound`.

## 3. Артефакты поставки

| Артефакт | Владелец | Требование |
| --- | --- | --- |
| Frontend static build | Frontend | `npm run build`, результат в `dist/` |
| Backend Docker image | Backend/DevOps | Image на базе `phone-tokens-backend/Dockerfile` |
| Docker Compose stack | DevOps | `phone-tokens-backend/docker-compose.yml` поднимает db, backend, frontend-builder, nginx |
| Production `.env` | DevOps | Хранится вне репозитория, без публикации секретов |
| GitHub release/tag | Бизнес/DevOps | Единый release tag со ссылками на frontend и backend commits |

Рекомендуемый формат тегов: `vYYYY.MM.DD`, например `v2026.06.10`.

## 4. Инструкция по развертыванию

### 4.1 Подготовить production environment

Создать production `.env` для backend/compose. Файл не коммитить в репозиторий.

Обязательные переменные:

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

### 4.2 Собрать и проверить frontend

```bash
cd phone-tokens-frontend
npm ci
npm test
npm run build
```

### 4.3 Собрать и проверить backend

```bash
cd phone-tokens-backend
go test ./...
docker build -t phone-tokens-backend:<release-tag> .
```

### 4.4 Поднять полный стек

```bash
cd phone-tokens-backend
docker compose up -d --build
docker compose ps
```

Ожидаемый результат:

- `db` здоров по healthcheck;
- `backend` доступен внутри Docker network на `backend:8080`;
- Postgres не опубликован наружу через host port;
- CA certificate/private key примонтированы в backend container как read-only volume;
- `nginx` публикует HTTP/HTTPS наружу;
- `/` отдает frontend, `/api/` проксируется в backend.

## 5. Безопасность

В рамках подготовки релиза уже сделано:

- CORS ограничен: backend разрешает только настроенный `FRONTEND_URL`, чужой preflight получает `403`;
- Postgres закрыт от внешней сети: из `docker-compose.yml` убрана публикация `5432:5432`;
- `cert.pem` и `key.pem` больше не копируются в backend image, а монтируются как read-only volume;
- `.dockerignore` исключает `.env`, `*.env`, `*.pem`, `*.key`, `*.csr`, чтобы секреты не попадали в Docker build context;
- `env.example` заменен на placeholder-значения без слабых секретов `postgres`, `babababa`, `key`;
- backend image проверен: `/app/cert.pem` и `/app/key.pem` внутри образа отсутствуют.

Перед релизом нужно подтвердить:

- production `.env` хранится вне репозитория;
- `JWT_SECRET`, `POSTGRES_PASSWORD`, `STRIPE_KEY`, `WEBHOOK_SECRET`, `API_KEY` заданы production-секретами;
- HTTPS включен, TLS certificate валиден;
- admin endpoints проверены на role-based access.

## 6. Календарный план

| Дата | Работа | Владелец | Критерий готовности |
| --- | --- | --- | --- |
| 2026-06-04 | Зафиксировать scope релиза, участников и release document draft | Бизнес + DevOps | Документ опубликован как `RELEASE.md` в backend-репозитории |
| 2026-06-05 | Подготовить production env checklist, Docker build checklist и список release artifacts | DevOps | Все обязательные env vars описаны, секреты вынесены из repo |
| 2026-06-06 | Устранить backend test blocker и подтвердить `go test ./...` | Backend | Backend tests проходят полностью |
| 2026-06-07 | Прогнать frontend `npm test`, backend `go test ./...`, Docker Compose smoke test | QA + DevOps | Все release gates зеленые |
| 2026-06-08 | Выполнить security acceptance | DevOps | Нет критичных security blockers |
| 2026-06-09 | Провести ручную приемку с бизнесом и DevOps | Бизнес + QA + DevOps | Приемка подтверждена |
| 2026-06-10 | Создать GitHub release/tag, опубликовать Docker images, выполнить production handoff | Бизнес + DevOps | Система передана в эксплуатацию |

## 7. Чеклист перед handoff

- [ ] Участники релиза назначены поименно в release notes или issue.
- [ ] Frontend tests проходят.
- [ ] Frontend production build проходит.
- [ ] Backend tests проходят.
- [ ] Backend Docker image собран.
- [ ] Full Docker Compose smoke test пройден.
- [ ] Production `.env` создан вне репозитория.
- [ ] Production-секреты не совпадают с примерами из `env.example`.
- [ ] Postgres не опубликован наружу.
- [ ] CORS ограничен production domain.
- [ ] Stripe webhook secret проверен.
- [ ] HTTPS включен, TLS certificate валиден.
- [ ] Admin endpoints проверены на role-based access.
- [ ] Release tag и release notes опубликованы.
