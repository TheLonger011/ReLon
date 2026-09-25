# ReLon

Backend API для Reddit-подобной платформы: посты, сообщества, комментарии и голосования. Написан на Go.

## Возможности

- Регистрация и вход с подтверждением email-кодом
- JWT-аутентификация
- Посты: создание, получение, поиск, обновление, удаление
- Голосование за посты (поставить голос / снять голос)
- Комментарии к постам: создание, получение, обновление, удаление
- Сообщества: создание, обновление, удаление, список, просмотр по ID
- Заявки на вступление в сообщество: подача, одобрение, отклонение владельцем
- Выход участника из сообщества
- Rate limiting по IP на чувствительных эндпоинтах (регистрация, вход, верификация, создание постов и комментариев)

## Технологии

- **Go 1.26**
- **go-chi/chi** - HTTP-роутер
- **PostgreSQL** (jackc/pgx) - основное хранилище данных
- **Redis** (go-redis) - хранение кодов email-верификации
- **golang-jwt/jwt** - JWT-токены
- **golang-migrate** - миграции базы данных
- **net/smtp (Gmail SMTP)** - отправка писем с кодом подтверждения
- **Docker / Docker Compose** - контейнеризация

## Основные эндпоинты

| Метод | Путь | Описание | Авторизация |
| --- | --- | --- | --- |
| POST | `/register` | Регистрация пользователя | - |
| POST | `/login` | Вход, получение JWT | - |
| POST | `/verify` | Подтверждение email по коду | - |
| GET | `/me` | Данные текущего пользователя | JWT |
| GET | `/posts` | Список постов | - |
| GET | `/posts/search` | Поиск постов | - |
| GET | `/posts/{id}` | Пост по ID | - |
| POST | `/posts` | Создать пост | JWT |
| PUT | `/posts/{id}` | Обновить пост | JWT |
| DELETE | `/posts/{id}` | Удалить пост | JWT |
| POST | `/posts/{id}/vote` | Проголосовать за пост | JWT |
| DELETE | `/posts/{id}/vote` | Снять голос | JWT |
| GET | `/posts/{id}/comments` | Комментарии к посту | - |
| POST | `/posts/{id}/comments` | Добавить комментарий | JWT |
| PUT | `/comments/{id}` | Обновить комментарий | JWT |
| DELETE | `/comments/{id}` | Удалить комментарий | JWT |
| GET | `/communities` | Список сообществ | - |
| GET | `/communities/{id}` | Сообщество по ID | - |
| POST | `/communities` | Создать сообщество | JWT |
| PUT | `/communities/{id}` | Обновить сообщество | JWT |
| DELETE | `/communities/{id}` | Удалить сообщество | JWT |
| DELETE | `/communities/{id}/leave` | Покинуть сообщество | JWT |
| POST | `/communities/{id}/join` | Подать заявку на вступление | JWT |
| GET | `/communities/{id}/requests` | Заявки на вступление | JWT |
| POST | `/communities/{id}/requests/{requestId}/approve` | Одобрить заявку | JWT |
| POST | `/communities/{id}/requests/{requestId}/reject` | Отклонить заявку | JWT |

## Установка и запуск

### 1. Клонирование репозитория

```
git clone https://github.com/TheLonger011/ReLon.git
cd ReLon
```

### 2. Настройка переменных окружения

Скопируйте `.env.example` в `.env` и заполните значения:

```
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_HOST=
DB_PORT=

REDIS_PORT=
REDIS_ADDR=

SERVER_PORT=

JWT_SECRET=

EMAIL_FROM=
EMAIL_PASSWORD=
```

`EMAIL_FROM` и `EMAIL_PASSWORD` - Gmail-адрес и пароль приложения, с которых отправляются коды подтверждения (используется `smtp.gmail.com:587`).

### 3. Запуск через Docker Compose

```
docker-compose up -d
```

Поднимутся сам сервис, PostgreSQL и Redis. Миграции нужно применить отдельно (см. ниже).

### 4. Применение миграций

Требуется установленный [golang-migrate](https://github.com/golang-migrate/migrate):

```
make migrate-up
```

Другие полезные команды:

```
make migrate-down     # откатить последнюю миграцию
make migrate-version   # показать текущую версию
make migrate-create name=имя_миграции
```

### 5. Запуск локально

```
go run cmd/api/main.go
```

Перед локальным запуском убедитесь, что PostgreSQL и Redis доступны по адресам из `.env`, а миграции применены.

---

# ReLon (English)

Backend API for a Reddit-like platform: posts, communities, comments, and voting. Built with Go.

## Features

- Registration and login with email verification code
- JWT authentication
- Posts: create, fetch, search, update, delete
- Voting on posts (cast vote / remove vote)
- Comments on posts: create, fetch, update, delete
- Communities: create, update, delete, list, fetch by ID
- Community join requests: submit, approve, reject by the owner
- Leaving a community
- IP-based rate limiting on sensitive endpoints (register, login, verify, creating posts and comments)

## Technologies

- **Go 1.25**
- **go-chi/chi** - HTTP router
- **PostgreSQL** (jackc/pgx) - primary data store
- **Redis** (go-redis) - email verification code storage
- **golang-jwt/jwt** - JWT tokens
- **golang-migrate** - database migrations
- **net/smtp (Gmail SMTP)** - sending verification emails
- **Docker / Docker Compose** - containerization

## Main Endpoints

| Method | Path | Description | Auth |
| --- | --- | --- | --- |
| POST | `/register` | Register a user | - |
| POST | `/login` | Log in, receive a JWT | - |
| POST | `/verify` | Confirm email with a code | - |
| GET | `/me` | Current user info | JWT |
| GET | `/posts` | List posts | - |
| GET | `/posts/search` | Search posts | - |
| GET | `/posts/{id}` | Get a post by ID | - |
| POST | `/posts` | Create a post | JWT |
| PUT | `/posts/{id}` | Update a post | JWT |
| DELETE | `/posts/{id}` | Delete a post | JWT |
| POST | `/posts/{id}/vote` | Vote on a post | JWT |
| DELETE | `/posts/{id}/vote` | Remove a vote | JWT |
| GET | `/posts/{id}/comments` | Get comments for a post | - |
| POST | `/posts/{id}/comments` | Add a comment | JWT |
| PUT | `/comments/{id}` | Update a comment | JWT |
| DELETE | `/comments/{id}` | Delete a comment | JWT |
| GET | `/communities` | List communities | - |
| GET | `/communities/{id}` | Get a community by ID | - |
| POST | `/communities` | Create a community | JWT |
| PUT | `/communities/{id}` | Update a community | JWT |
| DELETE | `/communities/{id}` | Delete a community | JWT |
| DELETE | `/communities/{id}/leave` | Leave a community | JWT |
| POST | `/communities/{id}/join` | Submit a join request | JWT |
| GET | `/communities/{id}/requests` | List pending join requests | JWT |
| POST | `/communities/{id}/requests/{requestId}/approve` | Approve a join request | JWT |
| POST | `/communities/{id}/requests/{requestId}/reject` | Reject a join request | JWT |

## Installation and Running

### 1. Clone the repository

```
git clone https://github.com/TheLonger011/ReLon.git
cd ReLon
```

### 2. Configure environment variables

Copy `.env.example` to `.env` and fill in the values:

```
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_HOST=
DB_PORT=

REDIS_PORT=
REDIS_ADDR=

SERVER_PORT=

JWT_SECRET=

EMAIL_FROM=
EMAIL_PASSWORD=
```

`EMAIL_FROM` and `EMAIL_PASSWORD` are the Gmail address and app password used to send verification codes (via `smtp.gmail.com:587`).

### 3. Run with Docker Compose

```
docker-compose up -d
```

This starts the app, PostgreSQL, and Redis. Migrations still need to be applied separately (see below).

### 4. Apply migrations

Requires [golang-migrate](https://github.com/golang-migrate/migrate) to be installed:

```
make migrate-up
```

Other useful commands:

```
make migrate-down     # roll back the last migration
make migrate-version   # show the current migration version
make migrate-create name=migration_name
```

### 5. Run locally

```
go run cmd/api/main.go
```

Before running locally, make sure PostgreSQL and Redis are reachable at the addresses from `.env` and that migrations have been applied.
