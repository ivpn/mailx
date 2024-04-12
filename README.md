# Email Service

## API Stack

- Go
- Fiber (API, middleware)
- Gorm (ORM)
- MariaDB (Database)
- Redis (Cache)
- Docker (Containerization)

## App Stack

- TypeScript
- Vue.js
- Vite (Bundler)
- Tailwind (Styling)
- Docker (Containerization)

## Run locally

### Move to api directory
```bash
cd api
```

### Config
```bash
cp .env.sample .env
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
