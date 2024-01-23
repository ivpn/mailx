# Email Service

## Inbound SMPT service

Inbound SMPT service is an incoming message handler. It should be used to route/forward or drop mail messages as defined by filters and rules.

### Config
```bash
cp services/inbound/.env.sample services/inbound/.env
```

### Run
```bash
go run services/inbound/main.go
```

### Test
Run unit tests:  
```bash
go test ./services/inbound/... --timeout 5s
```

Send test email:  
```bash
go run services/inbound/test/send_test_mail.go
```

## API service
