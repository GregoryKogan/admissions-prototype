# admissions

Service for applicants to the L2SH

## Project structure

- `server.go` - entry point
- `handlers` - api request handlers
- `ui/` - frontend

## Build and run

Server listens on `localhost:8080` by default.

### Docker (recommended)

```bash
docker-compose up --build
```

### Manual

```bash
# Build frontend
cd ui
yarn install
yarn build
cd ..
# Build and run server
go mod download
go build -o server .
./server
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
