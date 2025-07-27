# video-rental-api

![overview](overview.png)

## Frontend
- [https://github.com/rstoltzm-profile/video-rental-portal-django](https://github.com/rstoltzm-profile/video-rental-portal-django)

## Changes
* [Change log in Docs](docs/)

## Landing Page - localhost:8080
![landing](docs/landing_page.png)

## Install
1. This API uses the database from [pagila](https://github.com/devrimgunduz/pagila)
2. In another directory, clone the pagila db, and run docker-compose up

## Set ENV
```
export DATABASE_URL="postgres://postgres:123456@localhost:5432/postgres"
export PORT=8080
export API_KEY="secure-dev-key-123"
```

## Swagger Setup
```
go install github.com/swaggo/swag/cmd/swag@latest
export PATH=$PATH:$(go env GOPATH)/bin
make docs-swagger
```

### Optional port change if needed on the pg db 6543:5432
```
export DATABASE_URL="postgres://postgres:123456@localhost:6543"
```

## API
### Health Check
```
http://localhost:8080/
http://localhost:8080/health
```

### Swagger
```
http://localhost:8080/swagger
```
![swagger1](docs/swagger/swagger-1.png)
![swagger2](docs/swagger/swagger-2.png)
![swagger3](docs/swagger/swagger-3.png)