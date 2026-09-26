# go-microservices

Two Go services and a CLI. `auth` stores users, issues JWT, and checks access. `chat-server` stores chats and messages and pushes new messages to connected clients. Shared Postgres code lives in [platform_common](https://github.com/KornilovIvan/platform_common) `v1.0.0`.

Each service is its own module. `chat-server` does not import `auth`. It keeps a copy of the `Access` proto and a generated client.

## Layout

| Path | Module | Role |
| --- | --- | --- |
| `auth/` | `github.com/ivankornilov/auth` | Users, login, tokens, access checks, HTTP gateway, Swagger, logs, metrics, traces |
| `chat-server/` | `github.com/ivankornilov/chat-server` | Chats, messages, server stream, access check before unary RPCs |
| `cli/` | `github.com/ivankornilov/cli` | Cobra commands `create user` and `delete user` |

Inside each service the call path is `api` → `service` → `repository` → Postgres. `serviceProvider` wires the dependencies. Configuration is loaded from an env file with `godotenv` (`--config-path`, default `local.env`).

## Auth

gRPC `50051`, HTTP gateway `8080`, Swagger `8082`, metrics `2112`.

**Users.** `Create`, `Get`, `Update`, `Delete`. Role is `USER` or `ADMIN`. `Update` takes optional `name` and `email` as `google.protobuf.StringValue`. `Create`, `Update`, and `Delete` write a row to `user_logs` in the same transaction as the user change. `user_logs` has no foreign key to `users`, so a delete can be logged after the user row is gone.

**Tokens.** `Login` loads the user by `name`, compares the password with the stored value, and returns a refresh token. `GetRefreshToken` checks that refresh token and returns a new one. `GetAccessToken` exchanges a refresh token for an access token. Both are HS256 JWTs (`dgrijalva/jwt-go`) signed with different secrets. Refresh TTL is 60 minutes, access TTL is 5 minutes. Claims hold the user name and role.

**Access.** `AccessV1.Check` reads `authorization: Bearer <access>` from gRPC metadata and verifies the access token. The role map allows `ADMIN` on `/note_v1.NoteV1/Get`. Any other method is allowed when the token is valid.

**HTTP.** `protoc-gen-validate` checks proto rules in a unary interceptor before the handler: non-empty name, email format, non-empty password, `id > 0`. Password confirmation is checked in the user service. grpc-gateway exposes the same RPCs over HTTP. `protoc-gen-openapiv2` writes `pkg/swagger/api.swagger.json`. Swagger UI is embedded with `statik` and served on its own port. The spec host is `localhost:8080`.

## Chat-server

gRPC `50052`. Its own database is `chat`.

`Create` inserts a chat and its members, then opens an in-memory channel for that id. `SendMessage` inserts the message and publishes it to the channel. `ConnectChat` is a server stream: the client sends `chat_id` and `username` and receives messages until it disconnects. `Delete` removes the chat and closes the channel.

A unary interceptor copies incoming metadata and calls `auth` `AccessV1.Check` with the full method name. `Create`, `Delete`, and `SendMessage` run only after that call succeeds. The auth address is `AUTH_GRPC_HOST` and `AUTH_GRPC_PORT` (`localhost:50051` locally, `host.docker.internal:50051` in the chat container).

## Data

Migrations use goose.

`auth` tables: `users`, `user_logs`.  
`chat-server` tables: `chats`, `chat_users`, `messages`. A username is unique within a chat. Messages and members are removed with the chat.

Queries are built with squirrel. Transactions go through `TxManager` at `Read Committed`, with the transaction stored in context. The Postgres client, query wrapper, and transaction manager come from `platform_common`.

## Observability

`auth` unary interceptor chain: tracing, logging, metrics, validation.

- **Logs.** `zap` records method, request, response, and duration for every gRPC call. Output goes to the console and to `logs/app.log`. lumberjack rotates the file at 10 MB, keeps 3 backups, and drops files older than 7 days. Level comes from `-l` and defaults to `info`.
- **Metrics.** Request counter, response counter (`status`, `method`), and a response-time histogram. Served at `GET /metrics`. Compose runs Prometheus on `9090` (scrape target `app:2112`) and Grafana on `3000`. `alerts.yml` raises `TargetIsDown` when a target is down for 30 seconds.
- **Traces.** Jaeger `all-in-one` `1.48`, UI on `16686`. Sampler is `const/1`. Each span is named with the full method name, and the response carries `x-trace-id`. The gateway dial to the local gRPC server continues the same trace. The agent host is `localhost` locally and `jaeger` in Compose.

## CLI

`my-app create user -u <name>` prints that the user was created. `my-app delete user -u <name>` prints that the user was deleted. `--username` is required.

```bash
cd cli
go run ./cmd/main.go create user -u ivan
```

`make build` writes `cli/bin/my_app`.

## Tests and CI

Handler tests in `auth` and `chat-server` use `testify` suites and minimock. The service is mocked.

On push to `auth/**` or `chat-server/**`, GitHub Actions runs `go test ./...` in both modules. On `main`, after tests, both images are pushed to GHCR (`ghcr.io/kornilovivan/auth`, `ghcr.io/kornilovivan/chat-server`). The deploy job SSHs to the server, pulls `auth`, and runs it as `auth-app` on port `50051`.

## Local run

From a service directory:

```bash
make local-up
```

`auth` starts Postgres on host port `54321`, runs migrations, and starts the app, Prometheus, Grafana, and Jaeger.  
`chat-server` starts Postgres on host port `54322` and the app. Access checks call `auth` on `50051`.
