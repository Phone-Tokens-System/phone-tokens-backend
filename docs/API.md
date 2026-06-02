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

