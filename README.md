# Скачать образ
docker pull nilidushka/urlsh:latest

# Запустить контейнер с in-memory хранилищем
docker run -d -p 8080:8080 --name url-shortener nilidushka/urlsh:latest -storage=memory

# Или с PostgreSQL (предварительно запустите PostgreSQL)
## Пример запуска у себя PostgreSQL
docker run -d --name postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=shortener -p 5432:5432 postgres:15-alpine
## Подключение приложения
docker run -d -p 8080:8080 --name url-shortener nilidushka/urlsh:latest -storage=postgres -pg-conn="postgres://postgres:password@host.docker.internal:5432/shortener?sslmode=disable"

#В случае ошибки с PostgreSQL
## Создать сеть
docker network create shortener-network
## PostgreSQL
docker run -d --name postgres --network shortener-network -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=shortener -p 5432:5432 postgres:15-alpine
## Подключение приложения
docker run -d -p 8080:8080 --name url-shortener --network shortener-network nilidushka/urlsh:latest -storage=postgres -pg-conn="postgres://postgres:password@postgres:5432/shortener?sslmode=disable"
