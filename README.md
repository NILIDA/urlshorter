# Скачать образ
docker pull nilidushka/urlsh:latest

# Запустить контейнер с in-memory хранилищем
docker run -d -p 8080:8080 --name url-shortener nilidushka/urlsh:latest -storage=memory

# Или с PostgreSQL (предварительно запустите PostgreSQL)
docker run -d -p 8080:8080 --name url-shortener nilidushka/urlsh:latest -storage=postgres -pg-conn="postgres://postgres:password@host.docker.internal:5432/shortener?sslmode=disable"
