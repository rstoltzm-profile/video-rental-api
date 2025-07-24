# 13-db-pool

## moved to db pooling
http://localhost:8080/health/pool
```
{
  "connections": {
    "acquired": 0,
    "idle": 3,
    "total": 3
  },
  "database": "ok",
  "status": "ok"
}
```

## added basic db migrations
```
2025/07/23 17:24:27 Connected to database on attempt 1
2025/07/23 17:24:27 Migration pagila-initial already applied.
2025/07/23 17:24:27 Migration 2025-07-23-placeholder-table already applied.
2025/07/23 17:24:27 Migration 2025-07-23-drop-placeholder-table already applied.
2025/07/23 17:24:27 Database schema version: 2025-07-23-drop-placeholder-table (applied at 2025-07-24T00:18:48Z)
2025/07/23 17:24:27 Starting server on port 8080...
```