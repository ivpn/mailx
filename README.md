# Email Service

## MTA package

MTA (message transport agent) package is used for accepting, processing and forwarding emails:

1. Accept incoming emails from other SMTP servers
2. Process emails based on defined rules and filters
3. Transport filtered emails to an outbound SMTP server

## API package

## Run locally

### Config
```bash
cp .env.sample .env
```

### Run
```bash
go run cmd/main.go
```

### Test
Run unit tests:  
```bash
go test ./... --timeout 5s
```

Send test email:  
```bash
go run test/send_test_mail.go
```
