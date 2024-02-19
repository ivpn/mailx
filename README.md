# Email Service

## Run locally

### Config
```bash
cp ./config/.env.sample ./config/.env
```

### Build
```bash
docker build -t email-app .
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
