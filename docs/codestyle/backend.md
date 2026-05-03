# Backend Codestyle

Документ описывает текущие правила backend-части `chattery`. Область: `cmd`, `internal`, `queries`, `migrations`, `Makefile`, `sqlc.yaml`.

## Команды

Основные команды запускаются из корня репозитория:

```sh
make run            # backend dev server через air
make build          # сборка ./bin/chattery из ./cmd/main.go
make up             # docker-compose, backend migrations, e2e migrations
make up-docker      # только docker-compose up -d
make up-migrate     # goose migrations из migrations
make up-e2e         # goose migrations из e2e/migrations
make generate-sqlc  # генерация internal/client/postgres из queries и migrations
make lint           # golangci-lint
make test           # go test ./...
```

После изменения `queries/*.sql` или `migrations/*.sql` обязательно выполнять:

```sh
make generate-sqlc
make test
make lint
```

## Структура проекта

`cmd/main.go` - composition root. Здесь создаются config, logger, Postgres/Redis connections, transaction manager, adapters, stores, services, API servers и запускается HTTP server. Бизнес-логика сюда не добавляется.

`internal/domain` - доменные типы и value objects. Здесь лежат `UserID`, `ServerID`, `DMID`, `Session`, `Password`, `Cursor`, domain-модели. Доменные ID должны быть отдельными типами поверх `int64` и иметь метод `I64()`.

```go
type ServerID int64

func (id ServerID) I64() int64 { return int64(id) }
```

`internal/api/<feature>` - HTTP слой. Пакет содержит:

- `server.go` с локальными интерфейсами зависимостей, `Server`, `New`, `Pattern`, `Route`.
- `model.go` с request/response DTO и converters.
- файлы handlers по одному endpoint/action: `post_create_user.go`, `get_messages.go`, `delete_topic.go`.

`internal/service/<feature>` - бизнес-логика. Пакет содержит:

- `service.go` с локальными интерфейсами зависимостей, `Service`, `New`.
- отдельные файлы use-case/action: `create_dm.go`, `update_server.go`, `first_page_of_messages.go`.
- `validation.go` для приватных проверок, если проверок больше одной-двух.

`internal/adapter/postgres/<feature>` - ручной adapter между domain и sqlc. Пакет содержит:

- `adapter.go` с `Adapter`, `New`, методами доступа к DB.
- `converter.go` с преобразованиями `postgres.*` -> `domain.*`.

`internal/client/postgres` - код, сгенерированный sqlc. Файлы `*.sql.go`, `models.go`, `querier.go`, `db.go` не редактируются вручную.

`internal/client/redis` - низкоуровневый Redis client wrapper.

`internal/adapter/redis` - adapter поверх Redis client, который знает доменные значения и формат ключей.

`internal/store/<feature>` - in-memory read cache. Store синхронизируется через `internal/store/syncer`, защищает состояние `sync.RWMutex`, имеет `Name()` и `Sync(ctx) error`.

`internal/utils/<name>` - маленькие инфраструктурные helper-пакеты: `bind`, `render`, `validate`, `errutil`, `transaction`, `database`, `logger`.

`queries/*.sql` - sqlc-запросы, сгруппированные по feature: `user.sql`, `dm.sql`, `server.sql`.

`migrations/*.sql` - goose migrations с именами `000NN_create_<area>_table.sql`.

## Go naming

Пакеты называются коротко, lower_snake_case только если в имени несколько слов:

```go
package text_topic
package voice_topic
package websocket_manager
```

Для import aliases используется `<feature>_<layer>`, если имя пакета конфликтует или становится неясным:

```go
dm_adapter "chattery/internal/adapter/postgres/dm"
user_store "chattery/internal/store/user"
ws_manager "chattery/internal/service/websocket_manager"
```

Публичные типы слоя называются стандартно:

```go
type Server struct { ... }  // API package
type Service struct { ... } // service package
type Adapter struct { ... } // adapter package
```

Конструктор пакета называется `New` и возвращает pointer:

```go
func New(user userService, dm dmService, cache userCache) *Server
func New(dbAdapter db, transaction txManager, redisAdapter redis, cfg *config.Config) *Service
```

Локальные интерфейсы объявляются на стороне потребителя, в нижнем регистре:

```go
type userService interface {
    AuthRequiredMiddleware(next http.Handler) http.Handler
}
```

Методы handlers называются по HTTP action и сущности:

```go
func (s *Server) PostCreateUser(w http.ResponseWriter, r *http.Request)
func (s *Server) GetDMs(w http.ResponseWriter, r *http.Request)
func (s *Server) DeleteTopic(w http.ResponseWriter, r *http.Request)
```

Приватные validation/converter функции называются по request/response:

```go
func validatePostCreateUserRequest(req *PostCreateUserRequest) error
func convertPostCreateUserRequest(req *PostCreateUserRequest) *domain.User
func convertGetMessagesResponse(cursor *domain.DMCursor, messages []*domain.DMMessage, users map[domain.UserID]*domain.User) *GetMessagesResponse
```

## API layer

Каждый API package реализует интерфейс `internal/api.service`:

```go
func (*Server) Pattern() string {
    return "/v1/user"
}

func (s *Server) Route(router chi.Router) {
    router.Post("/create", s.PostCreateUser)
}
```

Auth-required endpoints группируются через `router.Group` и `AuthRequiredMiddleware`:

```go
router.Group(func(withAuthRouter chi.Router) {
    withAuthRouter.Use(s.user.AuthRequiredMiddleware)
    withAuthRouter.Get("/me", s.GetMe)
})
```

Handler flow:

1. `ctx := r.Context()`.
2. Получить `userID := domain.UserIDFromContext(ctx)`, если endpoint требует авторизацию.
3. Распарсить request через `bind.JSON[T](r)` для JSON body.
4. Провалидировать request локальной `validate...` функцией.
5. Сконвертировать request в domain-тип через `convert...`.
6. Вызвать service.
7. Ответить через `render.JSON(w, r, response)` или `render.Error(w, r, err)`.

Пример:

```go
func (s *Server) PostCreateServer(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    userID := domain.UserIDFromContext(ctx)

    request, err := bind.JSON[PostCreateServerRequest](r)
    if err != nil {
        render.Error(w, r, err)
        return
    }

    if err = validatePostCreateServer(request); err != nil {
        render.Error(w, r, err)
        return
    }

    serverID, err := s.server.CreateServer(ctx, request.Name, userID)
    if err != nil {
        render.Error(w, r, err)
        return
    }

    render.JSON(w, r, convertPostCreateServerResponse(serverID))
}
```

HTTP responses должны иметь JSON tags. Request/response типы именуются `<Method><Action><Entity>Request` и `<Method><Action><Entity>Response`, например `PostCreateServerRequest`.

API слой не должен обращаться к sqlc/postgres/redis напрямую. Он работает через локальные интерфейсы service/cache.

HTTP routes должны использовать verbs семантически:

- `GET` - чтение без изменения состояния;
- `POST` - создание сущности, command/action или read endpoint, которому принципиально нужен JSON body;
- `PUT` - полное/логическое обновление сущности;
- `DELETE` - удаление.

Если read endpoint использует cursor pagination, предпочтительный вариант - `GET` с query params:

```text
GET /v1/dm/messages?dm_id=...&message_id=...&timestamp=...
```

Если для cursor pagination осознанно используется JSON body, route должен оставаться `POST`, а handler/request/response naming должен явно отражать `Post...`, чтобы HTTP method и имя handler не расходились.

## Service layer

Service отвечает за бизнес-правила, транзакции и orchestration. Все DB/cache зависимости передаются через локальные интерфейсы в `service.go`.

Публичный use-case метод открывает транзакцию, если операция меняет состояние:

```go
func (s *Service) UpdateServer(ctx context.Context, serverID domain.ServerID, name string, userID domain.UserID) error {
    return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
        return s.updateServer(ctx, serverID, name, userID)
    })
}
```

Приватный метод содержит шаги use-case и вызывает validation:

```go
func (s *Service) createDM(ctx context.Context, participant1, participant2 domain.UserID) (domain.DMID, error) {
    if err := s.validateCreateDM(ctx, participant1, participant2); err != nil {
        return 0, err
    }
    ...
}
```

Ошибки нижнего слоя оборачиваются `errutil.E(err).Debug("dependency.Method")`. Бизнес-ошибки получают `Kind` и пользовательское `Message/Messagef`.

```go
if err != nil {
    return errutil.E(err).Debug("s.db.CreateServer")
}

return errutil.E().
    Kind(errutil.Permission).
    Message("only owners can perform this action")
```

Service не должен знать про HTTP DTO, `http.ResponseWriter`, chi routes или sqlc models. Исключение текущего проекта: session/cookie methods в `service/user` пока работают с `http.ResponseWriter` и middleware.

Паниковать в service layer и любых backend-пакетах за пределами `cmd/main.go` запрещено. Ошибки initialization/runtime setup должны возвращаться наверх как `error`.

Если приложение невозможно поднять, завершение происходит в `cmd/main.go` через `logger.Fatal` с читаемым контекстом:

```go
postgresConn, err := database.PostgresConnection(appCtx, cfg)
if err != nil {
    logger.Fatal(err, "database.PostgresConnection")
}
```

Конструкторы service/adapter/client, которым может не хватить внешнего ресурса или которые могут получить ошибку initialization, должны возвращать `(*Type, error)` или разделять `New` и `Start/Init` с явным `error`.

## Adapter layer

Postgres adapters принимают `context.Context` первым аргументом и работают через `queryProvider`:

```go
type queryProvider interface {
    Query(ctx context.Context) postgres.Querier
}
```

Каждый adapter method:

1. Конвертирует domain input в `postgres.*Params`.
2. Вызывает `a.db.Query(ctx).<SQLCMethod>(ctx, req)`.
3. Мапит `pgx.ErrNoRows` через `database.NotFound(err)` в `errutil.NotFound`.
4. Мапит constraint violations через `database.IsConstraintViolation`.
5. Конвертирует DB result в domain через `converter.go`.

Пример:

```go
func (a *Adapter) CreateUser(ctx context.Context, user *domain.User) (domain.UserID, error) {
    req := &postgres.CreateUserParams{
        Login:    user.Login.String(),
        Password: user.Password.Hashed(),
        Username: user.Username.String(),
    }
    id, err := a.db.Query(ctx).CreateUser(ctx, req)
    if database.IsConstraintViolation(err, uniqueLoginConstraint) {
        return domain.UserIsUnknown, errutil.E(err).
            Kind(errutil.Exist).
            Messagef("user with login %s already exists", user.Login)
    }
    if err != nil {
        return domain.UserIsUnknown, errutil.E(err).Debug("Query.CreateUser")
    }
    return domain.UserID(id), nil
}
```

## Domain layer

Domain structs не содержат JSON/DB tags. Они выражают внутреннюю модель приложения и используют доменные ID-типы.

```go
type DMMessage struct {
    CreatedAt time.Time
    Text      string
    ID        DMMessageID
    DMID      DMID
    SenderID  UserID
}
```

Строковые value objects должны иметь `String()`:

```go
type Username string

func (u Username) String() string { return string(u) }
```

## Store layer

Store - read cache поверх adapter. Он должен:

- иметь `Name() string` для syncer/logs;
- иметь `Sync(ctx context.Context) error`;
- защищать внутреннее состояние `sync.RWMutex`;
- строить map indexes через `sliceutil.SliceToMap`;
- возвращать `errutil.NotFound`, если сущность не найдена.

Пример:

```go
func (s *ServerStore) GetByID(id domain.ServerID) (*domain.Server, error) {
    s.m.RLock()
    defer s.m.RUnlock()

    server, ok := s.serversByID[id]
    if !ok {
        return nil, errutil.E().Kind(errutil.NotFound).Messagef("server id='%d' not found", id.I64())
    }
    return server, nil
}
```

## Errors and validation

Для ошибок используется только `internal/utils/errutil`:

- `InvalidRequest` -> 400
- `Unauthorized` -> 401
- `Permission` -> 403
- `Exist` -> 409
- `NotFound` -> 404
- default/Internal -> 500

Handlers не выбирают HTTP status вручную. Они вызывают `render.Error(w, r, err)`.

Validation request-полей делается в API слое через `internal/utils/validate`. Validation бизнес-состояния делается в service layer.

```go
func validatePostCreateUserRequest(req *PostCreateUserRequest) error {
    if err := validate.Username(req.Username); err != nil {
        return err
    }
    if err := validate.Password(req.Password); err != nil {
        return err
    }
    return validate.Login(req.Login)
}
```

## SQL queries

SQLC queries лежат в `queries/<feature>.sql`. Каждый запрос начинается с annotation:

```sql
-- name: CreateUser :one
INSERT INTO users(login, password, username)
VALUES ($1, $2, $3)
RETURNING id;
```

Правила именования:

- `Create<Entity>` для insert.
- `Update<Entity>` для update.
- `Delete<Entity>...` для delete.
- `Get<Entity>...` для чтения одной сущности.
- `List<Entities>` для полной выборки в cache/store.
- `<Feature>...` допускается, если метод относится к агрегату, например `UserDMs`.
- Название query должно точно описывать WHERE/filter.

Возвращаемость:

- `:one` для одного значения/строки.
- `:many` для списка.
- `:exec` для command без количества строк.
- `:execrows` если adapter должен проверить количество затронутых строк.

Pagination queries должны сортировать по `created_at DESC, id DESC` и использовать cursor condition:

```sql
WHERE topic_id = $1 AND (created_at < $2 OR (created_at = $2 AND id < $3))
ORDER BY created_at DESC, id DESC
LIMIT $4;
```

SQL queries не должны использовать `SELECT *`. Нужно явно перечислять возвращаемые колонки, чтобы generated code и API adapter не зависели от случайного расширения таблицы:

```sql
-- name: UserByID :one
SELECT id, username, login, password, avatar_id, created_at, updated_at
FROM users
WHERE id = $1;
```

После изменения SQL запускать `make generate-sqlc`; generated files вручную не править.

## Migrations

Миграции используют goose формат:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (...);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
```

Имена файлов:

```text
00001_create_users_table.sql
00002_create_dms_table.sql
00003_create_servers_table.sql
```

DB naming:

- таблицы и колонки в `snake_case`;
- primary key колонка называется `id`;
- foreign-id колонки называются `<entity>_id`;
- timestamp колонки: `created_at`, `updated_at`;
- auto increment IDs: `BIGSERIAL PRIMARY KEY`;
- IDs-ссылки: `BIGINT NOT NULL`;
- строковые поля: `TEXT NOT NULL`;
- timestamps: `TIMESTAMP NOT NULL DEFAULT NOW()`;
- индексы: `<table>_<columns>_idx`, например `dm_messages_dm_id_created_id_desc_idx`.

## Config

Config читается из env через `internal/utils/bind` в `internal/config/config.go`. Новые настройки добавляются в соответствующую вложенную struct (`HTTP`, `Redis`, `Postgres`, `Session`, `Cache`, `Voice`, `Chat`, `App`) и имеют default value.

Env names пишутся UPPER_SNAKE_CASE:

```go
bind.EnvString("POSTGRES_URL", "postgresql://user:password@localhost:5432/chattery?sslmode=disable")
bind.EnvDuration("CACHE_USER_SYNC_TIMEOUT", 30*time.Second)
```

## Logging

Для request errors используется `render.Error`, который логирует через `logger.ErrorCtx`.

`render.Timestamp` считает `today/yesterday` и форматирует входное время в UTC. Если понадобится user-local timezone, policy должна быть изменена явно во всем API.

Для startup fatal errors используется:

```go
logger.Fatal(err, "database.PostgresConnection")
```

Для store sync logs используется `slog.Info` с именем store:

```go
slog.Info("[user_store] update", slog.Int("len", len(users)))
```

## HTTP server lifecycle

HTTP server должен запускаться через `http.Server`, а не напрямую через `http.ListenAndServe`.

Обязательные свойства server lifecycle:

- `ReadHeaderTimeout`;
- `ReadTimeout` и/или осознанное объяснение, почему он не нужен;
- `WriteTimeout` и/или осознанное объяснение для streaming/websocket endpoints;
- `IdleTimeout`;
- graceful shutdown по signal/context;
- понятная ошибка startup/shutdown через `logger.Fatal` или возвращаемый `error`.

`Run` предпочтительно принимает `context.Context` или получает shutdown context от `cmd/main.go`:

```go
func (s *Server) Run(ctx context.Context) error {
    server := &http.Server{
        Addr:              s.address,
        Handler:           s.mux,
        ReadHeaderTimeout: 5 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    errCh := make(chan error, 1)
    go func() {
        errCh <- server.ListenAndServe()
    }()

    select {
    case <-ctx.Done():
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        return server.Shutdown(shutdownCtx)
    case err := <-errCh:
        return err
    }
}
```

## Tests and generated code

Минимальная проверка перед сдачей backend изменений:

```sh
make test
make lint
```

Generated code в `internal/client/postgres` не редактируется вручную. Источник правды - `queries`, `migrations`, `sqlc.yaml`.
