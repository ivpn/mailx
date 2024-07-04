# Email Service

## API

- Go
- Fiber (API, middleware)
- Gorm (ORM)
- MariaDB (Database)
- Redis (Cache)
- Docker (Containerization)
- Swagger (API Documentation)

## App

- TypeScript
- Vue.js
- Vite (Bundler)
- Tailwind (Styling)
- Docker (Containerization)

## Mailserver

- Postfix - Inbound SMTP

## Requirements

- Docker

## Run locally

### Config
```bash
cp api/.env.sample api/.env
cp app/src/env.sample.json app/src/env.json
cp mailserver/.env.sample mailserver/.env
cp mailserver/postfix-virtual.cf.sample mailserver/docker-data/dms/config/postfix-virtual.cf
cp mailserver/postfix-aliases.cf.sample mailserver/docker-data/dms/config/postfix-aliases.cf
```

### API + App

#### Move to api directory
```bash
cd api
```

#### Build
```bash
docker compose build
```

#### Run
```bash
docker compose up -d
```

App:  
http://localhost:3001  

API:  
http://localhost:3000  

### Mailserver

#### Move to mailserver directory
```bash
cd mailserver
```

#### Run
```bash
docker compose up -d
```

#### Update alias table
```bash
docker exec -it mailserver sh

# Build the db file
postmap /etc/postfix/virtual

# Update the alias table
newaliases

# Restart the postfix service
exit
docker compose up -d
```

### Test
Run API tests:  
```bash
go test ./... -v
go vet ./...
gosec ./...
```

Send test email:  
```bash
docker exec -it mailserver sh
echo "Test email body" | mail -s "Test subject" example.alias12@example.net
```

### API Documentation
API docs:  
http://localhost:3000/docs  

Generate API docs:  
```bash
swag init -g cmd/main.go
```
