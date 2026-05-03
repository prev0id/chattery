# Chatter - Group chatting app

![Logo](./docs/attachments/logo.svg)

Документация к проекту расположена в [docs](./docs)

## TODO

Список задач, по которому можно отслеживать прогресс проекта.

**Прогресс**: 84/108 (~78%)
Задач всего: 108
Выполнено: 84
Осталось: 24

### Infra (4/9 ~44%)
Задачи, связанные с инфраструктурой сервиса, деплоем и локальной разработкой.
- (4/4) Локальная разработка
    - [x] docker-compose
    - [x] Makefile
    - [x] Air (live reload)
    - [x] Миграции
- (0/5) Деплой
    - [ ] Домен
    - [ ] DNS
    - [ ] Postgres
    - [ ] Redis
    - [ ] Server

### Backend (39/51 ~76%)
Задачи на серверную часть проекта.
1. (13/13) User service
    - [x] Миграция БД
    - [x] Domain структуры
    - [x] SQLC запросы
    - [x] Создание профиля
    - [x] Удаление профиля
    - [x] Логин
    - [x] Разлогин
    - [x] Поиск пользователей
    - [x] Обновление профиля
    - [x] Создание сессии
    - [x] Удаление сессии
    - [x] Аутентификационная middleware
    - [x] Интеграция с фронтендом для login/signup/me
2. (8/8) DM service (личные сообщения)
    - [x] Миграция БД
    - [x] Domain структуры
    - [x] SQLC запросы
    - [x] Создание DM
    - [x] Список DM пользователя
    - [x] Сообщения (первая/следующие страницы)
    - [x] Ручка чтения DM: `POST /v1/dm/read`
    - [x] Интеграция с Hub (WebSocket)
3. (11/11) Server service (публичные сервера с топиками)
    - [x] Миграция БД (servers, topics)
    - [x] Domain структуры
    - [x] SQLC запросы
    - [x] Создание сервера
    - [x] Обновление сервера
    - [x] Удаление сервера
    - [x] Присоединение к серверу
    - [x] Выход из сервера
    - [x] CRUD топиков
    - [x] Сообщения в текстовых топиках
    - [x] Интеграция с Hub (WebSocket)
4. (7/13) WebSocket Hub (realtime сообщения + будущий signaling)
    - [x] WebSocket endpoint `/ws/`
    - [x] Ping/pong heartbeat
    - [x] Redis pub/sub доставка событий пользователю
    - [x] Join канала DM/text topic с проверкой доступа
    - [x] Leave канала
    - [x] Broadcast сообщений участникам DM/text topic
    - [x] Глобальная доставка DM-событий пользователю для live sidebar/notifications
    - [ ] Join канала voice topic с проверкой доступа
    - [ ] Call events (start/end/join/leave)
    - [ ] WebRTC signaling (offer, answer, ice candidates)
    - [ ] Call room management для voice топиков
    - [ ] Нотификация дисконекта участников звонка
    - [ ] Broadcast участникам звонка
5. (0/6) Voice Topic Service
    - [ ] Domain: Call, CallParticipant структуры
    - [ ] Redis: хранение активных звонков
    - [ ] Создание/удаление комнаты звонка
    - [ ] Вход/выход участников
    - [ ] Signaling: обработка offer/answer
    - [ ] Signaling: обработка ICE candidates

### Frontend (41/48 ~85%)
Задачи на клиентскую часть проекта (SolidJS).
1. (4/4) Общее
    - [x] Vite + SolidJS зависимость
    - [x] Общие стили (Tailwind CSS)
    - [x] Toast уведомления
    - [x] UI компоненты (Button, FormTextInput, Modal, avatar)
2. (3/3) Страницы
    - [x] Страница логина (UI + API)
    - [x] Страница регистрации (UI + API)
    - [x] Главная страница App (UI + роутинг)
3. (17/20) Sidebar
    - [x] Tab switching (direct/servers)
    - [x] Компоненты SidebarServer, SidebarDM
    - [x] API: загрузка списка серверов
    - [x] API: загрузка списка DM
    - [x] API/UI: поиск серверов
    - [x] API/UI: поиск пользователей и создание DM из поиска
    - [x] API: создание сервера
    - [x] API/UI: присоединение к серверу
    - [x] API/UI: выход из сервера
    - [x] UI: сервер с топиками
    - [x] UI/API: создание, обновление и удаление сервера
    - [x] UI/API: создание, обновление и удаление топика
    - [x] UI: каркас модалки настроек профиля
    - [ ] API/UI: обновление профиля
    - [ ] API/UI: загрузка и удаление аватара
    - [ ] API/UI: разлогин
    - [x] JS: menu-bar + content-bar communication
    - [x] JS: состояние вкладок через URL (`/dm`, `/server`)
    - [x] JS: получение данных профиля
    - [x] JS: обновление списков после мутаций без ручного reload
4. (11/11) Чат
    - [x] Компонент Chat
    - [x] Компонент ChatMessage
    - [x] Компонент ChatInput
    - [x] WebSocket клиент с переподключением
    - [x] Единое WebSocket подключение без сброса при переходах между чатами
    - [x] API: отправка сообщения
    - [x] API: загрузка истории
    - [x] API: пометка DM прочитанным при live-сообщении в открытом чате
    - [x] JS: авто-прокрутка к новым сообщениям
    - [x] JS: догрузка старых сообщений по cursor при скролле вверх
    - [x] JS: live preview/unread для DM и toast-нотификации входящих сообщений
5. (6/10) Звонки
    - [x] Media hook для локальных устройств
    - [x] Поддержка микрофона
    - [x] Поддержка камеры
    - [x] Поддержка трансляции экрана
    - [x] UI: кнопки звонка и настройки устройств
    - [x] UI: локальный preview stream
    - [ ] WebRTC peer/signaling hook
    - [ ] UI: отображение удаленных участников
    - [ ] Интеграция с signaling WebSocket
    - [ ] Call UI подключает пользователя к backend-комнате voice topic

### To fix
- [x] sort services, topics, dms
- [ ] Подключить profile settings modal к реальным `/v1/user/*` endpoint'ам вместо неподключенных HTML form handlers
- [ ] Дописать закомментированный `internal/service/voice_topic`
- [ ] Добавить e2e/unit тесты для сервисов и API handlers
- [ ] перенести toast в entry

### Extra (на подумать)
- [ ] Turn server
- [ ] Unit tests
- [ ] E2E tests
- [x] Notifications
- [ ] Desc package
- [x] Password as struct

## работа очередей

### Messaging

-> function call
<- PUBSUB by user_id

user -> ws :  {join, chat_id}
ws -> rabbitmq : PUBLISH chat_events {chat_id, join, user_id}
rabbitmq -> chat : CONSUME chat_events {chat_id, join, user_id}
chat -> reddis : SET chat:chat_id:active

user -> chat : POST /v1/api/chat/message {chat_id, message}
chat -> reddis : GET chat:chat_id:active
chat -> rabbitmq : PUBLISH user_events_exchange {chat_id, message}
rabbitmq -> ws : CONSUME user_events_queue {chat_id, join, user_id}

### Call service

user -> ws : {join, call_id}
ws -> call : JoinUser(call_id, user_id, conn_id)
call -> reddis : SET call:{call_id}:lock {service_id} NX PX 10000
call -> rabbitmq : PUBLISH user_events_exchange {conn_id, call, message}
rabbitmq -> ws : CONSUME user_events_queue {chat_id, join, user_id}
user -> ws :  {join_ok, chat_id}
