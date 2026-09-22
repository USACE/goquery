# Quick Start

### 1. Basic Connection (PostgreSQL)

```go
package main

import (
    "context"
    "log"
    _ "github.com/jackc/pgx/v4/stdlib"
    _ "github.com/usace/goquery/adapters/postgres/v3"
    "github.com/usace/goquery/v3"
)

func main() {
    ctx := context.Background()
    // Create configuration
    config := goquery.RdbmsConfig{
        Dbuser:   "postgres",
        Dbpass:   "password",
        Dbhost:   "localhost",
        Dbport:   "5432",
        Dbname:   "mydb",
        DbDriver: "pgx",
        DbStore:  "pgx",
    }

    // Connect to database
    store, err := goquery.NewRdbmsDataStore(&config)
    if err != nil {
        log.Fatal(err)
    }

    // Use the store with context...
    _ = ctx 
}
```

### 2. Simple Query

```go
ctx := context.Background()

type User struct {
    ID    int    `db:"id"`
    Name  string `db:"name"`
    Email string `db:"email"`
}

// Query into struct slice
var users []User
err := store.Select("SELECT id, name, email FROM users").
    Context(ctx).
    Dest(&users).
    Fetch()
if err != nil {
    log.Fatal(err)
}

for _, user := range users {
    fmt.Printf("%s (%s)\n", user.Name, user.Email)
}
```

### 3. Parameterized Query

```go
ctx := context.Background()

var user User
err := store.Select("SELECT id, name, email FROM users WHERE id = $1").
    Context(ctx).
    Params(42).
    Dest(&user).
    Fetch()
```
