# Migration Guide

## Overview of Breaking Changes

Version 3.0 introduces breaking changes that require code updates. The main changes are:

1. **Module path includes `/v3`**
2. **Adapters must be explicitly imported**
3. **`OnInit` deprecated in favor of `OnConnect`**

---

## Migration from database/sql or GORM

### From database/sql

**Before (database/sql)**

```go
// Before (database/sql)
rows, err := db.Query("SELECT id, name FROM users WHERE active = $1", true)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

var users []User
for rows.Next() {
    var user User
    err := rows.Scan(&user.ID, &user.Name)
    if err != nil {
        log.Fatal(err)
    }
    users = append(users, user)
}
```

**After (goquery)**

```go
// After (goquery)
var users []User
err := store.Select("SELECT id, name FROM users WHERE active = $1").
    Params(true).
    Dest(&users).
    Fetch()
```

### From GORM

**Before (GORM)**

```go
// Before (GORM)
var users []User
db.Where("active = ?", true).Find(&users)
```

**After (goquery)**

```go
// After (goquery)
var users []User
store.Select("SELECT * FROM users WHERE active = $1").
    Params(true).
    Dest(&users).
    Fetch()
```

---

## v3 Migration Guide (Upgrade from v1)

### 1. Update Import Paths

**Before (v1):**

```go
import "github.com/usace/goquery"
```

**After (v3):**

```go
import "github.com/usace/goquery/v3"
```

### 2. Import Database Adapters

In v3, you must explicitly import the adapter for your database.

**Before (v1):**

```go
import (
    _ "github.com/jackc/pgx/v4/stdlib"
    "github.com/usace/goquery"
)
```

**After (v3):**

```go
import (
    _ "github.com/jackc/pgx/v4/stdlib"
    _ "github.com/usace/goquery/adapters/postgres/v3"  // New: adapter import required
    "github.com/usace/goquery/v3"
)
```

**Adapter import paths:**

| Database | Adapter Import |
|----------|----------------|
| PostgreSQL | `_ "github.com/usace/goquery/adapters/postgres/v3"` |
| DuckDB | `_ "github.com/usace/goquery/adapters/duckdb/v3"` |
| SQLite | `_ "github.com/usace/goquery/adapters/sqlite/v3"` |
| Oracle | `_ "github.com/usace/goquery/adapters/oracle/v3"` |

### 3. Update go.mod

Run these commands:

```bash
# Remove old version
go get github.com/usace/goquery@none

# Add v3
go get github.com/usace/goquery/v3

# Add adapter(s) you need
go get github.com/usace/goquery/adapters/postgres/v3
go get github.com/usace/goquery/adapters/duckdb/v3
go get github.com/usace/goquery/adapters/sqlite/v3
go get github.com/usace/goquery/adapters/oracle/v3

# Clean up
go mod tidy
```

### 4. Replace OnInit with OnConnect (Oracle users)

If you were using `OnInit` for Oracle initialization:

**Before (v1):**

```go
config := goquery.RdbmsConfig{
    DbDriver:    "godror",
    DbStore:     "sqlx",
    ExternalLib: "/usr/lib/oracle/instantclient",
    OnInit:      "ALTER SESSION SET NLS_DATE_FORMAT='YYYY-MM-DD HH24:MI:SS'",
}
```

**After (v3):**

```go
config := goquery.RdbmsConfig{
    DbDriver:    "godror",
    DbStore:     "sqlx",
    ExternalLib: "/usr/lib/oracle/instantclient",
    OnConnect: func(db goquery.DataStore) error {
        return db.Exec(goquery.NoTx, "ALTER SESSION SET NLS_DATE_FORMAT='YYYY-MM-DD HH24:MI:SS'")
    },
}
```

**Note:** `OnInit` still works for backward compatibility but is deprecated.

### 5. Update SQLite Driver Names (if applicable)

If you're using SQLite and want to use the native Go driver:

**Before (v1):**

```go
config := goquery.RdbmsConfig{
    DbDriver: "sqlite3",  // CGO driver
    DbStore:  "sqlx",
    Dbname:   "./mydb.db",
}
```

**After (v3) - Native Go (recommended):**

```go
import _ "modernc.org/sqlite"

config := goquery.RdbmsConfig{
    DbDriver: "sqlite",   // Native Go driver (no CGO)
    DbStore:  "sqlx",
    Dbname:   "./mydb.db",
}
```

**Or keep using CGO:**

```go
import _ "github.com/mattn/go-sqlite3"

config := goquery.RdbmsConfig{
    DbDriver: "sqlite3",     // CGO driver (still supported)
    DbStore:  "sqlx",
    Dbname:   "./mydb.db",
}
```
