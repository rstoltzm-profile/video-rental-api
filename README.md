# video-rental-api

## Set ENV
```
export DATABASE_URL="postgres://postgres:123456@localhost:5432/postgres"
export PORT=8080
```

## Health Check
```bash
curl -i http://localhost:8080/health

```
```text
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 13 Jul 2025 22:10:35 GMT
Content-Length: 16

{"status":"ok"}
```

## Tests
### Internal DB
```
go test ./internal/db
```