# video-rental-api

## Set ENV
```
export DATABASE_URL="postgres://postgres:123456@localhost:5432/postgres"
export PORT=8080
```

## Health Check
```
curl localhost:8080/health
{"status":"ok"}
```

## Tests
### Internal DB
```
go test ./internal/db
```