# admissions

Service for applicants to the L2SH

## ToC <!-- omit in toc -->

- [Project structure](#project-structure)
- [Build and run](#build-and-run)
  - [Development](#development)
  - [Production](#production)
  - [Ports](#ports)
  - [Secrets](#secrets)
- [Administration](#administration)
  - [PgAdmin](#pgadmin)
- [Testing](#testing)
  - [Run tests](#run-tests)
  - [Code coverage](#code-coverage)

## Project structure

- `cmd/` - application entry points
  - `cmd/admissions/main.go` - main entry point
- `internal/` - internal packages
- `tests/` - tests
- `ui/` - frontend
- `migrations/` - database migrations
- `secrets/` - secrets (ignored by git)
- `config.yml` - configuration file

## Build and run

Before running the application, make sure to create the necessary secrets files (see [Secrets](#secrets)).

### Development

```bash
docker compose --profile dev up --build --watch
```

- `--profile dev` - use `dev` profile that is bound to `database` and `pgadmin` services
- `--watch` - update the container on code changes

### Production

```bash
docker compose up --build
```

### Ports

Default ports:

- `server` - 8888 (set in `config.yml`)
- `pgadmin` - 5050
- `database` - 5432

### Secrets

Secrets are stored in `secrets/` directory.  
`docker-compose.yml` expects the following files:

- `secrets/db_password.txt` - database password
- `secrets/database_host.txt` - database host
- `secrets/database_uri.txt` - database URI (for migrations)

## Administration

### PgAdmin

- URL: http://localhost:5050

Credentials to connect to the development database are in `docker-compose.yml` and `secrets/db_password.txt`.

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
