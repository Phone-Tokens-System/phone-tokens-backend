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
| Security acceptance | Проверки из [Security](Security.md) | Требует подтверждения перед релизом |

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
