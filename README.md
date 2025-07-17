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

### Prerequisites

- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- Docker Compose: [Install Docker Compose](https://docs.docker.com/compose/install/)

> [!IMPORTANT]
> Docker Mailserver officially supports Linux. If you want to run it on macOS, please read [this](https://github.com/docker-mailserver/docker-mailserver/issues/3648).

### Config
```bash
cp api/.env.sample api/.env
cp app/.env.sample app/.env
cp mailserver/.env.sample mailserver/.env
mkdir -p mailserver/docker-data/dms/config/rspamd/override.d
cp mailserver/config/postfix-main.cf.sample mailserver/docker-data/dms/config/postfix-main.cf
cp mailserver/config/postfix-virtual.cf.sample mailserver/docker-data/dms/config/postfix-virtual.cf
cp mailserver/config/postfix-aliases.cf.sample mailserver/docker-data/dms/config/postfix-aliases.cf
cp mailserver/config/user-patches.sh.sample mailserver/docker-data/dms/config/user-patches.sh
cp mailserver/config/rspamd/override.d/milter_headers.conf.sample mailserver/docker-data/dms/config/rspamd/override.d/milter_headers.conf
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
# localhost
docker compose up -d

# Staging|Production
docker compose -f compose.deploy.yml up -d
```

#### Setup Postfix
```bash
docker exec -it mailserver sh

# Build the db file
postmap /etc/postfix/virtual

# Update the alias table
newaliases

# Enable DKIM signing, this will output the contents of DKIM TXT DNS record (mail._domainkey.domain.com)
setup config dkim selector mail

# Restart Postfix
supervisorctl restart postfix

# Show logs
setup debug show-mail-logs
```

#### Setup DKIM Signing

Update dkim_signing.conf:
```bash
nano docker-data/dms/config/rspamd/override.d/dkim_signing.conf
```

dkim_signing.conf:
```conf
# documentation: https://rspamd.com/doc/modules/dkim_signing.html

enabled = true;

sign_authenticated = true;
sign_local = true;
try_fallback = false;

use_domain = "header";
use_domain_sign_local = "header";
use_redis = false; # don't change unless Redis also provides the DKIM keys
use_esld = true;
allow_username_mismatch = true;

check_pubkey = false; # you want to use this in the beginning

domain {
    domain.com {
        path = "/tmp/docker-mailserver/rspamd/dkim/rsa-2048-mail-domain.com.private.txt";
        selector = "mail";
    }
}
```

Restart Mailserver:
```bash
docker compose down
docker compose up -d
```

Test DKIM signing with https://dkimvalidator.com:
```bash
echo "Test email" | mail -s "Test email" wyygMXeSfnzl5l@dkimvalidator.com
```

#### Update Mailserver
```bash
docker compose pull
docker compose down
docker compose up -d
```


## Add New Domain

### Mailserver

#### Update postfix-virtual.cf:
```bash
nano docker-data/dms/config/postfix-virtual.cf
@domain.com curl_email
```

#### Setup Postfix:
```bash
docker exec -it mailserver sh
postmap /etc/postfix/virtual
newaliases
setup config dkim selector mail domain domain.com
supervisorctl restart postfix
```

#### Update dkim_signing.conf:
```bash
nano docker-data/dms/config/rspamd/override.d/dkim_signing.conf
domain.com {
    path = "/tmp/docker-mailserver/rspamd/dkim/rsa-2048-mail-domain.com.private.txt";
    selector = "mail";
}
```

#### Restart mailserver:
```bash
docker compose down
docker compose -f compose.deploy.yml up -d
```

### DNS Records

#### MX Records:
```
domain.com. . 14400 IN MX 10 mail.domain.com.
mail.domain.com. . 14400 IN MX 10 MAIL_SERVER_IPV4
```

#### SPF Records:
```
domain.com. 3600 IN TXT "v=spf1 ip4:MAIL_SERVER_IPV4 -all"
mail.domain.com. 3600 IN TXT "v=spf1 ip4:MAIL_SERVER_IPV4 -all"
```

#### DMARC (TXT record):
```
_dmarc.mail.domain.com. 3600 IN TXT v=DMARC1; p=quarantine
````

DKIM (TXT record):
```
mail._domainkey.domain.com. 3600 IN TXT v=DKIM1;k=rsa;p=DKIM_PUBLIC_KEY
```

### API

#### Update .env:
```bash
nano .env
DOMAINS=domain1.com,domain2.com
```

#### Restart api:
```bash
docker compose down
docker compose up -d
```

### App

#### Update GitHub Actions env variables:
```bash
STAGING_VITE_DOMAINS=domain1.com,domain2.com
PROD_VITE_DOMAINS=domain1.com,domain2.com
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
