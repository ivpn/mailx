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

- [Docker Mailserver](https://github.com/docker-mailserver/docker-mailserver)

## Installation

### Requirements

- Docker

> [!IMPORTANT]
> Docker Mailserver officially supports Linux. If you want to run it on macOS, please read [this](https://github.com/docker-mailserver/docker-mailserver/issues/3648).

### Config
```bash
cp api/.env.sample api/.env
cp app/src/env.sample.json app/src/env.json
cp mailserver/.env.sample mailserver/.env
mkdir -p mailserver/docker-data/dms/config
cp mailserver/config/postfix-main.cf.sample mailserver/docker-data/dms/config/postfix-main.cf
cp mailserver/config/postfix-virtual.cf.sample mailserver/docker-data/dms/config/postfix-virtual.cf
cp mailserver/config/postfix-aliases.cf.sample mailserver/docker-data/dms/config/postfix-aliases.cf
```

> [!IMPORTANT]
> Make sure to set up the required config:
> - api/.env: `DOMAINS`, `SMTP_CLIENT_*`
> - app/src/env.json: `DOMAINS`
> - mailserver/.env: `HOSTNAME`
> - mailserver/docker-data/dms/config/postfix-virtual.cf: `@your-domain.net curl_email`

> [!TIP]
> For local testing, you can use [MailHog](https://github.com/mailhog/MailHog) or [MailTrap](https://mailtrap.io/email-sandbox/) as outbound SMTP client (`SMTP_CLIENT_*`).

### API + App

#### Move to api directory
```bash
cd api
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

#### Setup Postfix
```bash
docker exec -it mailserver sh

# Build the db file
postmap /etc/postfix/virtual

# Update the alias table
newaliases

# Restart Postfix
supervisorctl restart postfix

# Show logs
setup debug show-mail-logs
```

#### Update Mailserver
```bash
docker compose pull
docker compose down
docker compose up -d
```

## Test
Run API tests:  
```bash
go test ./... -v
go vet ./...
gosec ./...
```

Send test email:  
```bash
docker exec -it mailserver sh
echo "Test email body" | mail -s "Test subject" example.alias@example.net
```

## API Documentation
API docs:  
http://localhost:3000/docs  

Generate API docs:  
```bash
cd api
swag init -g cmd/main.go
```

> [!TIP]
> With [Task](https://github.com/go-task/task), run `task docs` to generate API documentation.
