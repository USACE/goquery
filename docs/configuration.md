# Configuration

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
