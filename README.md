# Chatter - Group chatting app

![Logo](./docs/attachments/logo.svg)

Документация к проекту расположена в [docs](./docs)

## TODO

Список задач по котором можно отслеживать прогресс проекта

**Прогресс**: (~67%)
Задач всего: 92
Выполнено: 62
Осталось: 30

### Infra (4/9 ~44%)
Задачи связанные с инфраструктурой сервиса, деплоем и локальной разработкой
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

### Backend (35/48 ~73%)
Задачи на серверную часть проекта
1. (13/13) User service
    - [x] Миграция БД
    - [x] Domain структуры
    - [x] SQLC запросы
    - [x] Создание профиля
    - [x] Удаление профиля
    - [x] Логин
    - [x] Разлогин
    - [x] Поиск
    - [x] Обновление профиля
    - [x] Создание сессии
    - [x] Удаление сессии
    - [x] Аутентификационная мидделвара
    - [x] Интеграция с фронтендом
2. (7/7) DM service (личные сообщения)
    - [x] Миграция БД
    - [x] Domain структуры
    - [x] SQLC запросы
    - [x] Создание DM
    - [x] Список DM пользователя
    - [x] Сообщения (первая/следующие страницы)
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
    - [x] Сообщения в топиках
    - [x] Интеграция с Hub (WebSocket)
4. (4/10) WebSocket Hub (realtime сообщения + signaling)
    - [x] WebSocket endpoint `/v1/signaling/ws`
    - [x] Join канала (DM/Server)
    - [x] Leave канала
    - [x] Broadcast сообщений участникам
    - [ ] Call events (EventCallStart, EventCallEnd)
    - [ ] WebRTC signaling (offer, answer, ice candidates)
    - [ ] Call room management (voice топики)
    - [ ] Нотификация дисконекта участников
    - [ ] Broadcast участникам звонка
    - [ ] Интеграция звонка с фронтендом
5. (0/7) Call Service
    - [ ] Domain: Call, CallParticipant структуры
    - [ ] Redis: хранение активных звонков
    - [ ] Hub: создание/удаление комнаты звонка
    - [ ] Hub: вход/выход участников
    - [ ] Signaling: обработка offer/answer
    - [ ] Signaling: обработка ICE candidates
    - [ ] Запись звонков (опционально)


### Frontend (23/40 ~58%)
Задачи на клиентскую часть проекта (SolidJS)
1. (4/4) Общее
    - [x] Vite + SolidJS зависимость
    - [x] Общие стили (Tailwind CSS)
    - [x] Toast уведомления
    - [x] UI компоненты (Button, FormTextInput, Modal, avatar)
2. (3/3) Страницы
    - [x] Страница логина (UI)
    - [x] Страница регистрации (UI)
    - [x] Главная страница App (UI)
3. (10/18) Sidebar
    - [x] Tab switching (direct/servers)
    - [x] Компоненты SidebarServer, SidebarDM
    - [x] API: загрузка списка серверов
    - [x] API: загрузка списка DM
    - [ ] API: поиск серверов
    - [ ] API: поиск пользователей
    - [x] API: создание сервера
    - [ ] API: присоединение к серверу
    - [x] UI: сервер с топиками
    - [x] UI: создание/удаление сервера
    - [x] UI: создание/удаление топика
    - [ ] UI: модалка настроек профиля
    - [ ] API: обновление профиля
    - [ ] API: разлогин
    - [x] JS: menu-bar + content-bar communication
    - [ ] JS: запись табов в URL
    - [x] JS: получение данных профиля
4. (5/7) Чат
    - [x] Компонент Chat
    - [x] Компонент ChatMessage
    - [x] Компонент ChatInput
    - [ ] WebSocket клиент с переподключением
    - [x] API: отправка сообщения
    - [x] API: загрузка истории
    - [ ] JS: авто-прокрутка к новым сообщениям
5. (0/8) Звонки
    - [ ] WebRTC хук (useWebRTC)
    - [ ] Поддержка микрофона
    - [ ] Поддержка камеры
    - [ ] Поддержка трансляции экрана
    - [ ] UI: компонент звонка
    - [ ] UI: отображение участников
    - [ ] Интеграция с signaling WebSocket
    - [ ] Call UI для voice топиков

To fix:
- [x] sort services, topics, dms

### Extra (на подумать)
- [ ] Turn server
- [ ] Unit tests
- [ ] E2E tests
- [ ] Push notifications
