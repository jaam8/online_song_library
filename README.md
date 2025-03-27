# online_song_library

[![wakatime](https://wakatime.com/badge/user/018badf6-44ca-4a0f-82e9-9b27db08764a/project/86150661-2ef5-4ed6-8382-d5601d8bde98.svg)](https://wakatime.com/badge/user/018badf6-44ca-4a0f-82e9-9b27db08764a/project/86150661-2ef5-4ed6-8382-d5601d8bde98)  
Проект представляет собой API для работы с библиотекой песен.  
Сервис реализован на Go и использует БД PostgreSQL и Swagger-документацией   
Миграции для базы данных находятся в каталоге `/db/migrations`.

<details><summary>Техническое задание</summary>
1. Выставить rest методы:

  - Получение данных библиотеки с фильтрацией по всем полям и пагинацией  
  - Получение текста песни с пагинацией по куплетам  
  - Удаление песни  
  - Изменение данных песни  
  - Добавление новой песни в формате JSON  

```json
{
  "group": "Muse",
  "song": "Supermassive Black Hole"
}
```

2. При добавлении сделать запрос в АПИ, описанного сваггером.  
   API, описанный сваггером, будет поднят при проверке тестового задания.  
   Реализовывать его отдельно не нужно.

   <details><summary>Описанное API</summary>

   ```yaml
   openapi: 3.0.3
   info:
     title: Music info
     version: 0.0.1
   paths:
     /info:
       get:
         parameters:
           - name: group
             in: query
             required: true
             schema:
               type: string
           - name: song
             in: query
             required: true
             schema:
               type: string
         responses:
           '200':
             description: Ok
             content:
               application/json:
                 schema:
                   $ref: '#/components/schemas/SongDetail'
           '400':
             description: Bad request
           '500':
             description: Internal server error
   components:
     schemas:
       SongDetail:
         required:
           - releaseDate
           - text
           - link
         type: object
         properties:
           releaseDate:
             type: string
             example: "16.07.2006"
           text:
             type: string
             example: "Ooh baby, don't you know I suffer?\\nOoh baby, can you hear me moan?\\nYou caught me under false pretenses\\nHow long before you let me go?\\n\\nOoh\\nYou set my soul alight\\nOoh\\nYou set my soul alight"
           link:
             type: string
             example: "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
    ```
   </details>

3. Обогащенную информацию положить в БД postgres (структура БД должна
   быть создана путем миграций при старте сервиса)
4. Покрыть код debug- и info-логами
5. Вынести конфигурационные данные в .env-файл
6. Сгенерировать сваггер на реализованное АПИ
</details>

## Структура проекта

```
online_song_library
├── cmd
│   └── main.go               # Точка входа в приложение
├── db
│   └── migrations            # Миграции для базы данных
│       ├── 000001_init.down.sql
│       └── 000001_init.up.sql
├── docker-compose.yml        # Конфигурация Docker Compose
├── Dockerfile                # Dockerfile для сборки контейнера
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml          # Swagger-документация API
├── .env.example              # Пример файла конфигурации переменных окружения
├── go.mod                    # Файл зависимостей Go
├── go.sum                    # Контрольные суммы зависимостей
├── internal                  # Внутренняя логика сервиса
│   ├── api                   # Обработчики запросов
│   │   ├── middleware.go
│   │   └── song_handler.go
│   ├── config                # Конфигурации приложения
│   │   └── config.go
│   ├── models                # Описание моделей данных
│   │   └── song.go
│   ├── repository            # Логика работы с базой данных
│   │   └── song_repo.go
│   └── service               # Бизнес-логика
│       └── song_service.go
├── pkg                       # Вспомогательные модули
│   ├── logger                # Логирование
│   │   └── logger.go
│   └── postgres              # Подключение к базе данных
│       └── postgres.go
├── README.md                 # Основной файл с документацией
├── test_for_online_song_library.json  # Тесты для Postman
└── wait-for-it.sh            # Скрипт ожидания запуска зависимостей
```

## Требования

- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/) и [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Настройка

1. Скопируйте файл `.env.example` и настройте переменные окружения по своему усмотрению:
```bash
cp .env.example .env
```

| Переменная           | Значение по умолчанию       | Описание                              |
|----------------------|-----------------------------|---------------------------------------|
| `REST_PORT`          | `8080`                      | Порт, на котором будет доступно API   |
| `POSTGRES_USER`      | `root`                      | Логин пользователя базы данных        |
| `POSTGRES_PASSWORD`  | `1234`                      | Пароль пользователя базы данных       |
| `POSTGRES_DB`        | `postgres`                  | Название базы данных                  |
| `POSTGRES_HOST`      | `postgres`                  | Хост базы данных                      |
| `POSTGRES_PORT`      | `5432`                      | Порт базы данных                      |
| `LOG_LEVEL`          | `info`                      | Уровень логирования (`debug`, `info`) |
| `SWAGGER_URL`        |                             | URL для получения информации о песне  |
| `PATH_TO_MIGRATIONS` | `file:///app/db/migrations` | Путь к миграциям для базы данных      |

2. Убедитесь, что путь к миграциям указан верно:
    - В Docker используется `file:///app/db/migrations`
    - В локальном окружении можно использовать `file://db/migrations`

## Запуск с Docker

Для сборки и запуска контейнеров выполните:

```bash
docker-compose build
docker-compose up
```

## Миграции

При старте приложения автоматически запускаются миграции базы данных.  
Если миграции не применяются, проверьте правильность пути в переменной `PATH_TO_MIGRATIONS`.

## Тестирование API через Postman

- Импортируйте файл `test_for_online_song_library.json` в Postman.
- Задайте переменные окружения:
    - `base_url` — URL вашего API (например, `http://localhost:8080/api/v1/songs`).

## Swagger-документация

Swagger-документация доступна для вашего API, чтобы облегчить взаимодействие с сервисом.  
Документацию можно просматривать по URL, который зависит от вашей конфигурации, например, 
`http://localhost:8080/swagger/index.html`.

