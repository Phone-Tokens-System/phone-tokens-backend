# Phone Tokens Wiki

Phone Tokens - система личных кабинетов для пользователей, агентов и администраторов. Продукт работает как split-приложение: Go backend, Vue/Vite frontend, PostgreSQL, Docker Compose и nginx.

## Навигация

| Страница | Назначение |
| --- | --- |
| [Installation](Installation.md) | Локальный запуск backend, frontend и тестовой БД |
| [API](API.md) | Swagger, группы endpoint и авторизация |
| [Architecture](Architecture.md) | Компоненты системы, роли и основные потоки |
| [Database](Database.md) | Миграции, основные таблицы и эксплуатационные требования |
| [Deployment](Deployment.md) | Production env, Docker Compose и nginx |
| [Security](Security.md) | Сделанные security-исправления и проверки перед релизом |
| [Release Plan](Release-Plan.md) | Release gates, артефакты, календарный план и handoff checklist |

## Текущий статус проверок

| Проверка | Статус |
| --- | --- |
| Frontend tests | Пройдено: `npm test` |
| Frontend build | Пройдено: `npm run build` |
| Backend targeted tests | Пройдено: `go test ./cmd ./internal/adapter/in ./internal/adapter/in/http` |
| Backend full tests | Блокер: ошибки компиляции тестов `internal/service/tokens`, `internal/service/users` |
