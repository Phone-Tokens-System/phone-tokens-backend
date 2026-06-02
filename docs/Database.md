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
