# Frontend (Docker) для вашего Go-бэкенда

Этот фронтенд реализует CRUD по расходам (expenses) под ваши эндпоинты:

- `GET /expenses`
- `POST /expense`
- `PUT /expense`
- `DELETE /expense` (тело `{ "id": <number> }`)

## Почему без CORS

В контейнере фронта используется Nginx, который:
- отдаёт SPA (React build),
- проксирует запросы `/api/*` на бэкенд.

Поэтому браузер видит один origin (например, `http://localhost:3000`) и CORS на Go-сервере не требуется.

## Быстрый старт (вся связка через docker compose)

1) Скопируйте папку `frontend/` к себе в репозиторий.

2) В `docker-compose.yml`:
- либо оставьте сервис `backend` (если у вас есть `./backend/Dockerfile`),
- либо удалите `backend` и укажите в `frontend/nginx.conf` правильный `proxy_pass`.

3) Запуск:

```bash
docker compose up --build
```

Откройте: http://localhost:3000

## Фронт отдельно (если бэкенд уже запущен)

Если бэкенд работает, например, на `http://localhost:8080`, есть два варианта:

### Вариант A (рекомендуется): прокси через Nginx
Оставьте `VITE_API_BASE="/api"` (по умолчанию) и измените `proxy_pass` в `frontend/nginx.conf`, например:

```
proxy_pass http://host.docker.internal:8080;
```

### Вариант B: прямой вызов API (нужен CORS на бэкенде)
Сборка фронта с переменной:

```bash
cd frontend
VITE_API_BASE=http://localhost:8080 npm run build
```

Но тогда Go должен отдавать корректные CORS-заголовки.

## Расширение на payments
В вашем репозитории уже есть методы `AddPayment/UpdatePayment/DeletePayment/GetAllPayments`,
но в роутере их нет. Чтобы добавить UI под платежи — добавьте эндпоинты в Go и я дам готовую страницу/таблицу.
