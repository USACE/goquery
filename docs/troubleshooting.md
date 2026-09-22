# Troubleshooting

### Common Issues

#### 1. "no rows in result set"

```go
ctx := context.Background()
var user User
err := store.Select("SELECT * FROM users WHERE id = $1").
    Params(999).
    Context(ctx).
    Dest(&user).
    Fetch()
```

**err will be "no rows in result set" if no record exists with id=999**

**Solution:** Check if record exists or handle the error:

```go
if err != nil {
    if strings.Contains(err.Error(), "no rows") {
        return nil, fmt.Errorf("user not found")
    }
    return nil, err
}
```

#### 2. "connection refused"

```
Unable to connect to pgx datastore: connection refused
```

**Checklist:**
- ✅ Is the database running? (`pg_ctl status` or `systemctl status postgresql`)
- ✅ Is the host correct? (use `localhost` or `127.0.0.1` for local)
- ✅ Is the port correct? (default PostgreSQL is `5432`)
- ✅ Firewall blocking the connection?
- ✅ Check `pg_hba.conf` for PostgreSQL access rules

#### 3. "unsupported store type"

```
Unsupported store type: pgxx
```

**Solution:** Check `DbStore` value - must be exactly `"pgx"` or `"sqlx"`:

```go
config := goquery.RdbmsConfig{
    DbStore: "pgx",  // Not "pgxx" or "PGX"
}
```

#### 4. "missing database name"

**Solution:** Ensure all required fields are set:

```go
config := goquery.RdbmsConfig{
    Dbuser:   "postgres",
    Dbpass:   "password",
    Dbhost:   "localhost",
    Dbport:   "5432",
    Dbname:   "mydb",      // Required!
    DbDriver: "pgx",
    DbStore:  "pgx",
}
```

#### 5. "too many open connections"

```
FATAL: sorry, too many clients already
```

**Solutions:**

1. Reduce pool size:
```go
config.PoolMaxConns = 10  // Lower value
```

2. Increase database max connections (PostgreSQL):
```sql
ALTER SYSTEM set max_connections = 200;
SELECT pg_reload_conf();
```

3. Ensure connections are closed:
```go
rows, err := store.Select("...").FetchRows()
defer rows.Close()  // Always close!
```

#### 6. Type Mismatch Errors

```
sql: Scan error: converting NULL to string is unsupported
```

**Solution:** Use pointer types for nullable columns:

```go
type User struct {
    Name  string  `db:"name"`      // NOT NULL column
    Email *string `db:"email"`     // Nullable
}
```

### Debug Logging

Enable SQL logging to see what queries are executed:

```go
err := store.Select("SELECT * FROM users WHERE active = $1").
    Params(true).
    LogSql(true).  // Prints SQL to console
    Dest(&users).
    Fetch()
```

### Performance Profiling

```go
import (
    "time"
)

start := time.Now()
err := store.Select("SELECT * FROM large_table").Dest(&results).Fetch()
duration := time.Since(start)
log.Printf("Query took %v", duration)
```
