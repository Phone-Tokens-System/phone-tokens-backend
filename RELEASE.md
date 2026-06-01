# План релиза Phone Tokens

## 1. Назначение документа

План релиза фиксирует договоренность между разработкой, бизнесом и эксплуатацией о том, как передать Phone Tokens в эксплуатацию. Документ описывает цели релиза, участников, артефакты поставки, порядок развертывания, проверки, безопасность и календарный план.

Основной целевой формат - GitHub Wiki-страница публичного проекта. На момент подготовки документа wiki отключена в репозиториях `Phone-Tokens-System/phone-tokens-frontend` и `Phone-Tokens-System/phone-tokens-backend` (`has_wiki=false`), поэтому текущий релизный артефакт хранится как `RELEASE.md` рядом с backend `README`.

## 2. Участники релиза

| Сторона | Роль в релизе | Ответственность |
| --- | --- | --- |
| Разработка | Представитель frontend/backend команды | Подготовить сборки, подтвердить тесты, описать ограничения и известные дефекты |
| Бизнес | Заказчик или куратор продукта | Подтвердить цели релиза, пользовательские сценарии и критерии приемки |
| Эксплуатация | DevOps или команда эксплуатации | Подготовить окружение, секреты, сеть, TLS, Docker Compose и мониторинг |

Перед production handoff участники должны подтвердить готовность релиза в одном месте: GitHub release notes, issue, wiki-комментарий или согласованный release checklist.

## 3. Цели релиза

Релиз передает в эксплуатацию продукт Phone Tokens: личный кабинет пользователя, агента и администратора для работы с телефонными токенами, сертификатами, SMS, биллингом и SSO.

В релиз входят:

- авторизация и регистрация пользователей с ролями `user`, `agent`, `admin`;
- пользовательский профиль и управление user tokens;
- агентский кабинет: certificates, SMS logs, billing, tokens;
- административные разделы: CSR approval и SMS logs;
- SSO-сценарий привязки пользовательского токена к внешнему агенту;
- интеграция backend API через `/api/v1/*`;
- поставка frontend static build, backend Docker image и Docker Compose конфигурации.

В релиз не входит исправление backend-кода в рамках этого документа. Backend проверяется как компонент поставки, а найденные блокеры передаются владельцам backend.

## 4. Release Gate

Релиз нельзя передавать в эксплуатацию, пока не выполнены все условия:

| Gate | Команда или проверка | Текущий статус |
| --- | --- | --- |
| Frontend unit tests | `cd phone-tokens-frontend && npm test` | Пройдено: 9 файлов, 25 тестов |
| Frontend production build | `cd phone-tokens-frontend && npm run build` | Пройдено: Vite build completed |
| Backend tests | `cd phone-tokens-backend && go test ./...` | Блокер: тесты не компилируются |
| Backend Docker image | `cd phone-tokens-backend && docker build -t phone-tokens-backend:<tag> .` | Пройдено локально для `phone-tokens-backend:security-check` |
| Full compose smoke | `cd phone-tokens-backend && docker compose up -d --build` | Требует проверки перед релизом |
| Security acceptance | Проверки из раздела 8 | Требует закрытия критичных пунктов |
| Business acceptance | Tutorial из раздела 9 | Требует ручной приемки |

Текущий backend-блокер:

- `internal/service/tokens/service_test.go`: тестовый `memoryRepo` не реализует метод `GetTokensByUserIdAndAgentId`;
- `internal/service/users/service_test.go`: тесты ссылаются на `ErrNotFound`, который не определен в пакете `users`.

До устранения этих ошибок backend `go test ./...` завершается с `FAIL`, поэтому production handoff запрещен.

## 5. Артефакты поставки

К релизу должны быть подготовлены не только исходники, но и готовые сборки:

| Артефакт | Владелец | Требование |
| --- | --- | --- |
| Frontend static build | Разработка | `npm run build`, результат в `dist/` |
| Frontend Docker build stage | Разработка/DevOps | `phone-tokens-frontend/Dockerfile` копирует `dist` в shared volume |
| Backend Docker image | Разработка/DevOps | Image на базе `phone-tokens-backend/Dockerfile` |
| Docker Compose stack | DevOps | `phone-tokens-backend/docker-compose.yml` поднимает db, backend, frontend-builder, nginx |
| GitHub release/tag | Разработка | Единый release tag с ссылками на frontend и backend commits |
| Production `.env` | DevOps | Хранится вне репозитория, без публикации секретов |
| Release notes | Разработка/бизнес | Список фич, ограничений, блокеров и инструкции отката |

Рекомендуемый формат тегов: `vYYYY.MM.DD`, например `v2026.06.10`.

## 6. Инструкция по развертыванию

### 6.1 Подготовить production environment

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

Запрещено использовать значения из `phone-tokens-backend/env.example` в production: `postgres`, `babababa`, `key`.

### 6.2 Собрать frontend

```bash
cd phone-tokens-frontend
npm ci
npm test
npm run build
```

Ожидаемый результат:

- unit tests завершились успешно;
- каталог `dist/` создан;
- build не содержит ошибок Vite.

### 6.3 Проверить backend

```bash
cd phone-tokens-backend
go test ./...
docker build -t phone-tokens-backend:<release-tag> .
```

Ожидаемый результат:

- `go test ./...` завершился успешно;
- Docker image собран без ошибок;
- backend image не содержит production secrets.

### 6.4 Поднять полный стек

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
- frontend static files отдаются из `/usr/share/nginx/html`;
- `/api/` проксируется в backend.

### 6.5 Проверить nginx routes

```bash
curl -I https://<production-domain>/
curl -I https://<production-domain>/api/v1/dictionary/countries
```

Ожидаемый результат:

- `/` отдает frontend;
- `/api/` проходит через nginx reverse proxy;
- HTTP перенаправляется на HTTPS;
- TLS-сертификат валиден для production domain.

## 7. План отката

Откат выполняется DevOps после согласования с разработкой и бизнесом.

1. Остановить текущий compose stack:

```bash
cd phone-tokens-backend
docker compose down
```

2. Вернуть предыдущий backend image tag и frontend build artifact.
3. Поднять предыдущую версию:

```bash
docker compose up -d
```

4. Проверить `/`, `/api/v1/me`, login/register и основные role-based routes.
5. Зафиксировать инцидент в release notes: причина отката, затронутые пользователи, дальнейшие действия.

Если релиз включает необратимые миграции БД, откат запрещен до отдельного database rollback plan.

## 8. Безопасность

Безопасность является обязательным release gate. Релиз нельзя передавать в эксплуатацию, если есть незакрытые критичные угрозы.

### 8.1 Угрозы

| Угроза | Риск | Обязательное действие |
| --- | --- | --- |
| Утечка JWT secret | Выпуск валидных токенов злоумышленником | Использовать сильный случайный `JWT_SECRET`, хранить вне repo/image, определить rotation policy |
| Слабые пароли БД | Доступ к пользовательским данным | Запретить `postgres/postgres`, использовать secret storage |
| Публичный Postgres | Компрометация БД через интернет | Убрать публикацию `5432:5432` в production, оставить DB только во внутренней Docker network |
| Широкий CORS | Вызовы API с произвольных origin | Заменить `Access-Control-Allow-Origin: *` на production domain |
| Утечка SMS/phone data | Нарушение приватности пользователей | Ограничить доступ role-based endpoints, не логировать токены/телефоны без маскирования |
| Подмена Stripe webhook | Фальшивые пополнения баланса | Проверить `WEBHOOK_SECRET`, принимать webhook только через HTTPS |
| Утечка TLS/private key | Компрометация HTTPS | Не копировать production TLS keys в backend image, хранить сертификаты на nginx/secret volume |
| Доступ к admin endpoints | Несанкционированное управление CSR/SMS | Проверить роль `admin`, протестировать 403 для `user` и `agent` |
| Секреты внутри контейнеров | Утечка при публикации image | Перед сборкой проверить Dockerfile и image history |

### 8.2 Выявленные уязвимости и действия

| Наблюдение | Источник | Действие до релиза |
| --- | --- | --- |
| `env.example` содержал слабые значения `postgres`, `babababa`, `key` | `phone-tokens-backend/env.example` | Исправлено: пример заменен на placeholder-значения и переменные secret paths |
| Backend CORS разрешал любой origin | `internal/adapter/in/http/middleware.go` | Исправлено: CORS разрешает только настроенный `FRONTEND_URL`, чужой preflight получает 403 |
| Compose публиковал Postgres наружу | `phone-tokens-backend/docker-compose.yml` | Исправлено: host port `5432` убран, DB остается внутри Docker network |
| Backend Dockerfile копировал `cert.pem` и `key.pem` | `phone-tokens-backend/Dockerfile`, `.dockerignore`, `docker-compose.yml` | Исправлено: CA files исключены из build context/image и монтируются read-only volume |
| Backend tests не проходят | `go test ./...` | Закрыть до production handoff |

## 9. Tutorial и acceptance-сценарий

Этот сценарий должен пройти представитель бизнеса или QA вместе с DevOps после развертывания staging/prod-like окружения.

### 9.1 User flow

1. Открыть `https://<production-domain>/register`.
2. Зарегистрировать пользователя с ролью `user`.
3. Войти пользователем.
4. Заполнить пользовательский профиль.
5. Создать user token.
6. Проверить список токенов, изменение TTL, freeze/unfreeze и delete для тестового токена.

Ожидаемый результат: пользователь видит только user-разделы, API отвечает без 401/403 для разрешенных действий.

### 9.2 Agent flow

1. Зарегистрировать или подготовить пользователя с ролью `agent`.
2. Войти агентом.
3. Указать `agent_id` в Agent Context, если backend не возвращает его автоматически.
4. Загрузить CSR.
5. Проверить раздел certificates.
6. Проверить SMS logs.
7. Проверить billing balance/top-up flow.

Ожидаемый результат: агент видит agent-разделы, не получает доступ к admin-only маршрутам, billing и SMS endpoints отвечают согласно роли.

### 9.3 Admin flow

1. Войти пользователем с ролью `admin`.
2. Открыть admin CSR page.
3. Найти CSR агента.
4. Approve CSR request.
5. Открыть admin SMS logs.
6. Проверить refresh/sync SMS from provider.

Ожидаемый результат: admin endpoints доступны только роли `admin`, agent и user получают 403.

### 9.4 SSO flow

1. Открыть `/sso/authorize` с параметрами тестового внешнего сервиса.
2. Авторизоваться пользователем.
3. Завершить `/api/v1/sso/complete`.
4. Проверить `/api/v1/sso/me` для полученного токена.

Ожидаемый результат: пользовательский token связан с агентом, внешний сервис получает валидный результат проверки.

## 10. Календарный план

| Дата | Работа | Владелец | Критерий готовности |
| --- | --- | --- | --- |
| 2026-06-04 | Зафиксировать scope релиза, участников, целевые repositories и release document draft | Разработка + бизнес + DevOps | Документ опубликован как `RELEASE.md` в backend-репозитории |
| 2026-06-05 | Подготовить production env checklist, Docker build checklist и список release artifacts | DevOps + разработка | Все обязательные env vars описаны, секреты вынесены из repo |
| 2026-06-06 | Устранить backend test blocker владельцами backend и подтвердить `go test ./...` | Backend-разработка | Backend tests проходят полностью |
| 2026-06-07 | Прогнать frontend `npm test`, backend `go test ./...`, Docker Compose smoke test | Разработка + DevOps | Все release gates зеленые |
| 2026-06-08 | Выполнить security review и закрыть критичные пункты | DevOps + разработка | Нет критичных security blockers |
| 2026-06-09 | Провести ручной acceptance tutorial с бизнесом и DevOps | Бизнес + DevOps + разработка | Acceptance-сценарий пройден и подтвержден |
| 2026-06-10 | Создать GitHub release/tag, опубликовать Docker images, выполнить production handoff | Разработка + DevOps | Система передана в эксплуатацию |

## 11. Чеклист перед handoff

- [ ] Участники релиза назначены поименно в release notes или issue.
- [ ] Frontend tests проходят.
- [ ] Frontend production build проходит.
- [ ] Backend tests проходят.
- [ ] Backend Docker image собран.
- [ ] Full Docker Compose smoke test пройден.
- [ ] Production `.env` создан вне репозитория.
- [ ] Слабые секреты из примеров не используются.
- [ ] Postgres не опубликован наружу в production.
- [ ] CORS ограничен production domain.
- [ ] Stripe webhook secret проверен.
- [ ] HTTPS включен, TLS certificate валиден.
- [ ] Admin endpoints проверены на role-based access.
- [ ] Tutorial/use case пройден бизнесом.
- [ ] Release tag и release notes опубликованы.
- [ ] Rollback path подтвержден DevOps.
