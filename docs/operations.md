# DataStore Operations

### SELECT Queries

#### Basic Select

```go
// Query all rows
ctx := context.Background()
var users []User
err := store.Select("SELECT * FROM users").
    Context(ctx).
    Dest(&users).
    Fetch()

// Query single row
ctx := context.Background()
var user User
err := store.Select("SELECT * FROM users WHERE id = $1").
    Context(ctx).
    Params(123).
    Dest(&user).
    Fetch()
```

#### With DataSet

```go
// With DataSet
ctx := context.Background()
var users []User
err := store.Select().
    DataSet(&usersDS).
    StatementKey("get-all").
    Context(ctx).
    Dest(&users).
    Fetch()
```

#### Dynamic Query Building

```go
// Apply() for internal query building (column names, table names)
ctx := context.Background()
var user User
err := store.Select("SELECT * FROM %s WHERE %s = $1").
    Apply("users", "id").
    Context(ctx).
    Params(123).
    Dest(&user).
    Fetch()
```

**⚠️ SECURITY WARNING:**  
`Apply()` uses `fmt.Sprintf()` for SQL templating. **NEVER** pass user input to `Apply()`. Use `Params()` for all user data.

```go
// ✅ CORRECT - Apply() for internal values, Params() for user input
tableName := "users"  // Internal constant
userID := getUserInput()
err := store.Select("SELECT * FROM %s WHERE id = $1").
    Apply(tableName).
    Params(userID).
    Dest(&user).
    Fetch()

// ❌ WRONG - SQL Injection vulnerability!
userTable := getUserInput()
err := store.Select("SELECT * FROM %s").
    Apply(userTable).  // NEVER DO THIS!
    Dest(&results).
    Fetch()
```

#### Query Suffixes

```go
// Add WHERE, ORDER BY, LIMIT dynamically
ctx := context.Background()
var users []User
err := store.Select("SELECT * FROM users").
    Suffix("WHERE active = true ORDER BY created_at DESC LIMIT 10").
    Context(ctx).
    Dest(&users).
    Fetch()

// With parameters
ctx := context.Background()
err := store.Select("SELECT * FROM users").
    Suffix("WHERE active = $1 AND role = $2").
    Params(true, "admin").
    Dest(&users).
    Fetch()
```

#### Logging SQL

```go
err := store.Select("SELECT * FROM users WHERE id = $1").
    Params(123).
    LogSql(true).  // Logs the final SQL statement
    Dest(&user).
    Fetch()
```

### Row Iteration

#### Manual Iteration

```go
rows, err := store.Select("SELECT * FROM users").FetchRows()
if err != nil {
    log.Fatal(err)
}
defer rows.Close()  // Always close!

for rows.Next() {
    var user User
    if err := rows.ScanStruct(&user); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User: %s <%s>\n", user.Username, user.Email)
}
```

#### ForEachRow (Recommended)

```go
// Automatic resource management
var user User
err := store.Select("SELECT * FROM users").
    ForEachRow(func(row goquery.Rows) error {
        if err := row.ScanStruct(&user); err != nil {
            return err
        }
        fmt.Printf("User: %s <%s>\n", user.Username, user.Email)
        return nil
    }).
    Fetch()
```

### INSERT Operations

#### Single Insert

```go
// Auto-generates INSERT statement from struct tags
ctx := context.Background()
user := User{
    Username: "john_doe",
    Email:    "john@example.com",
}

err := store.Insert(&usersDS).
    Records(&user).
    Context(ctx).
    Execute()
```

#### Bulk Insert

```go
users := []User{
    {Username: "alice", Email: "alice@example.com"},
    {Username: "bob", Email: "bob@example.com"},
    {Username: "charlie", Email: "charlie@example.com"},
}

ctx := context.Background()
err := store.Insert(&usersDS).
    Records(&users).
    Context(ctx).
    Execute()
```

#### Batch Insert (High Performance)

Batch inserts are more efficient for large datasets (pgx only):

```go
// Batch insert with error handling
ctx := context.Background()
users := make([]User, 10000)
// ... populate users ...

err := store.Insert(&usersDS).
    Records(&users).
    Context(ctx).
    Batch(true).
    BatchSize(500).  // 500 records per batch
    Execute()

if err != nil {
    log.Fatalf("Batch insert failed: %v", err)
}
```

**Batch vs Bulk:**
- **Bulk Insert** - Uses multiple transactions, one per record
- **Batch Insert** - Uses pgx batching, single network roundtrip per batch (much faster)

#### INSERT with RETURNING

```go
// Get generated ID
ctx := context.Background()
var newID int
err := store.Select().
    DataSet(&usersDS).
    StatementKey("insert-with-return").
    Params("john_doe", "john@example.com").
    Context(ctx).
    Dest(&newID).
    Fetch()

// Statement defined in DataSet:
// "insert-with-return": "INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id"
```

### UPDATE and DELETE

Use `Exec()` for operations that don't return data:

```go
// Update
ctx := context.Background()
err := store.Exec(
    goquery.NoTx,
    "UPDATE users SET email = $1 WHERE id = $2",
    "newemail@example.com",
    123,
)

// Delete
ctx := context.Background()
err = store.Exec(
    goquery.NoTx,
    "DELETE FROM users WHERE id = $1",
    123,
)

// Get rows affected
result, err := store.Execr(
    goquery.NoTx,
    "DELETE FROM users WHERE active = false",
)
if err == nil {
    fmt.Printf("Deleted %d rows\n", result.RowsAffected())
}
```

#### UPDATE/DELETE with RETURNING

```go
// Delete and get IDs of deleted rows
ctx := context.Background()
var deletedIDs []int
err := store.Select(`
    DELETE FROM users 
    WHERE active = false 
    RETURNING id
`).Context(ctx).Dest(&deletedIDs).Fetch()

fmt.Printf("Deleted user IDs: %v\n", deletedIDs)
```
