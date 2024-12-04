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

- [Postfix](http://www.postfix.org)
- [Rspamd](https://rspamd.com/)
- [Fail2ban](https://www.fail2ban.org/wiki/index.php/Main_Page)

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
cp mailserver/config/user-patches.sh.sample mailserver/docker-data/dms/config/user-patches.sh
```

> [!IMPORTANT]
> Make sure to set up the required config:
> - api/.env: `DOMAINS`, `SMTP_CLIENT_*`
> - app/src/env.json: `DOMAINS`
> - mailserver/.env: `HOSTNAME`
> - mailserver/docker-data/dms/config/postfix-virtual.cf: `@your-domain.net curl_email`

> [!TIP]
> Run `docker network inspect bridge`. Use "Gateway" value for `API_ALLOW_IP`. Set `API_ALLOW_IP="*"` to disable IP based access control.

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

# Enable DKIM signing, this will output the contents of DKIM TXT DNS record (mail._domainkey.domain.com)
setup config dkim

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

## Restore DB from backup

DB backup is stored locally on the host machine in the `${HOME}/backups` directory.

### Unpack backup
```bash
cd ${HOME}/backups
gpg -o backup.tar.gz -d backup-latest.tar.gz.gpg
tar -xvf backup.tar.gz
```

### Restore DB

```bash
# Stop the containers
docker compose down

# Clone the volume
docker volume create email_db_clone
docker run --rm -v email_db:/from -v email_db_clone:/to alpine sh -c "cp -a /from/. /to/"

# Remove the original volume
docker volume rm email_db

# Recreate the original volume from backup
docker run -d --name restore -v email_db:/email_db alpine
docker cp /unpacked_volume_dir/. restore:/email_db
docker stop restore && docker rm restore

# Start the containers
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
