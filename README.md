# L0WBTECHSCHOOL — Order Processing Service

Проект: учебный сервис на Go для обработки заказов.  
Он читает входящие заказы (Kafka), сохраняет их в PostgreSQL, кэширует, и позволяет получать данные через HTTP.

---

##  Структура проекта

├── cmd/
├── internal/
│ ├── cache/
│ ├── consumer/
│ ├── storage/
│ ├── api/
│ └── config/
├── migrations/
├── docker-compose.yml
├── init.sql
├── index.html
├── go.mod / go.sum
└── README.md

---

##  Технологии

- **Go** — язык разработки  
- **Apache Kafka** — поставка заказов  
- **PostgreSQL** — хранилище заказов  
- **In-memory Cache** — быстрый доступ к данным  
- **HTTP API & Web UI** — выдача заказа по ID  
- **Docker Compose** — поднимает Kafka + Postgres  
- **SQL миграции** — инициализация БД (`init.sql`, `migrations/`)

---

##  Установка и запуск

# Клонировать репозиторий
git clone https://github.com/Iodory/L0WBTECHSCHOOL.git
cd L0WBTECHSCHOOL

# Запустить инфраструктуру (Kafka, Postgres)
docker-compose up -d

# Инициализировать базу данных (если нужно)
psql -f init.sql -U postgres

# Запустить сервис
go run ./cmd

# Пример запроса к API:
curl http://localhost:8080/order/123
