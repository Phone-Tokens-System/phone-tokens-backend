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
