# Chatter - Group chatting app

![Logo](./docs/attachments/logo.svg)

Документация к проекту расположена в [docs](./docs)

## Архитектура

Микросервисная архитектура
- User Service - сервис для хранения пользовательских данных
- Chat Service - сервис роботы с чатами
- Signaling Service - сервис с websocket соединениями для чатов и сигналинга webrtc
- Call Service - сервис с WebRTC для обработки чатов

![architecture diagram](./docs/attachments/architecture.png)
