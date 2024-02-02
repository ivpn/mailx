# Email Service

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
