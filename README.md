# admissions

Сервис для поступающих в Лицей "Вторая школа"

## Структура проекта

- `server.go` - точка входа
- `ui/` - фронтенд
- `scripts/` - miscellaneous scripts
- `migrations/` - database migrations
- `secrets/` - secrets (ignored by git)

## Сборка и запуск

### Profiles

`docker-compose.yml` has `dev` profile that is bound to `database` and `pgadmin` services. Therefore in development run:

```bash
docker compose --profile dev up --build
```

For production without mock database and pgAdmin:

```bash
docker compose up --build
```

### Secrets

Secrets are stored in `secrets/` directory.  
`docker-compose.yml` expects the following files:

- `secrets/database_host.txt` - database host
- `secrets/database_uri.txt` - database URI (for migrations)

## Administration

### PgAdmin

- URL: http://localhost:5050

To connect to the database, use credentials from `docker-compose.yml`.  
How to find out the host name of the database:

```bash
docker inspect database | grep IPAddress
```

Or use the following command:

```bash
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' database
```

It's a bit more complicated, but returns only the IP address. It's convenient if you want to copy it into secrets for example.

```bash
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' database > secrets/database_host.txt
```
