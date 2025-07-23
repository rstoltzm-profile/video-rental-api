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

## Optional port change if needed on the pg db 6543:5432
```
export DATABASE_URL="postgres://postgres:123456@localhost:6543"
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

### Rental Routes
```
http://localhost:8080/v1/
```

| Method | Path | Description |
| ------ | ---- | ----------- |
| GET | /rentals | Get all rentals |
| GET | /rentals?late=true | Get all late rentals |
| GET | /rentals?customer_id={id} | Get rentals for customer |
| GET | /rentals?customer_id={id}&late=true | Get late rentals for customer |
| POST | /rentals | Rents inventory with payload --date rental.json |
| POST | /rentals/{id}/return | Returns a rental for {id} |

#### Rental payload
```json
{
  "inventory_id": 709,
  "customer_id": 397,
  "staff_id": 1
}
```

### Inventory Routes
```
http://localhost:8080/v1/
```

| Method | Path | Description |
| ------ | ---- | ----------- |
| GET | /inventory | Get all inventory |
| GET | /inventory?store_id=1 | Get all inventory by store_id |
| GET | /inventory/available?film_id=1&store_id=2 | Get all available inventory by store_id and film_id |

### Store Routes
```
http://localhost:8080/v1/
```

| Method | Path | Description |
| ------ | ---- | ----------- |
| GET | /stores/{id}/inventory/summary | Get Count of inventory by store  |

### Films Routes
```
http://localhost:8080/v1/
```

| Method | Path                               | Description                                   |
|--------|------------------------------------|-----------------------------------------------|
| GET    | /stores/{id}/inventory/summary     | Get count of inventory by store               |
| GET    | /films                             | Get all films                                 |
| GET    | /films/{id}                        | Get a single film by ID                       |
| GET    | /films/search?title={query}        | Search for films by title                     |
| GET    | /films/{id}/with-actors-categories | Get film details with actors and categories   |

