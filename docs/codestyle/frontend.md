# Web Codestyle

Документ описывает целевой стиль фронтенда `web`. Текущий код может нарушать эти правила, но новый код должен писаться по ним, а старый код нужно постепенно приводить к ним при любом изменении рядом.

Главная цель: единообразный SolidJS-код, предсказуемая структура проекта, тонкие UI-компоненты, централизованная работа с API и ошибками, понятные имена и минимальная скрытая магия.

## Стек и базовые договоренности

- Приложение пишется на SolidJS, Solid Router, Vite и Tailwind CSS.
- JSX-файлы используются только для компонентов, роутов и Solid-specific модулей.
- Обычная бизнес-логика, API-клиенты, константы, мапперы и утилиты пишутся в `.js`.
- Импорт из `src` делается через алиас `~/`, а не через длинные относительные пути.
- Относительные импорты допустимы только внутри одной папки, например `./MessageItem`.
- `web/dist` не редактируется вручную. Это артефакт сборки.
- Код должен быть форматируемым стандартным prettier-подобным стилем: 2 пробела, двойные кавычки, trailing comma в многострочных объектах/массивах/аргументах, точка с запятой.
- Не добавлять новые зависимости без причины. Для иконок использовать `lucide-solid`, если подходящая иконка уже есть.

## Целевая структура проекта

Новая структура должна быть доменно-ориентированной. Общие элементы лежат в `shared`, доменные модули в `features`, страницы в `routes`.

```text
web/
  index.html
  login.html
  signup.html
  package.json
  vite.config.js
  jsconfig.json
  src/
    app/
      entry.jsx
      entry-login.jsx
      entry-signup.jsx
      routes.jsx
      providers.jsx
    pages/
      LoginPage.jsx
      SignupPage.jsx
    routes/
      dm/
        DmLayout.jsx
        DmSelectPage.jsx
        DmSearchPage.jsx
        DmChatPage.jsx
      server/
        ServerLayout.jsx
        ServerSelectPage.jsx
        ServerCreatePage.jsx
        ServerEditPage.jsx
        ServerManagePage.jsx
        ServerTextTopicPage.jsx
        ServerVoiceTopicPage.jsx
    features/
      auth/
        api.js
        actions.js
        constants.js
        model.js
      dm/
        api.js
        constants.js
        model.js
        components/
      server/
        api.js
        actions.js
        constants.js
        model.js
        components/
      chat/
        api.js
        constants.js
        model.js
        components/
      voice/
        constants.js
        media.js
        call.js
        components/
    shared/
      api/
        client.js
        errors.js
        endpoints.js
      config/
        constants.js
        storage.js
      lib/
        cn.js
        form.js
        route.js
      stores/
        toast.js
        websocket.js
      ui/
        Button.jsx
        Modal.jsx
        TextField.jsx
        Toasts.jsx
        Header.jsx
      assets/
        icon.svg
    styles/
      index.css
```

Правила по структуре:

- `app` содержит только точки входа, конфигурацию роутера и глобальные провайдеры.
- `pages` содержит standalone-страницы вне основного `/app`, например login/signup.
- `routes` содержит route-компоненты. Route-компонент собирает данные, читает params, вызывает actions и компонует фичи, но не содержит низкоуровневый API-код.
- `features/<domain>` содержит доменную логику и компоненты конкретной области: `dm`, `server`, `chat`, `voice`, `auth`.
- `shared/ui` содержит переиспользуемые dumb-компоненты без знания домена.
- `shared/api` содержит общий fetch-клиент, endpoint-константы и типовые ошибки.
- `shared/stores` содержит глобальные singleton stores: toast, websocket, auth session. Доменные сторы должны жить в `features/<domain>`.
- `shared/config` содержит глобальные константы приложения, ключи localStorage, лимиты и default values.
- Нельзя складывать все компоненты в плоскую папку `components`, если компонент относится только к одной фиче.
- Нельзя держать API разных доменов в одном большом `lib/api.js`.
- Нельзя держать route-компоненты с именами `Wrapper`, `Select`, `Edit` без доменного контекста. Имя файла должно быть понятно вне папки.

## Именование файлов

Использовать единый стиль:

- Компоненты: `PascalCase.jsx`.
- Страницы и route-компоненты: `PascalCasePage.jsx` или `PascalCaseLayout.jsx`.
- Хуки/фабрики Solid-состояния: `camelCase.js`, имя начинается с `create`, если функция создает состояние или ресурс.
- API-модули: `api.js` внутри домена или конкретный `<domain>Api.js`, если файл лежит в общей папке.
- Actions Solid Router: `actions.js`.
- Константы: `constants.js`.
- Мапперы/нормализаторы: `model.js`, `mappers.js` или `normalize.js`.
- Утилиты: `camelCase.js`, например `formatDate.js`, `routeParams.js`.
- CSS: глобальные стили только в `styles/index.css`. Component-level CSS не добавлять без сильной причины.

Не использовать:

- `Wrapper.jsx`, `Select.jsx`, `Edit.jsx`, `Search.jsx` без доменного префикса.
- Смешение `snake_case` и `kebab-case` в именах файлов.
- Имена файлов, которые повторяются в разных доменах и плохо читаются в tabs/editor search.

Примеры целевых имен:

```text
routes/server/ServerEditPage.jsx
routes/server/ServerLayout.jsx
features/server/actions.js
features/server/api.js
features/voice/call.js
features/voice/media.js
shared/ui/TextField.jsx
```

## Именование в коде

### Компоненты

- Компоненты называются `PascalCase`: `ChatMessage`, `ServerSidebarItem`, `VoiceTopicGrid`.
- Route-компоненты заканчиваются на `Page`: `DmChatPage`, `ServerCreatePage`.
- Layout-компоненты заканчиваются на `Layout`: `ServerLayout`, `DmLayout`.
- UI-компоненты без домена называются нейтрально: `Button`, `Modal`, `TextField`.
- Внутренние маленькие компоненты в файле допустимы, если используются только в этом файле. Их имена тоже `PascalCase`.

### Функции и переменные

- Обычные функции и переменные: `camelCase`.
- Solid stores/hooks/factories: `createX`, например `createVoiceCall`, `createCallMedia`, `createDmMessages`.
- Accessor-сигналы называются существительными: `user`, `messages`, `isLoading`.
- Setter-сигналы называются `setX`: `setUser`, `setMessages`.
- Boolean-переменные начинаются с `is`, `has`, `can`, `should`: `isLoading`, `hasMore`, `canSubmit`, `shouldScroll`.
- Event handlers называются `handleX`: `handleSubmit`, `handleScroll`, `handleDeleteTopic`.
- Callback props называются `onX`: `onSubmit`, `onClose`, `onMessage`.
- Мапперы называются `mapXToY`, нормализаторы `normalizeX`.
- Fetch/API функции называются по действию: `getServers`, `createServer`, `updateTopic`, `markDmRead`.
- Не использовать `fetchX` для функций, которые не только получают данные или имеют side effects.

### Аббревиатуры

- Аббревиатуры в середине имени пишутся как обычные слова: `dmId`, `serverId`, `topicId`, `userId`.
- Не использовать одновременно `ID`, `Id` и `id`. Целевой стандарт: `id` в одиночку, `serverId`, `topicId`, `dmId`.
- Константы с аббревиатурами: `WS_EVENT_TYPE`, `WS_CHANNEL_TYPE`, если это plain object constants.

### Экспорты

- Компонент файла экспортируется default, если это основной компонент файла.
- Утилиты, actions, API-функции, константы экспортируются named exports.
- Не смешивать default и named export для одной и той же сущности.
- Не экспортировать внутренние helper-функции, если они не используются снаружи.

Плохо:

```js
export function UseServerContext() {}
export const GetServers = query(fetchServers, "user_servers");
const DMsType = "dms";
```

Хорошо:

```js
export function useServerContext() {}
export const getServersQuery = query(getServers, "server.list");
export const CHAT_TARGET = {
  dm: "dm",
  topic: "topic",
};
```

## SolidJS-правила

- Не деструктурировать reactive props напрямую, если это ломает реактивность. Использовать `splitProps`, accessors или функции.
- Props, которые являются accessors, должны явно читаться как функции: `props.user()`.
- Если компонент ожидает accessor, имя prop не меняется, но это должно быть видно из контекста или JSDoc-комментария для сложных shared-компонентов.
- Derived state делать через `createMemo`, а не через ручной `createSignal + createEffect`, если значение полностью вычисляемое.
- Side effects делать в `createEffect`, cleanup всегда регистрировать через `onCleanup`.
- Для async route data использовать `query`/`createAsync`, где это подходит Solid Router.
- Не вызывать async-функции прямо в JSX.
- Не создавать singleton side effects внутри компонента без cleanup.
- Не использовать `document.getElementById` для управления UI, кроме интеграции с нативным API, где нет нормальной альтернативы. Для модалок предпочтительнее controlled state.
- Не использовать `window.location.href` внутри компонентов и API-клиента, кроме hard redirect между отдельными entry pages. Внутри `/app` использовать `useNavigate` или централизованный auth redirect.
- `createRoot` допустим для настоящих singleton stores, например websocket/toast, но не для локального состояния компонента.
- Для списков использовать `<For>`; `<Index>` использовать только когда порядок и длина стабильны, а элементы обновляются по индексу.
- Для условного UI использовать `<Show>`, для множественного ветвления можно использовать `<Switch>/<Match>`.
- Не писать `key` prop в `<For>` как в React. Solid `<For>` не использует `key` таким образом.

## Компоненты

Компоненты делятся на три уровня:

- `shared/ui`: dumb-компоненты. Не знают API, router params, websocket, toast и доменные модели.
- `features/<domain>/components`: доменные компоненты. Могут знать модель домена, но не должны сами ходить в API без явного container-слоя.
- `routes`: container-компоненты. Читают route params, вызывают queries/actions, собирают layout.

Правила:

- Один файл - один основной компонент.
- Компонент длиннее 200-250 строк нужно разделить, если в нем смешаны форма, список, modal, API и state machine.
- JSX должен быть читаемым сверху вниз. Сложные вычисления выносить выше в `createMemo` или helper.
- Компонент не должен одновременно:
  - читать route params;
  - делать raw `fetch`;
  - хранить сложное состояние;
  - рендерить большой UI;
  - показывать toast.
- Shared UI-компоненты должны принимать `class` и аккуратно объединять классы через общий `cn`.
- Для `button` всегда задавать `type`, по умолчанию `button`.
- Интерактивные элементы должны иметь disabled/loading state, если действие async.
- Иконки в icon-only кнопках должны иметь доступное имя через `aria-label` или `title`.
- Не делать hover scale на базовых кнопках, если это двигает layout или ухудшает UX. Hover должен быть предсказуемым.
- Не прятать важные ошибки только в toast. Ошибка формы показывается рядом с формой.

## Роутинг

- Все route definitions держать в `app/routes.jsx`.
- Route params должны иметь единый нейминг: `serverId`, `topicId`, `dmId`.
- Парсинг params делать через helper, например `parseRouteId(params.serverId)`.
- Нельзя использовать `parseInt(..., 10)` в каждом компоненте вручную.
- `matchFilters` держать рядом с роутами в константе `ROUTE_PARAM_FILTERS`.
- Навигационные пути держать в route builder functions:

```js
export const routes = {
  dm: {
    list: () => "/dm",
    chat: (dmId) => `/dm/${dmId}`,
    search: () => "/dm/search",
  },
  server: {
    list: () => "/server",
    create: () => "/server/create",
    edit: (serverId) => `/server/${serverId}/edit`,
    textTopic: (serverId, topicId) => `/server/${serverId}/text/${topicId}`,
    voiceTopic: (serverId, topicId) => `/server/${serverId}/voice/${topicId}`,
  },
};
```

- Не писать URL строками по всему приложению.
- Внутри `/app` пути должны быть относительны base router, например `/server/create`, а не `/app/server/create`.
- External entry redirects (`/login`, `/signup`, `/app/dm`) централизовать в auth/navigation helper.

## API-клиент

Raw `fetch` должен быть только в API-слое. Компоненты и stores не должны напрямую вызывать `fetch`, кроме редких browser API случаев, не связанных с backend.

Целевой слой:

```text
shared/api/client.js
shared/api/errors.js
shared/api/endpoints.js
features/server/api.js
features/dm/api.js
features/auth/api.js
features/chat/api.js
```

Правила:

- Все endpoint paths хранить в `shared/api/endpoints.js`.
- Общий `apiRequest` отвечает за:
  - JSON serialization;
  - JSON parsing;
  - пустой response body;
  - 401;
  - network errors;
  - нормализацию server error message;
  - единый return/error contract.
- API-функции не показывают toast сами. Они возвращают данные или бросают/возвращают типовую ошибку.
- Toast показывает action/container-слой, где понятно, что делал пользователь.
- Не дублировать `handleResponse` в разных файлах.
- Не возвращать разные типы из одной API-функции без явного Result-типа. Например, нельзя иногда вернуть `[]`, иногда `null`, иногда объект.
- Для list endpoints возвращать массив только после нормализации: `return data.servers ?? []`.
- Для mutation endpoints возвращать созданную/обновленную сущность или `{ ok: true }`, но единообразно внутри домена.
- Query params строить через `URLSearchParams`.
- Request body должен использовать backend snake_case, но внутри фронтенда модель должна быть camelCase, если есть mapper.
- Backend DTO и frontend model не смешивать в компонентах.

Пример целевого контракта:

```js
export async function apiRequest(path, options = {}) {
  const response = await fetch(path, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  if (response.status === 401) {
    throw new AuthRequiredError();
  }

  const data = await parseJsonResponse(response);

  if (!response.ok) {
    throw ApiError.fromResponse(response, data);
  }

  return data;
}
```

## Ошибки

Ошибки должны быть предсказуемыми для пользователя и для кода.

### Типы ошибок

Минимальный набор:

- `ApiError`: backend вернул не-2xx.
- `AuthRequiredError`: backend вернул 401.
- `NetworkError`: запрос не дошел до сервера.
- `ValidationError`: ошибка валидации формы на клиенте.
- `MediaDeviceError`: ошибка камеры/микрофона/шеринга экрана.
- `WebSocketError`: ошибка websocket protocol/state.

### Правила обработки

- Не писать `console.log(err)` в production-коде. Если лог нужен, использовать единый debug/logger helper и не логировать чувствительные данные.
- Не глотать ошибки пустым `catch {}`. Допустимо только рядом с комментарием, почему ошибка безопасна и ожидаема.
- Не показывать пользователю `err.message` напрямую, если это техническая/browser ошибка. Сначала нормализовать.
- Ошибки формы показывать inline под конкретной формой или полем.
- Ошибки фоновой синхронизации можно показывать toast.
- Ошибки авторизации обрабатываются централизованно: reset auth state + redirect на login.
- Ошибки websocket не должны ломать весь UI. Показывать статус подключения и дать возможность reconnect.
- Ошибки media devices показывать рядом с voice UI, а не только toast.
- Async action всегда должен иметь loading/pending state и не допускать повторный submit.
- Любой `finally` должен восстанавливать loading state.
- Текст ошибок хранить в константах, если он повторяется.

Плохо:

```js
catch (err) {
  console.log(err);
  toast.error("Network error – please check your connection");
  return null;
}
```

Хорошо:

```js
catch (error) {
  throw normalizeApiError(error, ERROR_MESSAGES.network);
}
```

А toast:

```js
try {
  await createServer(values);
  toast.success(SERVER_MESSAGES.created);
} catch (error) {
  setFormError(getUserErrorMessage(error));
}
```

## Константы

Константы должны лежать максимально близко к месту использования, но не внутри компонента, если они не зависят от props/state.

Где хранить:

- Глобальные app constants: `shared/config/constants.js`.
- LocalStorage/sessionStorage keys: `shared/config/storage.js`.
- API endpoints: `shared/api/endpoints.js`.
- Route builders and route filters: `app/routes.js` или `shared/config/routes.js`.
- Доменные enum-like значения: `features/<domain>/constants.js`.
- UI variants shared-компонента: рядом с компонентом или в том же файле, если это private map.
- Повторяющиеся user-facing messages: в constants конкретного домена.
- Magic numbers для UI/логики: именованные константы рядом с модулем.

Именование:

- Immutable enum-like objects: `UPPER_SNAKE_CASE`.
- Primitive constants module-scope: `UPPER_SNAKE_CASE`, если это настоящая константа.
- Private maps рядом с компонентом: `camelCase` допустим, если это implementation detail, например `buttonVariantClass`.
- LocalStorage keys: `STORAGE_KEYS`.
- Error messages: `ERROR_MESSAGES`, `SERVER_ERROR_MESSAGES`.
- Таймауты: имя содержит единицу измерения, например `TOAST_TTL_MS`, `RECONNECT_DELAY_MS`, `SCROLL_BOTTOM_THRESHOLD_PX`.

Пример:

```js
export const STORAGE_KEYS = {
  voiceCameraId: "chattery.voice.camera_id",
  voiceMicId: "chattery.voice.mic_id",
  voiceSpeakerId: "chattery.voice.speaker_id",
};

export const WS_EVENT_TYPE = {
  ping: "ping",
  pong: "pong",
  join: "join",
  leave: "leave",
  message: "message",
  error: "error",
};

export const TOAST_TTL_MS = 10_000;
export const RECONNECT_DELAY_MS = 1_200;
```

Запрещено:

- Магические строки endpoint-ов в компонентах.
- Магические route strings в UI.
- Числа `50`, `80`, `10000`, `1200` без имени.
- Константы с разным стилем в одном домене: `DMsType`, `ServersType`, `TopicTypeText`.

## Состояние и stores

- Local UI state хранить в компоненте через `createSignal`.
- Derived state через `createMemo`.
- Общедоступное состояние домена выносить в feature store или Solid Router query.
- Глобальными singleton stores должны быть только действительно глобальные вещи: auth session, toast, websocket connection.
- Store не должен напрямую показывать toast, если это не сам toast store.
- Store не должен знать о route params. Params передаются снаружи.
- Для коллекций использовать `createStore`, если нужны точечные обновления вложенных данных; иначе достаточно `createSignal`.
- Не хранить один и тот же server/dm/topic state в нескольких местах без правила синхронизации.
- Нельзя смешивать resource/query state и ручной loading state для одного и того же запроса без необходимости.
- У каждого async процесса должны быть понятные состояния: `idle`, `loading`, `ready`, `error` или доменный state machine.

## Forms и actions

- Формы с backend mutation должны использовать один из двух подходов:
  - Solid Router `action` для route-level forms;
  - controlled submit handler для локальных interactive forms.
- В одном компоненте не смешивать `action`, `useAction`, `useSubmission` и ручной `fetch` для одинаковых операций.
- Валидацию формы делать до API-запроса.
- Все значения из `FormData` нормализовать: trim строк, Number для id, проверка required.
- Submit button disabled, пока идет отправка.
- Ошибку формы показывать inline.
- `confirm()` не использовать для важных destructive actions в новом коде. Использовать общий `ConfirmDialog`.
- Destructive action требует явного подтверждения и понятного текста.
- После успешной mutation инвалидировать/refetch соответствующую query или обновить store по единому правилу.

## WebSocket

- WebSocket protocol constants хранятся в `features/realtime/constants.js` или `shared/stores/websocket.js`, если store остается единственным владельцем протокола.
- В приложении должен быть один app-level websocket client, если backend protocol подразумевает одну активную сессию.
- Subscribe-функции всегда возвращают unsubscribe.
- Любой subscribe в компоненте должен вызываться внутри `createEffect`/`onMount` и чиститься через `onCleanup`.
- Join/leave channel должны быть идемпотентными.
- Pending events должны иметь лимит или стратегию очистки.
- Payload должен парситься и нормализоваться в одном месте.
- UI-компоненты не должны вручную собирать websocket event objects.
- Ошибки websocket отображаются через status/error state, а не только toast.

## Voice/WebRTC/media

- Media state и call signaling должны быть разделены: `media.js` отвечает за устройства и local streams, `call.js` за peer connection/signaling.
- Все browser capabilities проверять feature detection-ом: `navigator.mediaDevices`, `setSinkId`, `getDisplayMedia`.
- Все tracks должны останавливаться в cleanup.
- Ошибки permissions/devices нормализовать в user-friendly сообщения.
- Device ids хранить через `STORAGE_KEYS`, не строками в модуле.
- Magic delays для ICE batching, reconnect и т.п. именовать константами.
- Сложные WebRTC state transitions должны быть покрыты комментариями только там, где без комментария трудно понять инвариант.
- Не смешивать UI рендера участников и peer connection logic в одном файле.

## Styling и Tailwind

- Базовый стиль строится на Tailwind utilities.
- Глобальные custom utilities объявляются только в `styles/index.css`.
- Повторяющиеся наборы классов выносить в shared UI-компонент или class map, а не копировать в каждом компоненте.
- Для склейки классов использовать общий `cn(...classes)`.
- `classList` использовать для условных классов, если это читабельнее.
- Не использовать inline style, кроме динамических CSS variables или API, где Tailwind не подходит.
- Не добавлять arbitrary values без причины. Если значение повторяется, сделать константу или theme utility.
- Не использовать `tracking-wider` повсюду по умолчанию. Типографика должна быть осознанной и читабельной.
- Не использовать hover transform, который меняет layout или создает дерганье.
- Не делать вложенные карточки без необходимости.
- Все элементы должны корректно работать на узких экранах: не использовать `w-md`, `min-w-sm` без responsive fallback.
- Интерактивные элементы должны иметь focus state.
- Цвет variant-компонентов должен задаваться через variant map, а не ad hoc классами в местах использования.

## Accessibility

- У каждого input должен быть связанный `label`.
- Icon-only button обязан иметь `aria-label`.
- Модалка должна иметь понятный title, закрытие по Escape и фокус внутри модалки, если используется custom modal.
- Цвет не должен быть единственным способом передать состояние.
- Disabled state должен быть видимым.
- Loading state должен быть объявлен текстом или `aria-busy`, если действие длительное.
- Не использовать кликабельный `div`, если есть `button` или `a`.
- Ссылки используются для навигации, кнопки для действий.

## Assets

- SVG и статические файлы хранятся в `shared/assets` или `src/assets`, если пока нет миграции.
- Не импортировать asset из `dist`.
- Имена assets: `kebab-case`, например `app-icon.svg`.
- Доменные assets хранятся рядом с фичей только если не используются глобально.

## Комментарии

- Комментарии нужны для инвариантов, протоколов, browser quirks и сложной async-синхронизации.
- Не комментировать очевидный JSX или простые присваивания.
- Если нужен длинный комментарий, возможно код нужно разбить на функции с говорящими именами.
- Комментарии пишутся на английском, как и UI-код.

## JSDoc и типизация в JavaScript

Проект сейчас использует JavaScript, поэтому JSDoc нужен не для украшения, а для явных контрактов там, где без типов легко ошибиться. JSDoc должен помогать IDE и человеку понять форму данных, допустимые значения и side effects.

### Когда JSDoc обязателен

Добавлять JSDoc нужно для публичных контрактов модуля:

- Exported API-функции, которые вызывают backend.
- Exported Solid stores/factories: `createVoiceCall`, `createCallMedia`, `createDmMessages`.
- Shared UI-компоненты, если props не очевидны или включают accessors/callbacks.
- Route helpers и builders, если они принимают ids/params.
- Мапперы между backend DTO и frontend model.
- Функции, которые возвращают union/result object: `{ ok, error }`, `{ data, cursor }`, state machine status.
- WebSocket/WebRTC/media protocol payloads.
- Объекты констант, которые являются enum-like значениями.
- Нестандартные callbacks: `onMessage`, `onVolume`, `subscribeError`, `sendEvent`.
- Функции с важными side effects: redirect, localStorage, media tracks, websocket join/leave.

Пример:

```js
/**
 * @typedef {Object} ServerTopic
 * @property {number} id
 * @property {string} name
 * @property {"text" | "voice"} type
 */

/**
 * @typedef {Object} Server
 * @property {number} id
 * @property {string} name
 * @property {ServerTopic[]} topics
 */

/**
 * Loads servers available to the current user.
 *
 * @returns {Promise<Server[]>}
 * @throws {ApiError}
 * @throws {AuthRequiredError}
 * @throws {NetworkError}
 */
export async function getServers() {
  return apiRequest(API_ENDPOINTS.server.list);
}
```

### Когда JSDoc желателен

JSDoc стоит добавить, если код остается на JavaScript и тип нельзя быстро вывести из имени:

- Props большого компонента.
- Значение `props`, которое должно быть Solid accessor.
- Массивы/словари сложных объектов.
- Конфигурационные объекты.
- Private helper, если он работает с неочевидным browser API.
- Magic protocol shape, даже если helper не export-ится.

Пример для Solid component props:

```js
/**
 * @typedef {Object} VoiceTopicGridProps
 * @property {ReturnType<typeof createCallMedia>} media
 * @property {ReturnType<typeof createVoiceCall>} call
 * @property {User | null} user
 * @property {import("solid-js").Accessor<Record<number, number>>} volumes
 * @property {(userId: number, volume: number) => void} onVolumeChange
 */

/**
 * @param {VoiceTopicGridProps} props
 */
export default function VoiceTopicGrid(props) {
  // ...
}
```

### Когда JSDoc не нужен

Не добавлять JSDoc к очевидному коду:

- Маленькие local handlers: `handleSubmit`, `handleClose`, `handleScroll`, если их контракт ясен из JSX.
- Внутренние компоненты на 5-20 строк с простыми `children`.
- Простые signals: `const [isLoading, setIsLoading] = createSignal(false)`.
- Очевидные helpers: `cn`, `trimValue`, `parseJson`, если сигнатура простая.
- Функции, которые скоро будут удалены при рефакторинге и не являются контрактом.
- Дублирование имени функции словами: `/** Gets servers */ function getServers()`.

Плохо:

```js
/**
 * Handles submit.
 *
 * @param {Event} event
 */
const handleSubmit = (event) => {
  event.preventDefault();
};
```

Хорошо: оставить без JSDoc или вынести сложную часть в типизированный helper.

### Как писать JSDoc

- JSDoc пишется на английском.
- Описывать контракт, не реализацию.
- Сначала `@typedef`, затем функции, которые его используют.
- Имена typedef должны быть `PascalCase`: `Server`, `DmMessage`, `ApiErrorPayload`.
- Для enum-like значений использовать union literals или `@readonly`.
- Для Solid accessors использовать `import("solid-js").Accessor<T>`.
- Для callbacks писать полную сигнатуру: `(value: string) => void`.
- Для async функций всегда указывать `Promise<T>`.
- Если функция может бросить типовую ошибку, указать `@throws`.
- Если функция имеет side effect, написать одну короткую строку в описании.
- Не использовать `any`, кроме временного adapter-слоя. Если нужен `any`, рядом должен быть TODO на нормализацию модели.

Пример enum-like constants:

```js
/**
 * @readonly
 * @enum {string}
 */
export const WS_CHANNEL_TYPE = {
  dm: "dm",
  textTopic: "text_topic",
  voiceTopic: "voice_topic",
};
```

Пример result object:

```js
/**
 * @typedef {Object} ActionResult
 * @property {boolean} ok
 * @property {string=} error
 */

/**
 * @param {number} topicId
 * @param {FormData} formData
 * @returns {Promise<ActionResult>}
 */
export async function updateTopicActionHandler(topicId, formData) {
  // ...
}
```

### Где хранить typedef

- Типы, используемые только в одном файле, объявлять в начале этого файла после imports.
- Доменные типы, используемые в нескольких файлах, держать в `features/<domain>/model.js`.
- Shared API/error типы держать в `shared/api/errors.js` или `shared/api/client.js`.
- UI prop typedef держать рядом с компонентом, если props не используются снаружи.
- Не создавать один глобальный `types.js` для всего приложения.

### JSDoc и будущий TypeScript

JSDoc должен приближать проект к TypeScript, а не заменять архитектуру:

- Типизировать границы модулей и данные backend/frontend.
- Не пытаться JSDoc-ом компенсировать слишком большой компонент.
- Если JSDoc становится сложнее функции, нужно разбить функцию или ввести явную модель.
- При будущей миграции на TypeScript JSDoc typedef из `model.js` должен легко превращаться в `type`/`interface`.

## Тестируемость

Даже если тестовая инфраструктура еще не настроена, код должен быть написан так, чтобы его можно было тестировать:

- API-клиент отделен от компонентов.
- Мапперы и валидаторы - чистые функции.
- Сложные state machines вынесены из JSX.
- Компоненты получают данные через props, где это возможно.
- Browser APIs обернуты в маленькие функции, которые можно заменить в тестах.
- Для багфиксов в общей логике добавлять unit tests после появления test runner.
- Для критичных пользовательских сценариев нужны e2e: login, signup, server CRUD, DM, text topic, voice join.

## Анти-паттерны, которые нельзя повторять

- Большой общий `lib/api.js` со всеми endpoint-ами приложения.
- Дублирование API-функций в разных файлах.
- API-функции, которые одновременно делают request, показывают toast и управляют redirect.
- Компоненты с raw `fetch`.
- Разный нейминг id: `dmID`, `serverID`, `topicId`, `topicID`.
- Компоненты/файлы `Wrapper`, `Select`, `Edit`, `Search` без доменного имени.
- PascalCase для обычных функций: `UseServerContext`, `GetServers`.
- `console.log(err)` в catch.
- Пустой `catch {}` без объяснения.
- Магические строки URL, endpoint, localStorage keys и websocket event types внутри компонентов.
- `window.location.href` в глубине приложения вместо централизованной навигации.
- `document.getElementById(...).hidePopover()` из компонента вместо управляемого состояния.
- `confirm()` для destructive actions.
- Смешение route layer, API layer, store layer и UI в одном файле.
- Компонент, который сам решает, когда показывать глобальный toast, если ошибка относится к форме.
- Сложная WebRTC/WebSocket логика внутри UI-компонента.
- Tailwind-классы, скопированные в десятках input/button без общего UI-компонента.
- Неименованные timeout/delay/threshold numbers.
- Импорт через `../..` при наличии `~/`.

## Правила для при генерации нового кода

Перед изменением фронтенда LLM должна:

1. Определить домен изменения: `auth`, `dm`, `server`, `chat`, `voice`, `shared`.
2. Положить новый код в правильный слой: route, feature, shared ui, shared api, config.
3. Проверить, нет ли уже подходящего shared UI/API/helper.
4. Использовать существующий `~/` alias.
5. Не расширять плоскую папку `components`, если компонент доменный.
6. Не добавлять raw `fetch` в компонент.
7. Не добавлять новую строковую константу endpoint/route/storage key прямо в JSX.
8. Сохранять camelCase для frontend-моделей и делать mapper для backend snake_case.
9. Возвращать из API единый тип результата или бросать типовую ошибку.
10. Показывать ошибку на правильном уровне: inline для формы, toast для фонового события, status block для long-running connection/media.
11. Добавлять cleanup для subscriptions, timers, media tracks и websocket listeners.
12. Не менять unrelated код и не приводить весь файл к новому стилю без необходимости.

При рефакторинге старого кода:

- Сначала убрать дубли и централизовать API/error handling.
- Затем переименовать файлы/компоненты до доменно-понятных имен.
- Затем разнести route/container/UI.
- Затем нормализовать constants и route builders.
- Затем улучшать accessibility и responsive styling.

## Минимальный чеклист перед merge

- Нет raw `fetch` в компонентах.
- Нет новых endpoint strings вне `shared/api/endpoints.js`.
- Нет новых route strings вне route builders.
- Нет новых localStorage keys вне `STORAGE_KEYS`.
- Нет `console.log`/`console.error` в пользовательском коде без logger policy.
- Все async actions имеют loading/pending state.
- Ошибки не глотаются и показываются на правильном уровне.
- Все subscriptions/timers/media resources имеют cleanup.
- Имена файлов и exports соответствуют правилам.
- Компонент не смешивает больше одного уровня ответственности.
- UI доступен с клавиатуры, inputs имеют labels, icon buttons имеют `aria-label`.
- `npm run build` проходит.
