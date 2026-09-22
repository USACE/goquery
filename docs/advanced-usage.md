# Advanced Usage

### Custom Row Scanning

```go
ctx := context.Background()
rows, err := store.Select("SELECT id, name, email FROM users").FetchRows()
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var id int
    var name, email string
    
    err := rows.Scan(&id, &name, &email)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%d: %s <%s>\n", id, name, email)
}
```

### Column Metadata

```go
ctx := context.Background()
rows, err := store.Select("SELECT * FROM users").FetchRows()
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

// Get column names
columns, err := rows.Columns()
if err != nil {
    log.Fatal(err)
}
fmt.Println("Columns:", columns)

// Get column types
types, err := rows.ColumnTypes()
if err != nil {
    log.Fatal(err)
}

for i, col := range columns {
    fmt.Printf("Column %s has type %v\n", col, types[i])
}
```

### Working with NULL Values

Use pointer types for nullable columns:

```go
ctx := context.Background()
type User struct {
    ID          int32      `db:"id"`
    Name        string     `db:"name"`
    Email       *string    `db:"email"`        // Nullable
    PhoneNumber *string    `db:"phone_number"` // Nullable
    LastLogin   *time.Time `db:"last_login"`   // Nullable
}

// Inserting with NULL values
email := "user@example.com"
user := User{
    Name:        "John Doe",
    Email:       &email,      // Has value
    PhoneNumber: nil,         // NULL
    LastLogin:   nil,         // NULL
}

err := store.Insert(&usersDS).Records(&user).Context(ctx).Execute()
```

### Multiple Databases

```go
// Connect to multiple databases
ctx := context.Background()
pgConfig := goquery.RdbmsConfig{
    Dbhost:   "postgres-host",
    DbDriver: "pgx",
    DbStore:  "pgx",
    // ...
}
pgStore, err := goquery.NewRdbmsDataStore(&pgConfig)

sqliteConfig := goquery.RdbmsConfig{
    Dbname:   "/path/to/local.db",
    DbDriver: "sqlite",
    DbStore:  "sqlx",
}
sqliteStore, err := goquery.NewRdbmsDataStore(&sqliteConfig)

// Use them independently
var pgUsers []User
pgStore.Select("SELECT * FROM users").Context(ctx).Dest(&pgUsers).Fetch()

var sqliteCache []CacheEntry
sqliteStore.Select("SELECT * FROM cache").Context(ctx).Dest(&sqliteCache).Fetch()
```

### Schema Support

```go
// TableDataSet with schema
ctx := context.Background()
var productsDS = goquery.TableDataSet{
    Name:   "products",
    Schema: "sales",  // Queries will use "sales.products"
    Statements: goquery.Statements{
        "get-all": "SELECT * FROM sales.products",
    },
    TableFields: Product{},
}
```

### Generated SQL Inspection

```go
// See what SQL will be generated
ctx := context.Background()
stmt, err := store.(*goquery.RdbmsDataStore).RdbmsDb().InsertStmt(&usersDS)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated INSERT:", stmt)
// Output: INSERT INTO users (id, name, email) VALUES (nextval('users_id_seq'), $1, $2)
```
