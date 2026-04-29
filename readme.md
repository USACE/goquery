# goquery - Comprehensive Documentation

**Version:** 3.0  
**License:** MIT  
**Go Version:** 1.24+
---

## What's New in v3

Version 3.0 introduces significant architectural improvements and new database support:

- 🆕 **Module Versioning** - Import path is now `github.com/usace/goquery/v3`
- 🔌 **OnConnect Hook** - Execute initialization code when connections are established
- 🦆 **DuckDB Support** - Full support for DuckDB with spatial extensions
- 🗄️ **Dual SQLite Modes** - Choose between native Go (`sqlite`) or CGO (`sqlite3`) drivers
- 🔗 **Driver Connectors** - Direct `driver.Connector` support for advanced connection management
- 📦 **Modular Adapters** - Database adapters are now separate modules to reduce dependencies

See the [v3 Migration Guide](#v3-migration-guide) for upgrade instructions.

---

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [Configuration](#configuration)
5. [v3 New Features](#v3-new-features)
   - [OnConnect Hook](#onconnect-hook)
   - [DuckDB Support](#duckdb-support)
   - [SQLite: Native vs CGO](#sqlite-native-vs-cgo)
   - [Driver Connectors](#driver-connectors)
   - [Modular Adapters](#modular-adapters)
6. [Core Concepts](#core-concepts)
7. [DataStore Operations](#datastore-operations)
8. [Transactions](#transactions)
9. [Batch Operations](#batch-operations)
10. [Output Formats](#output-formats)
11. [Security Best Practices](#security-best-practices)
12. [Advanced Usage](#advanced-usage)
13. [Troubleshooting](#troubleshooting)
14. [API Reference](#api-reference)
15. [v3 Migration Guide](#v3-migration-guide)

---

## Overview

**goquery** is a lightweight Go library that simplifies database operations through a fluent, type-safe API. It provides a unified interface for multiple database backends while maintaining performance and safety.

### Key Features

- ✅ **Fluent API** - Chainable, readable query building
- ✅ **Multi-Database Support** - PostgreSQL (pgx), SQLite (native/CGO), Oracle, DuckDB
- ✅ **Type-Safe Mapping** - Automatic struct-to-row mapping via tags
- ✅ **Transaction Support** - Automatic rollback on panic, commit on success
- ✅ **Batch Operations** - High-performance bulk inserts (pgx)
- ✅ **Multiple Output Formats** - Structs, JSON, CSV
- ✅ **Connection Pooling** - Configurable pool settings
- ✅ **SQL Generation** - Auto-generate INSERT/SELECT from structs
- ✅ **Security First** - Parameterized queries prevent SQL injection
- ✅ **OnConnect Hooks** - Initialize connections with custom logic
- ✅ **Modular Architecture** - Import only the database drivers you need

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

### Basic Installation

```bash
go get github.com/usace/goquery/v3
```

### Database-Specific Installation

goquery v3 uses a modular adapter system. You must import the adapter for your database along with the driver:

#### PostgreSQL
```bash
go get github.com/usace/goquery/v3
go get github.com/usace/goquery/v3/adapters/postgres
go get github.com/jackc/pgx/v4
```

```go
import (
    _ "github.com/jackc/pgx/v4/stdlib"
    _ "github.com/usace/goquery/v3/adapters/postgres"
    "github.com/usace/goquery/v3"
)
```

#### DuckDB
```bash
go get github.com/usace/goquery/v3
go get github.com/usace/goquery/v3/adapters/duckdb
go get github.com/duckdb/duckdb-go/v2
```

```go
import (
    _ "github.com/duckdb/duckdb-go/v2"
    _ "github.com/usace/goquery/v3/adapters/duckdb"
    "github.com/usace/goquery/v3"
)
```

#### SQLite (Native Go - No CGO)
```bash
go get github.com/usace/goquery/v3
go get github.com/usace/goquery/v3/adapters/sqlite
go get modernc.org/sqlite
```

```go
import (
    _ "modernc.org/sqlite"
    _ "github.com/usace/goquery/v3/adapters/sqlite"
    "github.com/usace/goquery/v3"
)
```

#### SQLite (CGO)
```bash
go get github.com/usace/goquery/v3
go get github.com/usace/goquery/v3/adapters/sqlite
go get github.com/mattn/go-sqlite3
```

```go
import (
    _ "github.com/mattn/go-sqlite3"
    _ "github.com/usace/goquery/v3/adapters/sqlite"
    "github.com/usace/goquery/v3"
)
```

#### Oracle
```bash
go get github.com/usace/goquery/v3
go get github.com/usace/goquery/v3/adapters/oracle
go get github.com/godror/godror
```

```go
import (
    _ "github.com/godror/godror"
    _ "github.com/usace/goquery/v3/adapters/oracle"
    "github.com/usace/goquery/v3"
)
```

### Supported Databases

| Database | Driver Name | Adapter Import | Driver Import |
|----------|-------------|----------------|---------------|
| PostgreSQL | `pgx` | `github.com/usace/goquery/v3/adapters/postgres` | `github.com/jackc/pgx/v4/stdlib` |
| DuckDB | `duckdb` | `github.com/usace/goquery/v3/adapters/duckdb` | `github.com/duckdb/duckdb-go/v2` |
| SQLite (Native) | `sqlite` | `github.com/usace/goquery/v3/adapters/sqlite` | `modernc.org/sqlite` |
| SQLite (CGO) | `sqlite3` | `github.com/usace/goquery/v3/adapters/sqlite` | `github.com/mattn/go-sqlite3` |
| Oracle | `godror` | `github.com/usace/goquery/v3/adapters/oracle` | `github.com/godror/godror` |

### Core Dependencies

goquery uses these excellent libraries:

- **pgx** (v4) - High-performance PostgreSQL driver
- **sqlx** - Extensions to database/sql
- **scany** - Struct scanning for SQL rows
- **go-strcase** - String case conversion for JSON

---

## Quick Start

### 1. Basic Connection (PostgreSQL)

```go
package main

import (
    "log"
    _ "github.com/jackc/pgx/v4/stdlib"
    _ "github.com/usace/goquery/v3/adapters/postgres"
    "github.com/usace/goquery/v3"
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

## v3 New Features

### OnConnect Hook

The `OnConnect` hook allows you to execute initialization code when a database connection is established. This is perfect for loading extensions, setting session variables, or performing one-time setup.

#### Function Signature

```go
type RdbmsConfig struct {
    // ...
    OnConnect func(ds DataStore) error
}
```

**Parameters:**
- `ds DataStore` - The newly created DataStore instance
- **Returns:** `error` - Return an error to abort connection, or `nil` for success

#### How It Works

The `OnConnect` function is called automatically after the connection is established but before the DataStore is returned to the caller. If `OnConnect` returns an error, the connection is closed and the error is propagated.

**Implementation Location:** `rdbms_datastore.go:34-39`

```go
store := &RdbmsDataStore{db}
if config.OnConnect != nil {
    err := config.OnConnect(store)
    if err != nil {
        return nil, err
    }
}
return store, nil
```

#### Use Cases

##### 1. Loading DuckDB Extensions

```go
config := goquery.RdbmsConfig{
    DbDriver: "duckdb",
    DbStore:  "sqlx",
    Dbname:   "analytics.db",
    OnConnect: func(db goquery.DataStore) error {
        // Load spatial and HTTP filesystem extensions
        return db.Exec(goquery.NoTx, 
            "INSTALL spatial; LOAD spatial; INSTALL httpfs; LOAD httpfs")
    },
}

store, err := goquery.NewRdbmsDataStore(&config)
if err != nil {
    log.Fatal(err)
}
```

##### 2. Setting PostgreSQL Session Variables

```go
config := goquery.RdbmsConfig{
    DbDriver: "pgx",
    DbStore:  "pgx",
    Dbhost:   "localhost",
    Dbport:   "5432",
    Dbname:   "mydb",
    OnConnect: func(db goquery.DataStore) error {
        // Set search path and timezone
        err := db.Exec(goquery.NoTx, "SET search_path TO myschema, public")
        if err != nil {
            return err
        }
        return db.Exec(goquery.NoTx, "SET timezone TO 'UTC'")
    },
}
```

##### 3. Oracle NLS Settings

```go
config := goquery.RdbmsConfig{
    DbDriver: "godror",
    DbStore:  "sqlx",
    ExternalLib: "/usr/lib/oracle/instantclient",
    OnConnect: func(db goquery.DataStore) error {
        // Set date format for session
        return db.Exec(goquery.NoTx, 
            "ALTER SESSION SET NLS_DATE_FORMAT='YYYY-MM-DD HH24:MI:SS'")
    },
}
```

##### 4. Creating Temporary Tables

```go
config := goquery.RdbmsConfig{
    DbDriver: "sqlite",
    DbStore:  "sqlx",
    Dbname:   ":memory:",
    OnConnect: func(db goquery.DataStore) error {
        // Create temporary lookup table
        return db.Exec(goquery.NoTx, `
            CREATE TEMP TABLE session_cache (
                key TEXT PRIMARY KEY,
                value TEXT,
                expires INTEGER
            )`)
    },
}
```

##### 5. Enabling SQLite Extensions

```go
config := goquery.RdbmsConfig{
    DbDriver: "sqlite3",
    DbStore:  "sqlx",
    Dbname:   "mydb.sqlite",
    OnConnect: func(db goquery.DataStore) error {
        // Enable foreign keys (disabled by default in SQLite)
        return db.Exec(goquery.NoTx, "PRAGMA foreign_keys = ON")
    },
}
```

##### 6. Multiple Initialization Steps with Error Handling

```go
config := goquery.RdbmsConfig{
    DbDriver: "duckdb",
    DbStore:  "sqlx",
    Dbname:   "data.duckdb",
    OnConnect: func(db goquery.DataStore) error {
        // Multiple initialization steps
        initCommands := []string{
            "INSTALL spatial",
            "LOAD spatial",
            "INSTALL httpfs",
            "LOAD httpfs",
            "SET memory_limit='4GB'",
            "SET threads=4",
        }
        
        for _, cmd := range initCommands {
            if err := db.Exec(goquery.NoTx, cmd); err != nil {
                return fmt.Errorf("init command failed [%s]: %w", cmd, err)
            }
        }
        
        log.Println("DuckDB initialized with spatial and httpfs extensions")
        return nil
    },
}
```

#### OnConnect vs OnInit

**Deprecated: `OnInit` string field (Oracle only)**

The older `OnInit` string field is still supported for backward compatibility but is Oracle-specific and limited to a single SQL statement. The new `OnConnect` function is recommended because it:

- ✅ Works with **all database types**
- ✅ Supports **multiple commands**
- ✅ Provides **error handling and reporting**
- ✅ Allows **conditional logic** and **logging**
- ✅ Has access to the **full DataStore interface**

```go
// Old way (Oracle only, single statement)
config := goquery.RdbmsConfig{
    OnInit: "ALTER SESSION SET NLS_DATE_FORMAT='YYYY-MM-DD'",
}

// New way (all databases, multiple statements, error handling)
config := goquery.RdbmsConfig{
    OnConnect: func(db goquery.DataStore) error {
        err := db.Exec(goquery.NoTx, "ALTER SESSION SET NLS_DATE_FORMAT='YYYY-MM-DD'")
        if err != nil {
            return err
        }
        return db.Exec(goquery.NoTx, "ALTER SESSION SET NLS_TIMESTAMP_FORMAT='YYYY-MM-DD HH24:MI:SS'")
    },
}
```

---

### DuckDB Support

goquery v3 adds first-class support for DuckDB, the high-performance analytical database. DuckDB is perfect for:

- 📊 OLAP workloads and analytics
- 🗺️ Geospatial data processing (with spatial extension)
- 📁 Querying Parquet, CSV, and JSON files directly
- 🌐 Reading data from HTTP/S3 (with httpfs extension)
- 💾 Embedded analytics in Go applications

#### Installation

```bash
go get github.com/usace/goquery/v3
go get github.com/usace/goquery/v3/adapters/duckdb
go get github.com/duckdb/duckdb-go/v2
```

#### Required Imports

```go
import (
    _ "github.com/duckdb/duckdb-go/v2"               // DuckDB driver
    _ "github.com/usace/goquery/v3/adapters/duckdb"  // goquery adapter
    "github.com/usace/goquery/v3"
)
```

#### Basic DuckDB Configuration

```go
config := goquery.RdbmsConfig{
    Dbname:   "analytics.duckdb",  // File path
    DbDriver: "duckdb",            // Driver name
    DbStore:  "sqlx",              // Use sqlx store
}

store, err := goquery.NewRdbmsDataStore(&config)
if err != nil {
    log.Fatal(err)
}
```

#### In-Memory DuckDB

```go
config := goquery.RdbmsConfig{
    Dbname:   ":memory:",  // In-memory database
    DbDriver: "duckdb",
    DbStore:  "sqlx",
}
```

#### DuckDB with Extensions

```go
config := goquery.RdbmsConfig{
    Dbname:   "geo_analytics.duckdb",
    DbDriver: "duckdb",
    DbStore:  "sqlx",
    OnConnect: func(db goquery.DataStore) error {
        // Install and load DuckDB extensions
        return db.Exec(goquery.NoTx, 
            "INSTALL spatial; LOAD spatial; INSTALL httpfs; LOAD httpfs")
    },
}

store, err := goquery.NewRdbmsDataStore(&config)
```

#### Dialect Implementation

**Location:** `adapters/duckdb/dialect_duckdb.go`

```go
package duckdb

import (
    "fmt"
    "github.com/usace/goquery/v3"
)

const (
    registryName string = "duckdb"
)

func init() {
    // Auto-registers when adapter is imported
    goquery.DbRegistry[registryName] = DuckdbDialect
}

var DuckdbDialect = goquery.DbDialect{
    TableExistsStmt: `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)`,
    Bind: func(field string, i int) string {
        return fmt.Sprintf("$%d", i+1)  // PostgreSQL-style parameters
    },
    Url: func(config *goquery.RdbmsConfig) string {
        return config.Dbname  // File path or :memory:
    },
}
```

#### Complete DuckDB Example: Spatial Query

```go
package main

import (
    "fmt"
    "log"
    
    _ "github.com/duckdb/duckdb-go/v2"
    _ "github.com/usace/goquery/v3/adapters/duckdb"
    "github.com/usace/goquery/v3"
)

func main() {
    config := goquery.RdbmsConfig{
        Dbname:   "spatial.duckdb",
        DbDriver: "duckdb",
        DbStore:  "sqlx",
        OnConnect: func(db goquery.DataStore) error {
            log.Println("Loading spatial extensions...")
            return db.Exec(goquery.NoTx, "INSTALL spatial; LOAD spatial")
        },
    }
    
    store, err := goquery.NewRdbmsDataStore(&config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Query GeoPackage file directly
    resource := "data/countries.gpkg"
    layer := "boundaries"
    
    query := fmt.Sprintf(
        "SELECT name, population, ST_Area(geom) as area FROM ST_Read('%s', layer='%s')", 
        resource, layer)
    
    type Country struct {
        Name       string  `db:"name"`
        Population int64   `db:"population"`
        Area       float64 `db:"area"`
    }
    
    var countries []Country
    err = store.Select(query).
        Dest(&countries).
        Fetch()
    
    if err != nil {
        log.Fatal(err)
    }
    
    for _, c := range countries {
        fmt.Printf("%s: population=%d, area=%.2f\n", c.Name, c.Population, c.Area)
    }
}
```

#### Querying Parquet Files

```go
// Query Parquet file directly without loading into database
var results []map[string]interface{}

err := store.Select("SELECT * FROM 'data/events.parquet' WHERE date >= '2024-01-01'").
    ForEachRow(func(row goquery.Rows) error {
        rowMap, err := row.ToMap()
        if err != nil {
            return err
        }
        results = append(results, rowMap)
        return nil
    }).
    Fetch()
```

#### Reading from HTTP/S3

```go
config := goquery.RdbmsConfig{
    Dbname:   ":memory:",
    DbDriver: "duckdb",
    DbStore:  "sqlx",
    OnConnect: func(db goquery.DataStore) error {
        return db.Exec(goquery.NoTx, "INSTALL httpfs; LOAD httpfs")
    },
}

store, _ := goquery.NewRdbmsDataStore(&config)

// Query CSV from URL
var data []YourStruct
err := store.Select("SELECT * FROM 'https://example.com/data.csv'").
    Dest(&data).
    Fetch()
```

#### DuckDB Performance Tips

1. **Set Memory Limit:**
```go
OnConnect: func(db goquery.DataStore) error {
    return db.Exec(goquery.NoTx, "SET memory_limit='8GB'")
}
```

2. **Configure Thread Count:**
```go
OnConnect: func(db goquery.DataStore) error {
    return db.Exec(goquery.NoTx, "SET threads=8")
}
```

3. **Use Persistent Storage for Large Datasets:**
```go
// Instead of :memory:, use a file
Dbname: "analytics.duckdb",
```

4. **Enable Progress Bar for Long Queries:**
```go
OnConnect: func(db goquery.DataStore) error {
    return db.Exec(goquery.NoTx, "SET enable_progress_bar=true")
}
```

---

### SQLite: Native vs CGO

goquery v3 supports **two different SQLite drivers**, giving you flexibility based on your deployment requirements.

#### Comparison Table

| Feature | **Native Go** (`sqlite`) | **CGO** (`sqlite3`) |
|---------|--------------------------|---------------------|
| **Package** | `modernc.org/sqlite` | `github.com/mattn/go-sqlite3` |
| **CGO Required** | ❌ No | ✅ Yes |
| **C Compiler** | ❌ Not needed | ✅ Required |
| **Cross-compilation** | ✅ Simple | ❌ Complex |
| **Pure Go** | ✅ Yes | ❌ No |
| **Performance** | Good (90-95% of CGO) | Excellent (100%) |
| **Binary Size** | Larger (~10MB+) | Smaller (~2-3MB) |
| **Build Speed** | Fast | Slower (C compilation) |
| **Docker Alpine** | ✅ Works easily | ⚠️ Needs build-base |
| **Production Ready** | ✅ Yes | ✅ Yes |
| **Driver Name** | `"sqlite"` | `"sqlite3"` |
| **Best For** | Cloud deployments, cross-platform builds | Maximum performance, existing CGO setup |

#### Dialect Implementation

**Location:** `adapters/sqlite/dialect_sqlite.go`

Both drivers use the **same goquery adapter** - just import the driver you want:

```go
package sqlite

import "github.com/usace/goquery/v3"

const (
    registryNameCgo      string = "sqlite3"  // CGO driver
    registryNameNativeGo string = "sqlite"   // Native Go driver
)

func init() {
    // Register both drivers with same dialect
    goquery.DbRegistry[registryNameCgo] = SqliteDialect
    goquery.DbRegistry[registryNameNativeGo] = SqliteDialect
}

var SqliteDialect = goquery.DbDialect{
    TableExistsStmt: `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`,
    Bind: func(field string, i int) string {
        return "?"  // SQLite uses ? placeholders
    },
    Seq: func(sequence string) string {
        return ""  // No sequences in SQLite
    },
    Url: func(config *goquery.RdbmsConfig) string {
        return config.Dbname  // File path
    },
}
```

#### Using Native Go SQLite (No CGO)

**Installation:**
```bash
go get github.com/usace/goquery/v3
go get github.com/usace/goquery/v3/adapters/sqlite
go get modernc.org/sqlite
```

**Code:**
```go
package main

import (
    "log"
    
    _ "modernc.org/sqlite"                          // Native Go driver
    _ "github.com/usace/goquery/v3/adapters/sqlite"
    "github.com/usace/goquery/v3"
)

func main() {
    config := goquery.RdbmsConfig{
        Dbname:   "./myapp.db",
        DbDriver: "sqlite",      // Use native Go driver
        DbStore:  "sqlx",
        DbDriverSettings: "_journal_mode=WAL&_timeout=5000",
    }
    
    store, err := goquery.NewRdbmsDataStore(&config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use store...
}
```

**Dockerfile (No build dependencies needed):**
```dockerfile
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/myapp /myapp
CMD ["/myapp"]
```

#### Using CGO SQLite

**Installation:**
```bash
go get github.com/usace/goquery/v3
go get github.com/usace/goquery/v3/adapters/sqlite
go get github.com/mattn/go-sqlite3
```

**Code:**
```go
package main

import (
    "log"
    
    _ "github.com/mattn/go-sqlite3"                 // CGO driver
    _ "github.com/usace/goquery/v3/adapters/sqlite"
    "github.com/usace/goquery/v3"
)

func main() {
    config := goquery.RdbmsConfig{
        Dbname:   "./myapp.db",
        DbDriver: "sqlite3",     // Use CGO driver
        DbStore:  "sqlx",
        DbDriverSettings: "_journal_mode=WAL&_timeout=5000&_busy_timeout=10000",
    }
    
    store, err := goquery.NewRdbmsDataStore(&config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use store...
}
```

**Dockerfile (Requires build dependencies):**
```dockerfile
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
COPY --from=builder /app/myapp /myapp
CMD ["/myapp"]
```

#### Common SQLite Configuration

Both drivers support the same connection parameters via `DbDriverSettings`:

```go
config := goquery.RdbmsConfig{
    Dbname:   "./myapp.db",
    DbDriver: "sqlite",  // or "sqlite3"
    DbStore:  "sqlx",
    DbDriverSettings: "_journal_mode=WAL&_timeout=5000&_busy_timeout=10000&cache=shared",
}
```

**Common parameters:**
- `_journal_mode=WAL` - Write-Ahead Logging for better concurrency
- `_timeout=5000` - Busy timeout in milliseconds
- `_busy_timeout=10000` - Alternative busy timeout syntax
- `cache=shared` - Share cache between connections
- `mode=ro` - Read-only mode
- `mode=memory` - In-memory database

#### Enabling Foreign Keys

SQLite disables foreign key constraints by default. Enable them with `OnConnect`:

```go
config := goquery.RdbmsConfig{
    Dbname:   "./myapp.db",
    DbDriver: "sqlite",
    DbStore:  "sqlx",
    OnConnect: func(db goquery.DataStore) error {
        return db.Exec(goquery.NoTx, "PRAGMA foreign_keys = ON")
    },
}
```

#### In-Memory Database

```go
config := goquery.RdbmsConfig{
    Dbname:   ":memory:",  // In-memory database
    DbDriver: "sqlite",
    DbStore:  "sqlx",
}
```

#### Which One Should You Use?

**Use Native Go (`sqlite`) if:**
- ✅ Deploying to cloud platforms (AWS Lambda, Google Cloud Run)
- ✅ Cross-compiling for multiple platforms
- ✅ Want pure Go dependencies
- ✅ Building Docker images from scratch/alpine
- ✅ Performance is "good enough" (it usually is)

**Use CGO (`sqlite3`) if:**
- ✅ Need maximum performance (5-10% faster)
- ✅ Already have CGO in your build pipeline
- ✅ Building custom SQLite extensions
- ✅ Need specific SQLite compilation flags

**Default recommendation:** Start with **Native Go** (`sqlite`) for simplicity. Switch to CGO only if profiling shows SQLite as a bottleneck.

---

### Driver Connectors

goquery v3 adds support for `database/sql/driver.Connector`, allowing you to use custom connection logic and advanced driver features.

#### Configuration Field

```go
type RdbmsConfig struct {
    // ...
    
    // If Connector is populated, it will be used to create all connections
    // All other connection parameters (Dbhost, Dbport, etc.) will be IGNORED
    Connector driver.Connector
}
```

#### How It Works

When `Connector` is set, goquery bypasses the normal connection string generation and uses the Connector directly:

**Implementation Location:** `sqlx_db.go:88-95`

```go
func NewSqlxConnection(config *RdbmsConfig) (SqlxDb, error) {
    dialect, err := getDialect(config.DbDriver)
    if err != nil {
        return SqlxDb{}, err
    }
    
    if config.Connector != nil {
        // Use Connector directly
        sqlcon := sql.OpenDB(config.Connector)
        con := sqlx.NewDb(sqlcon, config.DbDriver)
        return SqlxDb{con, dialect}, nil
    } else {
        // Use connection string (standard path)
        dburl := dialect.Url(config)
        con, err := sqlx.Connect(config.DbDriver, dburl)
        return SqlxDb{con, dialect}, err
    }
}
```

#### Why Use Connectors?

Connectors provide advanced capabilities:

1. **Driver-Specific Configuration** - Set options not available via connection string
2. **Connection Callbacks** - Execute code for every connection in the pool
3. **Custom Authentication** - Implement complex auth logic
4. **Connection Lifecycle** - Manage connection creation and destruction
5. **Performance Tuning** - Configure thread counts, memory limits, cache sizes

#### DuckDB Connector Example

DuckDB heavily relies on Connectors for proper configuration:

```go
package main

import (
    "context"
    "database/sql"
    "database/sql/driver"
    "log"
    
    duckdb "github.com/duckdb/duckdb-go/v2"
    _ "github.com/usace/goquery/v3/adapters/duckdb"
    "github.com/usace/goquery/v3"
)

func main() {
    // Create DuckDB connector with custom configuration
    connector, err := duckdb.NewConnector("analytics.duckdb", func(execer driver.ExecerContext) error {
        // This function runs for EVERY connection in the pool
        bootQueries := []string{
            "SET memory_limit='8GB'",
            "SET threads=8",
            "INSTALL spatial",
            "LOAD spatial",
            "INSTALL httpfs",
            "LOAD httpfs",
            "SET enable_progress_bar=true",
        }
        
        for _, query := range bootQueries {
            _, err := execer.ExecContext(context.Background(), query, nil)
            if err != nil {
                return fmt.Errorf("boot query failed [%s]: %w", query, err)
            }
        }
        
        log.Println("DuckDB connection configured")
        return nil
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Use Connector with goquery
    config := goquery.RdbmsConfig{
        DbDriver:  "duckdb",
        DbStore:   "sqlx",
        Connector: connector,  // All other connection params ignored
    }
    
    store, err := goquery.NewRdbmsDataStore(&config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use store for analytics queries
    var results []AnalyticsRow
    err = store.Select("SELECT * FROM large_dataset WHERE date >= '2024-01-01'").
        Dest(&results).
        Fetch()
}
```

#### Connector vs OnConnect

Both `Connector` and `OnConnect` can execute initialization code, but they work at different levels:

| Feature | **Connector** | **OnConnect** |
|---------|---------------|---------------|
| **Executes When** | Every connection in pool | Once when DataStore is created |
| **Scope** | Individual connection | DataStore instance |
| **Access To** | driver.ExecerContext | Full DataStore interface |
| **Use Case** | Connection-level settings | DataStore-level initialization |
| **Driver Support** | Driver-specific | All drivers |
| **Typical Usage** | DuckDB thread/memory config | Load extensions, create temp tables |

**Example using both:**

```go
connector, _ := duckdb.NewConnector("data.duckdb", func(execer driver.ExecerContext) error {
    // Runs for EVERY connection in pool
    _, err := execer.ExecContext(context.Background(), "SET threads=4", nil)
    return err
})

config := goquery.RdbmsConfig{
    DbDriver:  "duckdb",
    DbStore:   "sqlx",
    Connector: connector,
    OnConnect: func(db goquery.DataStore) error {
        // Runs ONCE when DataStore is created
        // Load extensions (persistent across connections)
        err := db.Exec(goquery.NoTx, "INSTALL spatial; LOAD spatial")
        if err != nil {
            return err
        }
        
        // Create tables
        return db.Exec(goquery.NoTx, `
            CREATE TABLE IF NOT EXISTS analytics (
                id INTEGER PRIMARY KEY,
                event_date DATE,
                value DOUBLE
            )`)
    },
}
```

#### PostgreSQL Connector Example (Custom SSL)

```go
import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/stdlib"
    _ "github.com/usace/goquery/v3/adapters/postgres"
    "github.com/usace/goquery/v3"
)

func main() {
    // Load custom CA certificate
    caCert, err := ioutil.ReadFile("/path/to/ca-cert.pem")
    if err != nil {
        log.Fatal(err)
    }
    
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)
    
    // Create pgx config with custom TLS
    connConfig, err := pgx.ParseConfig("postgres://user:pass@localhost:5432/mydb")
    if err != nil {
        log.Fatal(err)
    }
    
    connConfig.TLSConfig = &tls.Config{
        RootCAs:            caCertPool,
        InsecureSkipVerify: false,
        ServerName:         "postgres.example.com",
    }
    
    // Create connector
    connector := stdlib.GetConnector(*connConfig)
    
    // Use with goquery
    config := goquery.RdbmsConfig{
        DbDriver:  "pgx",
        DbStore:   "pgx",
        Connector: connector,
    }
    
    store, err := goquery.NewRdbmsDataStore(&config)
    if err != nil {
        log.Fatal(err)
    }
}
```

---

### Modular Adapters

goquery v3 introduces a **modular adapter architecture** where database-specific code is separated into independent Go modules. This reduces dependencies and binary size.

#### Architecture Design

**Before v3:** All database adapters were bundled with goquery core, forcing every application to pull in dependencies for all supported databases.

**v3 Approach:** Each adapter is a separate Go module that applications import only when needed.

```
goquery/v3/
├── go.mod                          # Core module (no database drivers)
├── datastore.go
├── config.go
└── adapters/
    ├── postgres/
    │   ├── go.mod                  # module github.com/usace/goquery/v3/adapters/postgres
    │   ├── go.sum
    │   └── dialect_pg.go
    ├── duckdb/
    │   ├── go.mod                  # module github.com/usace/goquery/v3/adapters/duckdb
    │   ├── go.sum
    │   └── dialect_duckdb.go
    ├── sqlite/
    │   ├── go.mod                  # module github.com/usace/goquery/v3/adapters/sqlite
    │   ├── go.sum
    │   └── dialect_sqlite.go
    └── oracle/
        ├── go.mod                  # module github.com/usace/goquery/v3/adapters/oracle
        ├── go.sum
        └── dialect_oracle.go
```

#### Registry Pattern

Each adapter registers itself with the global `DbRegistry` using an `init()` function:

**Example:** `adapters/duckdb/dialect_duckdb.go`

```go
package duckdb

import (
    "fmt"
    "github.com/usace/goquery/v3"
)

const (
    registryName string = "duckdb"
)

func init() {
    // Automatically registers when package is imported
    goquery.DbRegistry[registryName] = DuckdbDialect
}

var DuckdbDialect = goquery.DbDialect{
    TableExistsStmt: `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)`,
    Bind: func(field string, i int) string {
        return fmt.Sprintf("$%d", i+1)
    },
    Url: func(config *goquery.RdbmsConfig) string {
        return config.Dbname
    },
}
```

#### Core Registry Implementation

**Location:** `datastore.go`

```go
type DialectRegistry map[string]DbDialect

// Global registry populated by adapter init() functions
var DbRegistry = make(DialectRegistry)

// Lookup with helpful error message
func getDialect(driver string) (DbDialect, error) {
    if dialect, ok := DbRegistry[driver]; ok {
        return dialect, nil
    }
    return DbDialect{}, fmt.Errorf(
        "uninitialized or unsupported driver '%s'. "+
        "Make sure you imported the adapter: "+
        "import _ \"github.com/usace/goquery/v3/adapters/%s\"",
        driver, driver)
}
```

#### Import Patterns

##### Single Database Application

```go
package main

import (
    "log"
    
    // Core goquery
    "github.com/usace/goquery/v3"
    
    // PostgreSQL driver + adapter
    _ "github.com/jackc/pgx/v4/stdlib"
    _ "github.com/usace/goquery/v3/adapters/postgres"
)

func main() {
    config := goquery.RdbmsConfig{
        DbDriver: "pgx",
        DbStore:  "pgx",
        // ... connection details
    }
    
    store, err := goquery.NewRdbmsDataStore(&config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use store...
}
```

**go.mod dependencies (only PostgreSQL):**
```go
require (
    github.com/usace/goquery/v3 v3.0.0
    github.com/usace/goquery/v3/adapters/postgres v3.0.0
    github.com/jackc/pgx/v4 v4.18.0
)
```

##### Multi-Database Application

```go
package main

import (
    "github.com/usace/goquery/v3"
    
    // Import only the adapters you need
    _ "github.com/jackc/pgx/v4/stdlib"
    _ "github.com/usace/goquery/v3/adapters/postgres"
    
    _ "github.com/duckdb/duckdb-go/v2"
    _ "github.com/usace/goquery/v3/adapters/duckdb"
    
    _ "modernc.org/sqlite"
    _ "github.com/usace/goquery/v3/adapters/sqlite"
)

func main() {
    // PostgreSQL for OLTP
    pgStore, _ := goquery.NewRdbmsDataStore(&goquery.RdbmsConfig{
        DbDriver: "pgx",
        DbStore:  "pgx",
        // ... connection details
    })
    
    // DuckDB for analytics
    duckStore, _ := goquery.NewRdbmsDataStore(&goquery.RdbmsConfig{
        DbDriver: "duckdb",
        DbStore:  "sqlx",
        Dbname:   "analytics.duckdb",
    })
    
    // SQLite for caching
    cacheStore, _ := goquery.NewRdbmsDataStore(&goquery.RdbmsConfig{
        DbDriver: "sqlite",
        DbStore:  "sqlx",
        Dbname:   "./cache.db",
    })
    
    // Use each store for its purpose...
}
```

#### All Available Adapters

| Database | Driver Name | Adapter Import | Driver Import |
|----------|-------------|----------------|---------------|
| PostgreSQL | `pgx` | `github.com/usace/goquery/v3/adapters/postgres` | `github.com/jackc/pgx/v4/stdlib` |
| DuckDB | `duckdb` | `github.com/usace/goquery/v3/adapters/duckdb` | `github.com/duckdb/duckdb-go/v2` |
| SQLite (Native) | `sqlite` | `github.com/usace/goquery/v3/adapters/sqlite` | `modernc.org/sqlite` |
| SQLite (CGO) | `sqlite3` | `github.com/usace/goquery/v3/adapters/sqlite` | `github.com/mattn/go-sqlite3` |
| Oracle | `godror` | `github.com/usace/goquery/v3/adapters/oracle` | `github.com/godror/godror` |

#### Benefits of Modular Adapters

1. **Smaller Dependencies**
   - Applications only pull in drivers they actually use
   - `go.mod` remains clean and focused

2. **Smaller Binary Size**
   - Unused drivers aren't compiled into the binary
   - Example: PostgreSQL-only app doesn't include DuckDB (~50MB savings)

3. **Faster Build Times**
   - Less code to compile
   - No CGO compilation for drivers you don't use

4. **Independent Versioning**
   - Each adapter can evolve independently
   - Update one adapter without affecting others

5. **Clear Error Messages**
   - Forget to import adapter? You get a helpful error:
   ```
   uninitialized or unsupported driver 'duckdb'.
   Make sure you imported the adapter:
   import _ "github.com/usace/goquery/v3/adapters/duckdb"
   ```

#### Error: Missing Adapter Import

If you forget to import the adapter, you'll see this error at runtime:

```go
// ❌ WRONG - Missing adapter import
package main

import (
    _ "github.com/duckdb/duckdb-go/v2"
    "github.com/usace/goquery/v3"
)

func main() {
    config := goquery.RdbmsConfig{
        DbDriver: "duckdb",
        DbStore:  "sqlx",
        Dbname:   "data.duckdb",
    }
    
    store, err := goquery.NewRdbmsDataStore(&config)
    // Error: uninitialized or unsupported driver 'duckdb'.
    // Make sure you imported the adapter:
    // import _ "github.com/usace/goquery/v3/adapters/duckdb"
}
```

```go
// ✅ CORRECT - Both driver AND adapter imported
package main

import (
    _ "github.com/duckdb/duckdb-go/v2"                   // Driver
    _ "github.com/usace/goquery/v3/adapters/duckdb"      // Adapter (required!)
    "github.com/usace/goquery/v3"
)

func main() {
    config := goquery.RdbmsConfig{
        DbDriver: "duckdb",
        DbStore:  "sqlx",
        Dbname:   "data.duckdb",
    }
    
    store, err := goquery.NewRdbmsDataStore(&config)
    // ✅ Works!
}
```

#### Creating Custom Adapters

You can create custom adapters for unsupported databases:

```go
package mydb

import (
    "fmt"
    "github.com/usace/goquery/v3"
)

func init() {
    goquery.DbRegistry["mydb"] = goquery.DbDialect{
        TableExistsStmt: `SELECT table_name FROM information_schema.tables WHERE table_name = $1`,
        Bind: func(field string, i int) string {
            return fmt.Sprintf("$%d", i+1)
        },
        Seq: func(sequence string) string {
            return fmt.Sprintf("nextval('%s')", sequence)
        },
        Url: func(config *goquery.RdbmsConfig) string {
            return fmt.Sprintf("mydb://%s:%s@%s:%s/%s",
                config.Dbuser,
                config.Dbpass,
                config.Dbhost,
                config.Dbport,
                config.Dbname)
        },
    }
}
```

Then import your custom adapter:

```go
import (
    _ "yourmodule/mydb"
    "github.com/usace/goquery/v3"
)
```



### RdbmsConfig Structure

```go
type RdbmsConfig struct {
    // Connection Settings
    Dbuser      string  // Database username
    Dbpass      string  // Database password
    Dbhost      string  // Database host (e.g., "localhost")
    Dbport      string  // Database port (e.g., "5432")
    Dbname      string  // Database name or file path (for SQLite/DuckDB)
    DbDriver    string  // Driver: "pgx", "sqlite", "sqlite3", "duckdb", "godror"
    DbStore     string  // Store type: "pgx" or "sqlx"
    
    // Advanced Settings
    ExternalLib      string  // Path to external libs (Oracle)
    OnInit           string  // Initialization SQL (Oracle - deprecated, use OnConnect)
    DbDriverSettings string  // Additional driver parameters
    
    // v3 New Features
    OnConnect  func(ds DataStore) error  // Hook function called when connection is established
    Connector  driver.Connector          // Direct driver.Connector (bypasses other connection settings)
    
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

## v3 Migration Guide

### Overview of Breaking Changes

Version 3.0 introduces breaking changes that require code updates. The main changes are:

1. **Module path includes `/v3`**
2. **Adapters must be explicitly imported**
3. **`OnInit` deprecated in favor of `OnConnect`**

### Step-by-Step Migration

#### 1. Update Import Paths

**Before (v1):**
```go
import "github.com/usace/goquery"
```

**After (v3):**
```go
import "github.com/usace/goquery/v3"
```

#### 2. Import Database Adapters

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
    _ "github.com/usace/goquery/v3/adapters/postgres"  // New: adapter import required
    "github.com/usace/goquery/v3"
)
```

**Adapter import paths:**

| Database | Adapter Import |
|----------|----------------|
| PostgreSQL | `_ "github.com/usace/goquery/v3/adapters/postgres"` |
| DuckDB | `_ "github.com/usace/goquery/v3/adapters/duckdb"` |
| SQLite | `_ "github.com/usace/goquery/v3/adapters/sqlite"` |
| Oracle | `_ "github.com/usace/goquery/v3/adapters/oracle"` |

#### 3. Update go.mod

Run these commands:

```bash
# Remove old version
go get github.com/usace/goquery@none

# Add v3
go get github.com/usace/goquery/v3

# Add adapter(s) you need
go get github.com/usace/goquery/v3/adapters/postgres
go get github.com/usace/goquery/v3/adapters/duckdb
go get github.com/usace/goquery/v3/adapters/sqlite
go get github.com/usace/goquery/v3/adapters/oracle

# Clean up
go mod tidy
```

#### 4. Replace OnInit with OnConnect (Oracle users)

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

#### 5. Update SQLite Driver Names (if applicable)

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
    DbDriver: "sqlite3",  // CGO driver (still supported)
    DbStore:  "sqlx",
    Dbname:   "./mydb.db",
}
```

### Complete Migration Example

**Before (v1):**

```go
package main

import (
    "log"
    
    _ "github.com/jackc/pgx/v4/stdlib"
    "github.com/usace/goquery"
)

func main() {
    config := goquery.RdbmsConfig{
        Dbuser:   "postgres",
        Dbpass:   "password",
        Dbhost:   "localhost",
        Dbport:   "5432",
        Dbname:   "mydb",
        DbDriver: "pgx",
        DbStore:  "pgx",
    }
    
    store, err := goquery.NewRdbmsDataStore(&config)
    if err != nil {
        log.Fatal(err)
    }
    
    var users []User
    err = store.Select("SELECT * FROM users").Dest(&users).Fetch()
    if err != nil {
        log.Fatal(err)
    }
}
```

**After (v3):**

```go
package main

import (
    "log"
    
    _ "github.com/jackc/pgx/v4/stdlib"
    _ "github.com/usace/goquery/v3/adapters/postgres"  // NEW: adapter import
    "github.com/usace/goquery/v3"                       // CHANGED: /v3 suffix
)

func main() {
    config := goquery.RdbmsConfig{
        Dbuser:   "postgres",
        Dbpass:   "password",
        Dbhost:   "localhost",
        Dbport:   "5432",
        Dbname:   "mydb",
        DbDriver: "pgx",
        DbStore:  "pgx",
        // OPTIONAL: Add initialization hook
        OnConnect: func(db goquery.DataStore) error {
            log.Println("Connected to PostgreSQL")
            return nil
        },
    }
    
    store, err := goquery.NewRdbmsDataStore(&config)
    if err != nil {
        log.Fatal(err)
    }
    
    var users []User
    err = store.Select("SELECT * FROM users").Dest(&users).Fetch()
    if err != nil {
        log.Fatal(err)
    }
}
```

### New Features to Consider Using

After migrating, consider taking advantage of these new v3 features:

#### 1. OnConnect Hook for Initialization

```go
config := goquery.RdbmsConfig{
    // ... connection details ...
    OnConnect: func(db goquery.DataStore) error {
        log.Println("Database connected")
        // Set session variables
        // Load extensions
        // Create temporary tables
        return nil
    },
}
```

#### 2. DuckDB for Analytics

```go
import (
    _ "github.com/duckdb/duckdb-go/v2"
    _ "github.com/usace/goquery/v3/adapters/duckdb"
)

duckConfig := goquery.RdbmsConfig{
    Dbname:   "analytics.duckdb",
    DbDriver: "duckdb",
    DbStore:  "sqlx",
    OnConnect: func(db goquery.DataStore) error {
        return db.Exec(goquery.NoTx, "INSTALL spatial; LOAD spatial")
    },
}

analyticsStore, _ := goquery.NewRdbmsDataStore(&duckConfig)
```

#### 3. Native Go SQLite (No CGO)

```go
import (
    _ "modernc.org/sqlite"
    _ "github.com/usace/goquery/v3/adapters/sqlite"
)

config := goquery.RdbmsConfig{
    Dbname:   "./cache.db",
    DbDriver: "sqlite",  // Native Go, no CGO needed
    DbStore:  "sqlx",
}
```

#### 4. Driver Connectors for Advanced Configuration

```go
import duckdb "github.com/duckdb/duckdb-go/v2"

connector, _ := duckdb.NewConnector("data.duckdb", func(execer driver.ExecerContext) error {
    _, err := execer.ExecContext(context.Background(), "SET threads=8", nil)
    return err
})

config := goquery.RdbmsConfig{
    DbDriver:  "duckdb",
    DbStore:   "sqlx",
    Connector: connector,
}
```

### Compatibility Matrix

| Feature | v1 | v3 | Notes |
|---------|----|----|-------|
| Import path | `github.com/usace/goquery` | `github.com/usace/goquery/v3` | Breaking |
| Adapter imports | Not required | **Required** | Breaking |
| `OnInit` | ✅ Supported | ⚠️ Deprecated | Use `OnConnect` |
| `OnConnect` | ❌ Not available | ✅ New feature | Recommended |
| `Connector` | ❌ Not available | ✅ New feature | Advanced use |
| DuckDB | ❌ Not supported | ✅ Supported | New adapter |
| Native SQLite | ❌ Not available | ✅ Supported | No CGO |
| CGO SQLite | ✅ Supported | ✅ Supported | Still works |
| PostgreSQL | ✅ Supported | ✅ Supported | Works same |
| Oracle | ✅ Supported | ✅ Supported | Works same |

### Troubleshooting Migration Issues

#### "uninitialized or unsupported driver"

**Error:**
```
uninitialized or unsupported driver 'pgx'.
Make sure you imported the adapter:
import _ "github.com/usace/goquery/v3/adapters/postgres"
```

**Solution:** Add the missing adapter import:
```go
import _ "github.com/usace/goquery/v3/adapters/postgres"
```

#### Import path conflicts

**Error:**
```
cannot use goquery v1 and v3 in the same module
```

**Solution:** Remove v1 completely:
```bash
go get github.com/usace/goquery@none
go mod tidy
```

#### CGO errors after upgrading

If you're using CGO SQLite (`sqlite3`) and encountering build errors:

**Option 1:** Switch to native Go SQLite:
```go
import _ "modernc.org/sqlite"

config.DbDriver = "sqlite"  // Change from "sqlite3"
```

**Option 2:** Ensure you have CGO enabled:
```bash
CGO_ENABLED=1 go build
```

### Gradual Migration Strategy

For large codebases, you can migrate gradually:

1. **Create a compatibility layer:**

```go
// compat/goquery.go
package compat

import (
    _ "github.com/jackc/pgx/v4/stdlib"
    _ "github.com/usace/goquery/v3/adapters/postgres"
    goquery "github.com/usace/goquery/v3"
)

// Re-export commonly used types
type DataStore = goquery.DataStore
type RdbmsConfig = goquery.RdbmsConfig
type FluentSelect = goquery.FluentSelect
type FluentInsert = goquery.FluentInsert

// Re-export functions
var NewRdbmsDataStore = goquery.NewRdbmsDataStore
var RdbmsConfigFromEnv = goquery.RdbmsConfigFromEnv
```

2. **Update imports gradually:**

```go
// Change this:
import "github.com/usace/goquery"

// To this:
import goquery "yourapp/compat"
```

3. **Once all imports use compat package, remove it and update to v3 directly.**

### Getting Help

If you encounter issues during migration:

- **GitHub Issues:** https://github.com/usace/goquery/issues
- **Documentation:** https://github.com/usace/goquery/tree/v3
- **Examples:** https://github.com/usace/goquery/tree/v3/examples

---

## License

MIT License - see LICENSE file for details

---

## Support

- **Issues:** https://github.com/usace/goquery/issues
- **Documentation:** https://github.com/usace/goquery
- **Email:** support@usace.army.mil

---

**Version:** 3.0  
**Last Updated:** 2024  
**Maintained by:** U.S. Army Corps of Engineers


