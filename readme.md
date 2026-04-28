# goquery - Comprehensive Documentation

**Version:** 1.0  
**License:** MIT  
**Go Version:** 1.18+


//new stuff
 - OnConnect function
 - DuckDb Support
 - support for native sqlite (sqlite) or cgo sqlite (sqlite3)
 - support for using driver Connectors with sqlx (important for duckdb)
 
---

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [Configuration](#configuration)
5. [Core Concepts](#core-concepts)
6. [DataStore Operations](#datastore-operations)
7. [Transactions](#transactions)
8. [Batch Operations](#batch-operations)
9. [Output Formats](#output-formats)
10. [Security Best Practices](#security-best-practices)
11. [Advanced Usage](#advanced-usage)
12. [Troubleshooting](#troubleshooting)
13. [API Reference](#api-reference)

---

## Overview

**goquery** is a lightweight Go library that simplifies database operations through a fluent, type-safe API. It provides a unified interface for multiple database backends while maintaining performance and safety.

### Key Features

- ✅ **Fluent API** - Chainable, readable query building
- ✅ **Multi-Database Support** - PostgreSQL (pgx), SQLite, Oracle, DuckDB
- ✅ **Type-Safe Mapping** - Automatic struct-to-row mapping via tags
- ✅ **Transaction Support** - Automatic rollback on panic, commit on success
- ✅ **Batch Operations** - High-performance bulk inserts (pgx)
- ✅ **Multiple Output Formats** - Structs, JSON, CSV
- ✅ **Connection Pooling** - Configurable pool settings
- ✅ **SQL Generation** - Auto-generate INSERT/SELECT from structs
- ✅ **Security First** - Parameterized queries prevent SQL injection

### Architecture

```
┌─────────────────┐
│   Your Code     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   DataStore     │  ← Unified interface
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌───────┐  ┌───────┐
│  PgxDb│  │SqlxDb │  ← Driver implementations
└───┬───┘  └───┬───┘
    │          │
    ▼          ▼
┌───────┐  ┌───────┐
│  pgx  │  │ sqlx  │  ← Underlying libraries
└───────┘  └───────┘
```

---

## Installation

```bash
go get github.com/usace/goquery
```

### Dependencies

goquery uses these excellent libraries:

- **pgx** (v4) - High-performance PostgreSQL driver
- **sqlx** - Extensions to database/sql
- **scany** - Struct scanning for SQL rows
- **go-strcase** - String case conversion for JSON

---

## Quick Start

### 1. Basic Connection

```go
package main

import (
    "log"
    "github.com/usace/goquery"
)

func main() {
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

    // Use the store...
}
```

### 2. Simple Query

```go
type User struct {
    ID    int    `db:"id"`
    Name  string `db:"name"`
    Email string `db:"email"`
}

// Query into struct slice
var users []User
err := store.Select("SELECT id, name, email FROM users").
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
var user User
err := store.Select("SELECT id, name, email FROM users WHERE id = $1").
    Params(42).
    Dest(&user).
    Fetch()
```

---

## Configuration

### RdbmsConfig Structure

```go
type RdbmsConfig struct {
    // Connection Settings
    Dbuser      string  // Database username
    Dbpass      string  // Database password
    Dbhost      string  // Database host (e.g., "localhost")
    Dbport      string  // Database port (e.g., "5432")
    Dbname      string  // Database name
    DbDriver    string  // Driver: "pgx", "sqlite", "godror"
    DbStore     string  // Store type: "pgx" or "sqlx"
    
    // Advanced Settings
    ExternalLib      string  // Path to external libs (Oracle)
    OnInit           string  // Initialization SQL (Oracle)
    DbDriverSettings string  // Additional driver parameters
    
    // Connection Pool Settings
    PoolMaxConns        int     // Maximum pool connections
    PoolMinConns        int     // Minimum pool connections
    PoolMaxConnLifetime string  // Max connection lifetime (e.g., "1h")
    PoolMaxConnIdle     string  // Max connection idle time (e.g., "30m")
}
```

### Configuration from Environment Variables

```go
config := goquery.RdbmsConfigFromEnv()
store, err := goquery.NewRdbmsDataStore(config)
```

Supported environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DBUSER` | Database username | (none) |
| `DBPASS` | Database password | (none) |
| `DBHOST` | Database host | (none) |
| `DBPORT` | Database port | `5432` |
| `DBNAME` | Database name | (none) |
| `DBDRIVER` | Database driver | (none) |
| `DBSTORE` | Store type (pgx/sqlx) | (none) |
| `DBDRIVER_PARAMS` | Additional parameters | (none) |
| `POOLMAXCONNS` | Max connections | (driver default) |
| `POOLMINCONNS` | Min connections | (driver default) |
| `POOLMAXCONNLIFETIME` | Max lifetime | (none) |
| `POOLMAXCONNIDLE` | Max idle time | (none) |

### DbDriverSettings - Advanced Configuration

The `DbDriverSettings` field allows you to pass additional parameters to the database driver:

#### PostgreSQL SSL Modes

```go
config := goquery.RdbmsConfig{
    // ... other settings ...
    DbDriverSettings: "sslmode=verify-full sslrootcert=/path/to/ca.crt",
}
```

**Available SSL modes:**
- `disable` - No SSL (insecure, development only)
- `require` - SSL required, no certificate verification (default)
- `verify-ca` - Verify server certificate against CA
- `verify-full` - Verify certificate and hostname (recommended for production)

#### Other PostgreSQL Parameters

```go
config.DbDriverSettings = "sslmode=require application_name=myapp connect_timeout=10"
```

**Common parameters:**
- `application_name` - Application name in logs
- `connect_timeout` - Connection timeout in seconds
- `statement_timeout` - Query timeout in milliseconds
- `search_path` - Default schema search path

#### Oracle (godror) Settings

```go
config := goquery.RdbmsConfig{
    DbDriver:    "godror",
    DbStore:     "sqlx",
    ExternalLib: "/usr/lib/oracle/instantclient",
    OnInit:      "ALTER SESSION SET NLS_DATE_FORMAT='YYYY-MM-DD HH24:MI:SS'",
    DbDriverSettings: "poolMinSessions=4 poolMaxSessions=100",
}
```

#### SQLite Settings

```go
config := goquery.RdbmsConfig{
    Dbname:           "/path/to/database.db",
    DbDriver:         "sqlite3",
    DbStore:          "sqlx",
    DbDriverSettings: "_journal_mode=WAL&_timeout=5000",
}
```

### Connection Pool Tuning

Duration strings use Go's `time.ParseDuration` format: `"300ms"`, `"1.5h"`, `"2h45m"`.

**Example configurations:**

```go
// High-traffic web application
config := goquery.RdbmsConfig{
    PoolMaxConns:        100,
    PoolMinConns:        10,
    PoolMaxConnLifetime: "1h",
    PoolMaxConnIdle:     "10m",
}

// Background worker
config := goquery.RdbmsConfig{
    PoolMaxConns:        10,
    PoolMinConns:        2,
    PoolMaxConnLifetime: "30m",
    PoolMaxConnIdle:     "5m",
}

// Development
config := goquery.RdbmsConfig{
    PoolMaxConns:        5,
    PoolMinConns:        1,
    PoolMaxConnLifetime: "5m",
    PoolMaxConnIdle:     "1m",
}
```

---

## Core Concepts

### DataStore Interface

The `DataStore` is your main entry point. All database operations go through it.

```go
type DataStore interface {
    // Query Operations
    Select(stmt ...string) *FluentSelect
    FetchRows(tx *Tx, input QueryInput) (Rows, error)
    
    // Insert Operations
    Insert(ds DataSet) *FluentInsert
    InsertRecs(tx *Tx, input InsertInput) error
    
    // Execute Operations
    Exec(tx *Tx, stmt string, params ...interface{}) error
    Execr(tx *Tx, stmt string, params ...interface{}) (ExecResult, error)
    MustExec(tx *Tx, stmt string, params ...interface{})
    MustExecr(tx *Tx, stmt string, params ...interface{}) ExecResult
    
    // Transaction Operations
    NewTransaction() (Tx, error)
    Transaction(tf TransactionFunction) error
}
```

### DataSet Pattern

DataSets organize your data structures and associated SQL statements:

```go
// 1. Define your struct with db tags
type Product struct {
    ID          int32   `db:"id" dbid:"SEQUENCE" idsequence:"products_id_seq"`
    Name        string  `db:"name"`
    Price       float64 `db:"price"`
    Description *string `db:"description"`  // Nullable field
}

// 2. Create a TableDataSet
var productsDS = goquery.TableDataSet{
    Name:   "products",
    Schema: "public",  // Optional schema
    Statements: goquery.Statements{
        "get-all":       "SELECT * FROM products ORDER BY name",
        "get-by-id":     "SELECT * FROM products WHERE id = $1",
        "search":        "SELECT * FROM products WHERE name ILIKE $1",
        "get-expensive": "SELECT * FROM products WHERE price > $1",
    },
    TableFields: Product{},  // Used for auto-generating INSERT statements
}

// 3. Use the dataset
var products []Product
err := store.Select().
    DataSet(&productsDS).
    StatementKey("get-all").
    Dest(&products).
    Fetch()
```

### Struct Tags

goquery uses struct tags to map between Go structs and database columns:

```go
type User struct {
    ID        int32   `db:"id" dbid:"SEQUENCE" idsequence:"users_id_seq"`
    Username  string  `db:"username"`
    Email     string  `db:"email"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt *time.Time `db:"updated_at"`  // Nullable
    Internal  string  `db:"-"`  // Ignored by goquery
}
```

**Tag reference:**

| Tag | Description | Example |
|-----|-------------|---------|
| `db:"column_name"` | Maps field to column | `db:"user_id"` |
| `db:"-"` | Ignores field | `db:"-"` |
| `dbid:"SEQUENCE"` | Auto-increment ID | `dbid:"SEQUENCE"` |
| `dbid:"AUTOINCREMENT"` | Auto-increment (SQLite) | `dbid:"AUTOINCREMENT"` |
| `idsequence:"seq_name"` | Sequence name (Postgres) | `idsequence:"users_id_seq"` |

---

## DataStore Operations

### SELECT Queries

#### Basic Select

```go
// Query all rows
var users []User
err := store.Select("SELECT * FROM users").
    Dest(&users).
    Fetch()

// Query single row
var user User
err := store.Select("SELECT * FROM users WHERE id = $1").
    Params(123).
    Dest(&user).
    Fetch()
```

#### With DataSet

```go
var users []User
err := store.Select().
    DataSet(&usersDS).
    StatementKey("get-all").
    Dest(&users).
    Fetch()
```

#### Dynamic Query Building

```go
// Apply() for internal query building (column names, table names)
var user User
err := store.Select("SELECT * FROM %s WHERE %s = $1").
    Apply("users", "id").
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
var users []User
err := store.Select("SELECT * FROM users").
    Suffix("WHERE active = true ORDER BY created_at DESC LIMIT 10").
    Dest(&users).
    Fetch()

// With parameters
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
    fmt.Printf("User: %s\n", user.Username)
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
        fmt.Printf("User: %s\n", user.Username)
        return nil
    }).
    Fetch()
```

### INSERT Operations

#### Single Insert

```go
user := User{
    Username: "john_doe",
    Email:    "john@example.com",
}

// Auto-generates INSERT statement from struct tags
err := store.Insert(&usersDS).
    Records(&user).
    Execute()
```

#### Bulk Insert

```go
users := []User{
    {Username: "alice", Email: "alice@example.com"},
    {Username: "bob", Email: "bob@example.com"},
    {Username: "charlie", Email: "charlie@example.com"},
}

err := store.Insert(&usersDS).
    Records(&users).
    Execute()
```

#### Batch Insert (High Performance)

Batch inserts are more efficient for large datasets (pgx only):

```go
users := make([]User, 10000)
// ... populate users ...

err := store.Insert(&usersDS).
    Records(&users).
    Batch(true).
    BatchSize(100).  // Send 100 records per batch
    Execute()
```

**Batch vs Bulk:**
- **Bulk Insert** - Uses multiple transactions, one per record
- **Batch Insert** - Uses pgx batching, single network roundtrip per batch (much faster)

#### INSERT with RETURNING

```go
// Get generated ID
var newID int
err := store.Select().
    DataSet(&usersDS).
    StatementKey("insert-with-return").
    Params("john_doe", "john@example.com").
    Dest(&newID).
    Fetch()

// Statement defined in DataSet:
// "insert-with-return": "INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id"
```

### UPDATE and DELETE

Use `Exec()` for operations that don't return data:

```go
// Update
err := store.Exec(
    goquery.NoTx,
    "UPDATE users SET email = $1 WHERE id = $2",
    "newemail@example.com",
    123,
)

// Delete
err := store.Exec(
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
var deletedIDs []int
err := store.Select(`
    DELETE FROM users 
    WHERE active = false 
    RETURNING id
`).Dest(&deletedIDs).Fetch()

fmt.Printf("Deleted user IDs: %v\n", deletedIDs)
```

---

## Transactions

### Automatic Transactions

The `Transaction()` method handles commit/rollback automatically:

```go
err := store.Transaction(func(tx goquery.Tx) {
    // All operations within this function are part of the transaction
    
    // Insert a user
    store.Insert(&usersDS).
        Records(&user).
        Tx(&tx).
        PanicOnErr(true).  // Panic on error to trigger rollback
        Execute()
    
    // Update related data
    store.MustExec(&tx, 
        "UPDATE accounts SET balance = balance + $1 WHERE user_id = $2",
        100.00, user.ID,
    )
    
    // If any operation panics, transaction is automatically rolled back
    // If function completes without panic, transaction is automatically committed
})

if err != nil {
    log.Printf("Transaction failed: %v\n", err)
}
```

**Transaction behavior:**
- ✅ Automatic commit if function completes successfully
- ✅ Automatic rollback if function panics
- ✅ Panic is caught and converted to error return value
- ✅ Errors are logged during rollback/commit failures

### Manual Transactions

For more control:

```go
tx, err := store.NewTransaction()
if err != nil {
    log.Fatal(err)
}

// Perform operations
err = store.Exec(&tx, "UPDATE users SET active = false WHERE id = $1", 123)
if err != nil {
    tx.Rollback()
    log.Fatal(err)
}

// Commit
err = tx.Commit()
if err != nil {
    log.Fatal(err)
}
```

### Error Handling in Transactions

Use `PanicOnErr()` or `MustExec()` to automatically rollback on error:

```go
err := store.Transaction(func(tx goquery.Tx) {
    // Option 1: PanicOnErr with fluent API
    store.Insert(&usersDS).
        Records(&user).
        Tx(&tx).
        PanicOnErr(true).
        Execute()
    
    // Option 2: MustExec (panics on error)
    store.MustExec(&tx, "UPDATE accounts SET balance = $1 WHERE user_id = $2", 0, user.ID)
    
    // Option 3: Manual panic
    result, err := store.Execr(&tx, "DELETE FROM temp_data")
    if err != nil {
        panic(err)  // Triggers rollback
    }
})
```

---

## Batch Operations

Batch operations dramatically improve performance for bulk inserts by reducing network round trips.

### When to Use Batches

- ✅ Inserting 100+ records
- ✅ High-throughput data pipelines
- ✅ ETL operations
- ✅ Importing large datasets

### Batch Insert Example

```go
// Generate test data
users := make([]User, 10000)
for i := range users {
    users[i] = User{
        Username: fmt.Sprintf("user_%d", i),
        Email:    fmt.Sprintf("user_%d@example.com", i),
    }
}

// Batch insert with error handling
err := store.Insert(&usersDS).
    Records(&users).
    Batch(true).
    BatchSize(500).  // 500 records per batch
    Execute()

if err != nil {
    log.Fatalf("Batch insert failed: %v", err)
}
```

### Batch Size Tuning

Choose batch size based on your data and network:

| Batch Size | Use Case | Notes |
|------------|----------|-------|
| 50-100 | Small records, slow network | Reduces packet size |
| 100-500 | Typical web applications | Good default |
| 500-1000 | Large records, fast network | Maximum throughput |
| 1000+ | Very large datasets, LAN | Diminishing returns |

### Performance Comparison

Example with 10,000 records:

| Method | Time | Network Trips |
|--------|------|---------------|
| Individual inserts | ~45s | 10,000 |
| Bulk insert (multi-tx) | ~25s | 10,000 |
| Batch insert (size 100) | ~2s | 100 |
| Batch insert (size 500) | ~1s | 20 |

### Batch Error Handling

Batch operations validate **each statement** in the batch:

```go
err := store.Insert(&usersDS).
    Records(&users).
    Batch(true).
    BatchSize(100).
    Execute()

if err != nil {
    // Error message includes which record failed
    // Example: "batch insert failed at record 342: duplicate key value"
    log.Printf("Batch failed: %v\n", err)
}
```

**Important:** If any record in a batch fails, the entire batch is rolled back (not the entire operation).

---

## Output Formats

### JSON Output

#### Stream to Writer (Recommended for large datasets)

```go
var buf bytes.Buffer
writer := bufio.NewWriter(&buf)

err := store.Select("SELECT * FROM users").
    OutputJson(writer).
    Fetch()

if err != nil {
    log.Fatal(err)
}

writer.Flush()
jsonBytes := buf.Bytes()
```

#### JSON Array

```go
err := store.Select("SELECT * FROM users").
    OutputJson(writer).
    IsJsonArray(true).  // Wrap results in []
    Fetch()

// Output: [{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]
```

#### Single Object

```go
err := store.Select("SELECT * FROM users WHERE id = $1").
    Params(1).
    OutputJson(writer).
    IsJsonArray(false).  // No array wrapper
    Fetch()

// Output: {"id":1,"name":"Alice"}
```

#### JSON Options

```go
err := store.Select("SELECT * FROM users").
    OutputJson(writer).
    CamelCase(true).      // Convert column names to camelCase
    OmitNull(true).       // Omit null fields
    DateFormat("2006-01-02").  // Custom date format
    Fetch()

// Output: {"userId":1,"userName":"Alice","createdAt":"2024-01-15"}
```

#### In-Memory JSON (Small datasets only)

```go
// Deprecated but still available
jsonBytes, err := store.Select("SELECT * FROM users LIMIT 10").
    FetchJSON()

if err != nil {
    log.Fatal(err)
}

fmt.Println(string(jsonBytes))
```

### CSV Output

```go
csv, err := store.Select("SELECT id, name, email FROM users").
    FetchCSV()

if err != nil {
    log.Fatal(err)
}

fmt.Println(csv)
// Output:
// "id","name","email"
// 1,"Alice","alice@example.com"
// 2,"Bob","bob@example.com"
```

**CSV Options:**

```go
csv, err := store.Select("SELECT * FROM users").
    CamelCase(true).      // Column headers in camelCase
    DateFormat("2006-01-02").  // Date formatting
    FetchCSV()
```

### Direct HTTP Response

```go
func usersHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    err := store.Select("SELECT * FROM users").
        OutputJson(w).  // Write directly to response
        IsJsonArray(true).
        CamelCase(true).
        Fetch()
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

---

## Security Best Practices

### 1. Always Use Parameterized Queries

**✅ CORRECT:**

```go
userID := getUserInput()
err := store.Select("SELECT * FROM users WHERE id = $1").
    Params(userID).  // Safe - uses parameter binding
    Dest(&user).
    Fetch()
```

**❌ WRONG:**

```go
userID := getUserInput()
sql := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID)  // SQL INJECTION!
err := store.Select(sql).Dest(&user).Fetch()
```

### 2. Apply() is for Internal Use Only

**✅ CORRECT:**

```go
// Internal constants only
tableName := "users"
columnName := "id"
userID := getUserInput()

err := store.Select("SELECT * FROM %s WHERE %s = $1").
    Apply(tableName, columnName).  // Internal values
    Params(userID).                // User input
    Dest(&user).
    Fetch()
```

**❌ WRONG:**

```go
// NEVER pass user input to Apply()
tableName := getUserInput()  // User could input: "users; DROP TABLE users--"
err := store.Select("SELECT * FROM %s").
    Apply(tableName).  // SQL INJECTION!
    Dest(&results).
    Fetch()
```

### 3. SSL/TLS Configuration

**Development:**

```go
config := goquery.RdbmsConfig{
    // ... other settings ...
    DbDriverSettings: "sslmode=disable",  // OK for local development
}
```

**Production (Recommended):**

```go
config := goquery.RdbmsConfig{
    // ... other settings ...
    DbDriverSettings: "sslmode=verify-full sslrootcert=/path/to/ca-cert.pem",
}
```

**Available SSL modes (PostgreSQL):**

| Mode | Encryption | Certificate Check | Use Case |
|------|-----------|-------------------|----------|
| `disable` | ❌ No | ❌ No | Local dev only |
| `require` | ✅ Yes | ❌ No | Default, basic security |
| `verify-ca` | ✅ Yes | ✅ CA only | Verify server identity |
| `verify-full` | ✅ Yes | ✅ CA + hostname | **Production (recommended)** |

### 4. Credential Management

**❌ DON'T** hardcode credentials:

```go
config := goquery.RdbmsConfig{
    Dbuser: "admin",
    Dbpass: "password123",  // NEVER hardcode passwords!
}
```

**✅ DO** use environment variables:

```go
config := goquery.RdbmsConfigFromEnv()
// Or:
config := goquery.RdbmsConfig{
    Dbuser: os.Getenv("DB_USER"),
    Dbpass: os.Getenv("DB_PASS"),
    // ...
}
```

**✅ BETTER** - Use secret managers:

```go
// Example with AWS Secrets Manager
secret := getSecretFromAWS("prod/db/credentials")
config := goquery.RdbmsConfig{
    Dbuser: secret["username"],
    Dbpass: secret["password"],
    // ...
}
```

### 5. Never Log Connection Strings

Connection strings contain passwords in plaintext. Never log them:

```go
// ❌ WRONG
log.Printf("Connecting with config: %+v", config)  // Logs password!

// ✅ CORRECT
log.Printf("Connecting to %s@%s:%s/%s", config.Dbuser, config.Dbhost, config.Dbport, config.Dbname)
```

### 6. Input Validation

Always validate user input before using it in queries:

```go
func getUser(idStr string) (*User, error) {
    // Validate input
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }
    if id <= 0 {
        return nil, errors.New("user ID must be positive")
    }
    
    // Safe to use in query
    var user User
    err = store.Select("SELECT * FROM users WHERE id = $1").
        Params(id).
        Dest(&user).
        Fetch()
    
    return &user, err
}
```

### 7. Least Privilege Principle

Use database users with minimal required permissions:

```sql
-- Create application user with limited permissions
CREATE USER app_user WITH PASSWORD 'secure_password';

-- Grant only what's needed
GRANT SELECT, INSERT, UPDATE ON users TO app_user;
GRANT SELECT ON products TO app_user;

-- Don't grant DELETE or DROP permissions unless necessary
```

---

## Advanced Usage

### Custom Row Scanning

```go
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

err := store.Insert(&usersDS).Records(&user).Execute()
```

### Multiple Databases

```go
// Connect to multiple databases
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
pgStore.Select("SELECT * FROM users").Dest(&pgUsers).Fetch()

var sqliteCache []CacheEntry
sqliteStore.Select("SELECT * FROM cache").Dest(&sqliteCache).Fetch()
```

### Schema Support

```go
// TableDataSet with schema
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
stmt, err := store.(*goquery.RdbmsDataStore).RdbmsDb().InsertStmt(&usersDS)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Generated INSERT:", stmt)
// Output: INSERT INTO users (id, name, email) VALUES (nextval('users_id_seq'), $1, $2)
```

---

## Troubleshooting

### Common Issues

#### 1. "no rows in result set"

```go
var user User
err := store.Select("SELECT * FROM users WHERE id = $1").
    Params(999).
    Dest(&user).
    Fetch()

// err will be "no rows in result set" if no user with id=999
```

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
ALTER SYSTEM SET max_connections = 200;
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
    Email *string `db:"email"`     // Nullable column
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

---

## API Reference

### FluentSelect Methods

| Method | Description | Example |
|--------|-------------|---------|
| `DataSet(ds DataSet)` | Use a DataSet | `.DataSet(&usersDS)` |
| `StatementKey(key string)` | Use named statement from DataSet | `.StatementKey("get-all")` |
| `Params(params ...interface{})` | Bind parameters | `.Params(123, "john")` |
| `Apply(vals ...interface{})` | Apply formatting (internal only) | `.Apply("users", "id")` |
| `Suffix(suffix string)` | Append to query | `.Suffix("LIMIT 10")` |
| `Dest(dest interface{})` | Set destination | `.Dest(&users)` |
| `Tx(tx *Tx)` | Use transaction | `.Tx(&tx)` |
| `LogSql(log bool)` | Log SQL statement | `.LogSql(true)` |
| `PanicOnErr(panic bool)` | Panic on error | `.PanicOnErr(true)` |
| `OutputJson(w io.Writer)` | Output as JSON | `.OutputJson(writer)` |
| `OutputCsv(w io.Writer)` | Output as CSV | `.OutputCsv(writer)` |
| `IsJsonArray(array bool)` | Wrap JSON in array | `.IsJsonArray(true)` |
| `CamelCase(camel bool)` | Convert to camelCase | `.CamelCase(true)` |
| `OmitNull(omit bool)` | Omit null fields | `.OmitNull(true)` |
| `DateFormat(format string)` | Custom date format | `.DateFormat("2006-01-02")` |
| `ForEachRow(fn RowFunction)` | Iterate rows | `.ForEachRow(func(r Rows) error {...})` |
| `Fetch()` | Execute query | `.Fetch()` |
| `FetchRows()` | Get row cursor | `.FetchRows()` |

### FluentInsert Methods

| Method | Description | Example |
|--------|-------------|---------|
| `Records(recs interface{})` | Set records to insert | `.Records(&user)` or `.Records(&users)` |
| `Tx(tx *Tx)` | Use transaction | `.Tx(&tx)` |
| `Batch(batch bool)` | Use batch mode | `.Batch(true)` |
| `BatchSize(size int)` | Set batch size | `.BatchSize(500)` |
| `PanicOnErr(panic bool)` | Panic on error | `.PanicOnErr(true)` |
| `Execute()` | Execute insert | `.Execute()` |

### Transaction Methods

| Method | Description |
|--------|-------------|
| `Commit()` | Commit transaction |
| `Rollback()` | Rollback transaction |
| `PgxTx()` | Get underlying pgx transaction |
| `SqlXTx()` | Get underlying sqlx transaction |
| `SqlTx()` | Get underlying sql.Tx |

### Rows Methods

| Method | Description |
|--------|-------------|
| `Next()` | Advance to next row |
| `Scan(...interface{})` | Scan into variables |
| `ScanStruct(interface{})` | Scan into struct |
| `Columns()` | Get column names |
| `ColumnTypes()` | Get column types |
| `Close()` | Close cursor |

---

## Examples

### Complete Web Application Example

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    
    "github.com/usace/goquery"
    "github.com/gorilla/mux"
)

type User struct {
    ID       int32   `db:"id" dbid:"SEQUENCE" idsequence:"users_id_seq"`
    Username string  `db:"username"`
    Email    string  `db:"email"`
    Active   bool    `db:"active"`
}

var usersDS = goquery.TableDataSet{
    Name: "users",
    Statements: goquery.Statements{
        "get-all":    "SELECT * FROM users ORDER BY username",
        "get-by-id":  "SELECT * FROM users WHERE id = $1",
        "get-active": "SELECT * FROM users WHERE active = true",
    },
    TableFields: User{},
}

var store goquery.DataStore

func main() {
    // Initialize database
    config := goquery.RdbmsConfigFromEnv()
    var err error
    store, err = goquery.NewRdbmsDataStore(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Setup routes
    r := mux.NewRouter()
    r.HandleFunc("/users", listUsers).Methods("GET")
    r.HandleFunc("/users/{id}", getUser).Methods("GET")
    r.HandleFunc("/users", createUser).Methods("POST")
    r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
    r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
    
    // Start server
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

func listUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    err := store.Select().
        DataSet(&usersDS).
        StatementKey("get-all").
        OutputJson(w).
        IsJsonArray(true).
        CamelCase(true).
        Fetch()
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func getUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    
    err = store.Select().
        DataSet(&usersDS).
        StatementKey("get-by-id").
        Params(id).
        OutputJson(w).
        CamelCase(true).
        Fetch()
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
    }
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    err := store.Insert(&usersDS).Records(&user).Execute()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    user.ID = int32(id)
    
    err = store.Exec(
        goquery.NoTx,
        "UPDATE users SET username = $1, email = $2, active = $3 WHERE id = $4",
        user.Username, user.Email, user.Active, user.ID,
    )
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    err = store.Exec(
        goquery.NoTx,
        "DELETE FROM users WHERE id = $1",
        id,
    )
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}
```

### Batch ETL Example

```go
package main

import (
    "encoding/csv"
    "log"
    "os"
    "strconv"
    
    "github.com/usace/goquery"
)

type Record struct {
    ID   int32  `db:"id" dbid:"AUTOINCREMENT"`
    Name string `db:"name"`
    Value float64 `db:"value"`
}

var recordsDS = goquery.TableDataSet{
    Name: "records",
    TableFields: Record{},
}

func main() {
    // Connect to database
    config := goquery.RdbmsConfigFromEnv()
    store, err := goquery.NewRdbmsDataStore(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Read CSV file
    file, err := os.Open("data.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    
    reader := csv.NewReader(file)
    rows, err := reader.ReadAll()
    if err != nil {
        log.Fatal(err)
    }
    
    // Parse records
    records := make([]Record, 0, len(rows)-1)
    for i, row := range rows {
        if i == 0 {
            continue // Skip header
        }
        
        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            log.Printf("Skipping row %d: %v", i, err)
            continue
        }
        
        records = append(records, Record{
            Name: row[1],
            Value: value,
        })
    }
    
    log.Printf("Importing %d records...", len(records))
    
    // Batch insert
    err = store.Insert(&recordsDS).
        Records(&records).
        Batch(true).
        BatchSize(500).
        Execute()
    
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Import complete!")
}
```

---

## Migration Guide

### From database/sql

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

// After (goquery)
var users []User
err := store.Select("SELECT id, name FROM users WHERE active = $1").
    Params(true).
    Dest(&users).
    Fetch()
```

### From GORM

```go
// Before (GORM)
var users []User
db.Where("active = ?", true).Find(&users)

// After (goquery)
var users []User
store.Select("SELECT * FROM users WHERE active = $1").
    Params(true).
    Dest(&users).
    Fetch()
```

---

## Contributing

We welcome contributions! Please:

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

---

## License

MIT License - see LICENSE file for details

---

## Support

- **Issues:** https://github.com/usace/goquery/issues
- **Documentation:** https://github.com/usace/goquery
- **Email:** support@usace.army.mil

---

**Version:** 1.0  
**Last Updated:** 2024  
**Maintained by:** U.S. Army Corps of Engineers
