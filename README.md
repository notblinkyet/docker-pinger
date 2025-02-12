# **Docker Pinger**  

Проект для мониторинга доступности IP-адресов используя **Docker**, **Go**, **PostgreSQL**, **Redis** и **React**.  

## **Структура проекта**  

```
├── backend                # Бэкенд на Go
│   ├── cmd                # Основные файлы запуска
│   ├── internal           # Логика работы сервиса
│   ├── configs            # Конфигурационные файлы
│   ├── migrations         # SQL-миграции
│   ├── Dockerfile.backend # Dockerfile для бэкенда
│   ├── go.mod / go.sum    # Go зависимости
├── pinger                 # Сервис для пинга IP-адресов
│   ├── internal           # Логика работы сервиса
│   ├── configs            # Конфигурационные файлы
│   ├── Dockerfile.pinger  # Dockerfile для пингера
│   ├── go.mod / go.sum    # Go зависимости
├── frontend               # Фронтенд на React
│   ├── src                # Исходный код
│   ├── public             # Статика
│   ├── Dockerfile.frontend # Dockerfile для фронта
│   ├── package.json       # Зависимости
├── docker-compose.main.yaml  # Docker-compose для запуска всех сервисов
└── docker-compose.migration-down.yaml # Для отката миграций
```

## **Технологии**  

- **Backend**: Golang, PostgreSQL, Redis  
- **Frontend**: React (Vite)  
- **Pinger**: Golang, Redis  
- **База данных**: PostgreSQL  
- **Кэширование**: Redis  
- **Контейнеризация**: Docker, Docker Compose  

---

## **API Эндпоинты**  

### **Бэкенд (`/backend`)**  
| Метод | Эндпоинт | Описание |
|--------|----------------------|----------------------------------|
| `POST` | `/backend/containers` | Добавить IP в отслеживаемые |
| `DELETE` | `/backend/containers/:ip` | Удалить IP из отслеживаемых |
| `GET` | `/backend/containers` | Получить список отслеживаемых IP |
| `GET` | `/backend/pings` | Получить все пинги |
| `GET` | `/backend/pings/last` | Получить последние успешные пинги |
| `POST` | `/backend/pings` | Отправить результаты пинга |

### **Связь сервисов**  
- **Frontend** → **Backend**  
- **Backend** → **Pinger**  
- **Pinger** → **Redis**  
- **Pinger** → **Backend** (для отправки пингов)  

---

## **Структура базы данных**  

```sql
CREATE TABLE IF NOT EXISTS containers(
    id SERIAL PRIMARY KEY,
    ip VARCHAR(15) UNIQUE,
    is_tracked BOOLEAN
);

CREATE TABLE IF NOT EXISTS pings(
    id SERIAL PRIMARY KEY,
    container_id INT,
    latency BIGINT NOT NULL,
    last_success_at TIMESTAMP DEFAULT NULL,
    ping_at TIMESTAMP,
    FOREIGN KEY (container_id) REFERENCES containers (id)
);

CREATE TABLE IF NOT EXISTS last_pings(
    id INT PRIMARY KEY,
    container_id INT UNIQUE,
    latency BIGINT NOT NULL,
    last_success_at TIMESTAMP DEFAULT NULL,
    ping_at TIMESTAMP,
    FOREIGN KEY (container_id) REFERENCES containers (id) ON DELETE CASCADE
);
```

---

## **Запуск проекта**  

### **Локально (без Docker)**  

```sh
# 1. Клонировать репозиторий
git clone https://github.com/notblinkyet/docker-pinger
cd docker-pinger

# 2. Настроить конфигурационные файлы
# backend/configs/local.yaml
# pinger/configs/local.yaml
# frontend/.env

# 3. Установить переменные окружения
cd backend && export POSTGRES_PASS= && export BACKEND_CONFIG_PATH= && cd ..
cd pinger && export REDIS_PASS= && export PINGER_CONFIG_PATH= && cd ..

# 4. Запустить миграции
go run backend/cmd/migrator/main.go

# 5. Запустить сервисы
go run backend/cmd/backend/main.go
go run pinger/cmd/main.go

# 6. Запустить фронтенд
cd frontend && npm install && npm run dev
```

Открыть в браузере: **[http://localhost:5173](http://localhost:5173)**  

---

### **Запуск в Docker**  

```sh
# 1. Клонировать репозиторий
git clone https://github.com/notblinkyet/docker-pinger
cd docker-pinger

# 2. Настроить конфигурационные файлы
# backend/configs/remote.yaml
# pinger/configs/remote.yaml
# frontend/.env

# 3. Запустить через Docker Compose
docker compose -f docker-compose.main.yaml up --build
```

Открыть в браузере: **[http://localhost:5173](http://localhost:5173)**  

---

## **TODO / Возможные улучшения**  

- Добавить **авторизацию**  
- Улучшить **логирование и мониторинг**  
- Добавить **очередь сообщений**
- Добавить **балансировщик**
- Оптимизировать запросы к **Redis**  

---

### **Автор**  

[notblinkyet](https://github.com/notblinkyet)  
