# admissions

Service for applicants to the L2SH

## Project structure

- `server.go` - entry point
- `handlers` - api request handlers
- `ui/` - frontend
- `scripts/` - miscellaneous scripts
- `migrations/` - database migrations
- `secrets/` - secrets (ignored by git)

## Build and run

Server listens on `localhost:8888` by default.

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
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' database
```

## Testing

### Run tests

```bash
go test -v ./...
```

### Code coverage

```bash
go test -v -coverprofile=coverage.out ./...
```

`coverage.out` file will be generated in the project root directory.  
To view the coverage report, run:

```bash
go tool cover -html=coverage.out
```
