# LongMusic

Самостоятельно хостящийся музыкальный стриминг-сервис на Go: загрузка треков, прослушивание, плейлисты, избранное и чарты - с собственным веб-плеером.

## Возможности

- Регистрация и вход, JWT-аутентификация
- Загрузка MP3-треков с автоматическим определением названия, исполнителя и длительности из ID3-тегов
- Стриминг треков и отдача обложки альбома (извлекается из ID3-тега)
- Поиск треков, список треков по исполнителю, страницы исполнителей
- Плейлисты: создание, добавление и удаление треков, публикация (публичный/приватный), поиск публичных плейлистов
- Избранное: добавление и удаление треков
- История прослушиваний
- Чарты: топ треков за всё время и топ треков за сегодня
- Профиль пользователя: публичный профиль, смена имени пользователя и аватара
- Встроенный веб-плеер (статические HTML/CSS/JS) с поддержкой offline через Service Worker

## Технологии

- **Go 1.25**
- **go-chi/chi** - HTTP-роутер
- **PostgreSQL** (sqlx + lib/pq) - хранилище данных
- **Redis** (go-redis) - кэш
- **golang-jwt/jwt** - JWT-токены
- **bogem/id3v2**, **dhowden/tag**, **tcolgate/mp3** - чтение ID3-тегов, обложек и подсчёт длительности треков
- **Docker / Docker Compose** - контейнеризация
- Статический фронтенд (HTML/CSS/JS) и Service Worker, отдаваемые тем же Go-сервером

## Основные эндпоинты

| Метод | Путь | Описание | Авторизация |
| --- | --- | --- | --- |
| POST | `/api/auth/register` | Регистрация | - |
| POST | `/api/auth/login` | Вход, получение JWT | - |
| GET | `/api/me` | Данные текущего пользователя | JWT |
| GET | `/api/tracks` | Список треков | - |
| GET | `/api/tracks/search` | Поиск треков | - |
| GET | `/api/tracks/{id}` | Трек по ID | - |
| GET | `/api/tracks/{id}/stream` | Стриминг аудио | - |
| POST | `/api/tracks` | Загрузить трек | JWT |
| GET | `/api/cover/{id}` | Обложка трека | - |
| GET | `/api/artists` | Список исполнителей | - |
| GET | `/api/artists/{id}/tracks` | Треки исполнителя | - |
| GET | `/api/charts/{id}` | Чарт (например, топ за всё время / за сегодня) | - |
| GET | `/api/playlists` | Плейлисты пользователя | JWT |
| POST | `/api/playlists` | Создать плейлист | JWT |
| GET | `/api/playlists/{id}` | Плейлист по ID | опционально |
| GET | `/api/playlists/{id}/tracks` | Треки плейлиста | опционально |
| POST | `/api/playlists/{id}/tracks` | Добавить трек в плейлист | JWT |
| DELETE | `/api/playlists/{id}/tracks` | Удалить трек из плейлиста | JWT |
| PATCH | `/api/playlists/{id}/publish` | Сделать плейлист публичным/приватным | JWT |
| GET | `/api/playlists/public/search` | Поиск публичных плейлистов | - |
| POST | `/api/favorites` | Добавить в избранное | JWT |
| DELETE | `/api/favorites/{id}` | Убрать из избранного | JWT |
| GET | `/api/favorites` | Список избранного | JWT |
| POST | `/api/plays` | Записать факт прослушивания | JWT |
| GET | `/api/plays` | История прослушиваний | JWT |
| GET | `/api/profile` | Профиль текущего пользователя | JWT |
| GET | `/api/users/{login}/profile` | Публичный профиль пользователя | - |
| PATCH | `/api/profile/username` | Сменить имя пользователя | JWT |
| POST | `/api/profile/avatar` | Сменить аватар | JWT |
| GET | `/ping` | Проверка работоспособности сервиса | - |

## Установка и запуск

### 1. Клонирование репозитория

```
git clone https://github.com/TheLonger011/LongMusic.git
cd LongMusic
```

### 2. Параметры подключения

Сейчас строка подключения к PostgreSQL задана прямо в `cmd/api/main.go`:

```
postgres://postgres:181818@localhost:5432/longmusic?sslmode=disable
```

При локальном запуске без Docker либо поднимите базу с такими же учётными данными, либо поменяйте строку подключения под себя. Для запуска через Docker Compose ничего менять не нужно - сервис `db` уже настроен на эти значения.

### 3. Запуск через Docker Compose

```
docker-compose up -d
```

Поднимутся приложение (порт `8080`), PostgreSQL и Redis. Загруженные файлы сохраняются в volume `uploads_data`.

### 4. Применение миграций

Миграции лежат в `migrations/`. Примените их вручную через [golang-migrate](https://github.com/golang-migrate/migrate) или другой инструмент, например:

```
migrate -path ./migrations -database "postgres://postgres:181818@localhost:5432/longmusic?sslmode=disable" up
```

### 5. Запуск локально

```
go run ./cmd/api/main.go
```

Приложение стартует на порту `8080` и раздаёт веб-плеер по адресу `http://localhost:8080`.

---

# LongMusic (English)

A self-hosted music streaming service written in Go: track uploads, playback, playlists, favorites, and charts - with a built-in web player.

## Features

- Registration and login, JWT authentication
- MP3 track upload with automatic title, artist, and duration detection from ID3 tags
- Track streaming and album cover retrieval (extracted from the ID3 tag)
- Track search, per-artist track listings, artist pages
- Playlists: create, add/remove tracks, publish (public/private), search public playlists
- Favorites: add and remove tracks
- Listening history
- Charts: all-time top tracks and today's top tracks
- User profile: public profile, username and avatar updates
- Built-in web player (static HTML/CSS/JS) with offline support via a Service Worker

## Technologies

- **Go 1.25**
- **go-chi/chi** - HTTP router
- **PostgreSQL** (sqlx + lib/pq) - data storage
- **Redis** (go-redis) - caching
- **golang-jwt/jwt** - JWT tokens
- **bogem/id3v2**, **dhowden/tag**, **tcolgate/mp3** - reading ID3 tags/covers and computing track duration
- **Docker / Docker Compose** - containerization
- A static frontend (HTML/CSS/JS) and Service Worker, served by the same Go server

## Main Endpoints

| Method | Path | Description | Auth |
| --- | --- | --- | --- |
| POST | `/api/auth/register` | Register | - |
| POST | `/api/auth/login` | Log in, receive a JWT | - |
| GET | `/api/me` | Current user info | JWT |
| GET | `/api/tracks` | List tracks | - |
| GET | `/api/tracks/search` | Search tracks | - |
| GET | `/api/tracks/{id}` | Get a track by ID | - |
| GET | `/api/tracks/{id}/stream` | Stream audio | - |
| POST | `/api/tracks` | Upload a track | JWT |
| GET | `/api/cover/{id}` | Get a track's cover art | - |
| GET | `/api/artists` | List artists | - |
| GET | `/api/artists/{id}/tracks` | Get an artist's tracks | - |
| GET | `/api/charts/{id}` | Get a chart (e.g. all-time / today's top) | - |
| GET | `/api/playlists` | Current user's playlists | JWT |
| POST | `/api/playlists` | Create a playlist | JWT |
| GET | `/api/playlists/{id}` | Get a playlist by ID | optional |
| GET | `/api/playlists/{id}/tracks` | Get a playlist's tracks | optional |
| POST | `/api/playlists/{id}/tracks` | Add a track to a playlist | JWT |
| DELETE | `/api/playlists/{id}/tracks` | Remove a track from a playlist | JWT |
| PATCH | `/api/playlists/{id}/publish` | Set a playlist public/private | JWT |
| GET | `/api/playlists/public/search` | Search public playlists | - |
| POST | `/api/favorites` | Add to favorites | JWT |
| DELETE | `/api/favorites/{id}` | Remove from favorites | JWT |
| GET | `/api/favorites` | List favorites | JWT |
| POST | `/api/plays` | Record a play | JWT |
| GET | `/api/plays` | Listening history | JWT |
| GET | `/api/profile` | Current user's profile | JWT |
| GET | `/api/users/{login}/profile` | Public user profile | - |
| PATCH | `/api/profile/username` | Update username | JWT |
| POST | `/api/profile/avatar` | Update avatar | JWT |
| GET | `/ping` | Health check | - |

## Installation and Running

### 1. Clone the repository

```
git clone https://github.com/TheLonger011/LongMusic.git
cd LongMusic
```

### 2. Connection settings

The PostgreSQL connection string is currently hardcoded in `cmd/api/main.go`:

```
postgres://postgres:181818@localhost:5432/longmusic?sslmode=disable
```

If running locally without Docker, either spin up a database with matching credentials or edit the connection string. No changes are needed for Docker Compose - the `db` service is already configured to match.

### 3. Run with Docker Compose

```
docker-compose up -d
```

This starts the app (port `8080`), PostgreSQL, and Redis. Uploaded files are stored in the `uploads_data` volume.

### 4. Apply migrations

Migrations live in `migrations/`. Apply them manually with [golang-migrate](https://github.com/golang-migrate/migrate) or a similar tool, e.g.:

```
migrate -path ./migrations -database "postgres://postgres:181818@localhost:5432/longmusic?sslmode=disable" up
```

### 5. Run locally

```
go run ./cmd/api/main.go
```

The app starts on port `8080` and serves the web player at `http://localhost:8080`.
