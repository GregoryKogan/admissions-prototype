# admissions

Сервис для поступающих в Лицей "Вторая школа"

## Структура проекта

- `server.go` - точка входа
- `ui/` - фронтенд

## Сборка и запуск

### С Docker (рекомендуется)

```bash
docker-compose up --build
```

### Без Docker

```bash
# Сборка фронтенда
cd ui
yarn install
yarn build
cd ..
# Сборка и запуск сервера
go mod download
go build -o server .
./server
```
