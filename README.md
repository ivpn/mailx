# Email Service

## API Stack

- Go
- Fiber (API, middleware)
- Gorm (ORM)
- MariaDB (Database)
- Redis (Cache)
- Docker (Containerization)
- Swagger (API Documentation)

## App Stack

- TypeScript
- Vue.js
- Vite (Bundler)
- Tailwind (Styling)
- Docker (Containerization)

## Requirements

- Docker

## Run locally

### Move to api directory
```bash
cd api
```

### Config
```bash
cp .env.sample .env
cp ../app/src/env.sample.json ../app/src/env.json
```

### Build
```bash
docker compose build
```

### Run
```bash
docker compose up
```

### Test
Run tests:  
```bash
go test ./... -v
go vet ./...
gosec ./...
```

Send test email:  
```bash
go run test/send_test_mail.go
```

### API Documentation
Access API docs:  
http://localhost:3000/docs/index.html

Generate API docs:  
```bash
swag init -g cmd/main.go
```
