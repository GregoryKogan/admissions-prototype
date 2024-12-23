# admissions

This is an early prototype of an admissions system for the L2SH (Лицей "Вторая школа") school.  
The system is built with Go (Echo framework) and Vue.js. It uses Redis for caching and PostgreSQL for the database. A React Admin panel is provided for administration and API for it is built with PostgREST.

## ToC <!-- omit in toc -->

- [🗂️ Project structure](#️-project-structure)
- [🚀 Build and run](#-build-and-run)
  - [🌱 Development](#-development)
  - [🛠️ Production](#️-production)
  - [🔌 Ports](#-ports)
  - [Secrets](#secrets)
- [🛎️ Administration](#️-administration)
  - [📈 Logging](#-logging)
  - [🌐 PgAdmin](#-pgadmin)
- [🎨 Admin Panel](#-admin-panel)
- [🧪 Testing](#-testing)
  - [✅ Run tests](#-run-tests)
  - [Code coverage](#code-coverage)

## 🗂️ Project structure

- `cmd/` - application entry points
- `internal/` - internal packages
- `ui/` - frontend
- `tests/` - tests
- `admin-panel/` - React Admin panel
- `secrets/` - secrets (ignored by git)
- `config.yml` - configuration

## 🚀 Build and run

Before running the application, make sure to set the required environment variables [secrets](#secrets).

### 🌱 Development

```bash
docker compose up --build --watch
```

- `--watch` - update the container on code changes

### 🛠️ Production

```bash
docker compose up --build
```

### 🔌 Ports

Default ports:

- `server` - 8888 (set in `config.yml`)
- `admin-panel` - 4444
- `ui` - 3000 (Vite development server)
- `pgadmin` - 5050
- `database` - 5432

### Secrets

Secrets are loaded from environment variables.  
Set the following variables before running the application:

- DB_PASSWORD - password for the database
- JWT_KEY - secret key for JWT signing
- MAIL_API_KEY - NotiSend API key
- ADMIN_PASSWORD - password for the default admin user

## 🛎️ Administration

### 📈 Logging

The application logs to `stdout`, which can be viewed with `docker logs` command.

```bash
docker logs admissions
```

### 🌐 PgAdmin

- URL: http://localhost:5050

Credentials to connect to the development database are in `docker-compose.yml` and `DB_PASSWORD` secret.

## 🎨 Admin Panel

The admin panel is a separate frontend built with PostgREST and React Admin. A Docker service is provided
in the docker-compose.yml under the "admin-panel" service. It can be accessed at http://localhost:4444
once the container is running.

## 🧪 Testing

### ✅ Run tests

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
