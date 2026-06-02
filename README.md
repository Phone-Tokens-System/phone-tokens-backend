# Phone Tokens

Phone Tokens - система личных кабинетов для пользователей, агентов и администраторов.

Что делает продукт:

- пользователи регистрируются, заполняют профиль и выпускают phone tokens;
- агенты работают с сертификатами, SMS logs, billing и токенами пользователей;
- администраторы подтверждают CSR-заявки и просматривают SMS logs.

## Запуск через Docker Compose

1. Создать `.env` рядом с `docker-compose.yml` и заполнить переменные по примеру из `env.example`.

```bash
cp env.example .env
```

После копирования заменить placeholder-значения на реальные секреты и настройки окружения.

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
