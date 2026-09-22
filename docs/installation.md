# Installation

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
