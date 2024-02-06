# Email Service

## Run locally

### Config
```bash
cp ./config/.env.sample ./config/.env
```

### Build
```bash
docker build -t email-service .
```

### Run
```bash
docker run -it --publish 8025:8025 email-service
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
