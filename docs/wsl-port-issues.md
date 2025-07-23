# Had an issue with connecting to localhost:5432 on different system

## Workaround adjust the port for docker compose, for some reason 5432:5432 can conflict for certain systems
```
    ports:
      - 6543:5432
```

```
export DATABASE_URL="postgres://postgres:123456@localhost:6543"
```