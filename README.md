# video-rental-api

![overview](overview.png)

## Changes
* [Change log in Docs](docs/)

## Install
1. This API uses the database from [pagila](https://github.com/devrimgunduz/pagila)
2. In another directory, clone the pagila db, and run docker-compose up

## Set ENV
```
export DATABASE_URL="postgres://postgres:123456@localhost:5432/postgres"
export PORT=8080
```

## API
### Health Check
```
http://localhost:8080/health
```

### Customer Routes
```
http://localhost:8080/v1/
```

| Method | Path | Description |
| ------ | ---- | ----------- |
| GET | /customers | Get all customers |
| GET | /customers/{id} | Get a customer by ID |
| POST | /customers | Create a new customer|
| DELETE | /customers/{id} | Delete customer by ID |

