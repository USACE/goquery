# AI-OPTIMIZED REFERENCE: goquery v3 (ULTIMATE)

## SYSTEM CONTEXT
- **Library:** `github.com/usace/goquery/v3`
- **Purpose:** Fluent, type-safe RDBMS abstraction for Go.
- **Core Pattern:** `DataStore` -> `Fluent[Op]` -> `.Fetch()` / `.Execute()`.

## DATABASE CONFIGURATION MAPPING
*AI Note: You MUST import both the driver AND the goquery adapter for any database to work.*

| Database | `DbDriver` | `DbStore` | Adapter Import (Required) | Driver Import (Required) |
| :--- | :--- | :--- | :--- | :--- |
| **PostgreSQL** | `"pgx"` | `"pgx"` | `github.com/usace/goquery/adapters/postgres/v3` | `github.com/jackc/pgx/v4/stdlib` |
| **DuckDB** | `"duckdb"` | `"sqlx"` | `github.com/usace/goquery/adapters/duckdb/v3` | `github.com/duckdb/duckdb-go/v2` |
| **SQLite (Native)** | `"sqlite"` | `"sqlx"` | `github.com/usace/goquery/adapters/sqlite/v3` | `modernc.org/sqlite` |
| **SQLite (CGO)** | `"sqlite3"` | `"sqlx"` | `github.com/usace/goquery/adapters/sqlite/v3` | `github.com/mattn/go-sqlite3` |
| **Oracle** | `"godror"` | `"sqlx"` | `github.com/usace/goquery/adapters/oracle/v3` | `github.com/godror/godror` |

### ⚠️ Connection Priority Rule
If `RdbmsConfig.Connector` is non-nil, **ALL other connection fields** (`Dbhost`, `Dbport`, `Dbname`, etc.) are **IGNORED**. Use the `Connector` for advanced logic like custom TLS/SSL.

---

## STRUCT TAGS & SCHEMA MAPPING
*AI Note: Use these tags to define how Go structs map to database columns and how identity columns behave during INSERTS.*

| Tag | Value | Description |
| :--- | :--- | :--- |
| `db` | `"column_name"` | Maps Go field to specific DB column name. |
| `db` | `"-"` | Tells goquery to **ignore** this field during mapping/inserts. |
| `dbid` | `"SEQUENCE"` | **PostgreSQL:** Marks field as a Primary Key managed by a sequence. |
| `dbid` | `"AUTOINCREMENT"`| **SQLite:** Marks field as an autoincrementing Primary Key. |
| `idsequence`| `"seq_name"` | **PostgreSQL:** Specifies the exact name of the sequence to use. |

**Example Schema-Aware Struct:**
```go
type User struct {
    ID       int32  `db:"id" dbid:"SEQUENCE" idsequence:"users_id_seq"`
    Username string `db:"username"`
    Email    string `db:"email"`
}
```

---

## THE ENTERPRISE PATTERN: DATASETS
*AI Note: For complex apps, do not pass raw SQL to `.Select()`. Use `TableDataSet` to centralize schema logic.*

**Blueprint: Struct $\rightarrow$ DataSet $\rightarrow$ Execution**

```go
// 1. Define Struct
type Product struct {
    ID    int32   `db:"id" dbid:"AUTOINCREMENT"`
    Name  string  `db:"name"`
    Price float64 `db:"price"`
}

// 2. Define DataSet (Centralized Schema)
var productDS = goquery.TableDataSet{
    Name:   "products",
    Schema: "inventory", 
    Statements: goquery.Statements{
        "get-all":   "SELECT * FROM products ORDER BY name",
        "get-by-id": "SELECT * FROM products WHERE id = $1",
    },
    TableFields: Product{}, // CRITICAL: Drives auto-generated INSERTs
}

// 3. Execute using StatementKey
var products []Product
err := store.Select().
    DataSet(&productDS).
    StatementKey("get-all").
    Dest(&products).
    Fetch()
```

---

## CORE API REFERENCE

### 1. SELECT (Read)
**Pattern:** `store.Select(query).Context(ctx).Params(args...).Dest(&target).Fetch()`

| Method | Purpose | Note |
| :--- | :--- | :--- |
| `.Select(stmt)` | Set SQL | Supports `%s` for `.Apply()` |
| `.DataSet(&ds)` | Use `TableDataSet` | Use `.StatementKey("key")` after |
| `.Params(...interface{})` | Bind parameters | **MANDATORY for user input** |
| `.Apply(...interface{})` | Template injection | **IDENTIFIER ONLY** (table/column names) |
| `.Suffix(string)` | Append SQL | e.g., `"LIMIT 10"` |
| `.Dest(&target)` | Set destination | Pointer to struct or slice |
| `.Fetch()` | Execute & Map | Standard usage for slices/structs |
| `.FetchRows()` | Get row cursor | Returns `(sql.Rows, error)`. **MUST Close()** |
| `.LogSql(bool)` | Debug SQL | **Essential for troubleshooting** |

### 2. INSERT (Write)
**Pattern (Auto-Generated):** `store.Insert(&ds).Records(&data).Context(ctx).Execute()`

| Method | Purpose | Note |
| :--- | :--- | :--- |
| `.Records(&target)` | Data source | Single struct or slice of structs |
| `.Batch(true)` | **PGX Batching** | High performance, single roundtrip |
| `.BatchSize(int)` | Batch limit | Recommended `100-500` |

**Pattern (Manual/Returning):** `store.Select(query).DataSet(&ds).StatementKey("key").Params(...).Dest(&target).Fetch()`
*Use this when you need `RETURNING id` or complex logic not covered by auto-generation.*

### 3. TRANSACTIONS
**Pattern (Safe/Automatic):** `store.Transaction(func(tx goquery.Tx) error { ... })`

| Method | Purpose | Note |
| :--- | :--- | :--- |
| `.PanicOnErr(true)` | **Rollback Trigger** | **MUST use inside Transaction()** to ensure rollback on error/panic. |
| `.MustExec(...)` | Panic on error | Useful inside `Transaction()` blocks. |

---

## SECURITY GUARDRAILS

### 1. Injection Prevention (The Golden Rule)
*   **❌ NEVER** use `fmt.Sprintf` or `+` for user data.
*   **❌ NEVER** pass user input to `.Apply()`.
*   **✅ ALWAYS** use `.Params()` for **Values** (e.g., `WHERE id = $1`).
*   **✅ ONLY** use `.Apply()` for **Identifiers** (e.g., `FROM %s`).

### 2. Credential Safety
*   **❌ NEVER** log the `RdbmsConfig` struct directly (it contains `Dbpass`).
*   **✅ ALWAYS** use `RdbmsConfigFromEnv()` or a secret manager.

---

## TROUBLESHOOTING & DEBUGGING

| Error Message | Likely Cause | Fix |
| :--- | :--- | :--- |
| `"uninitialized or unsupported driver"` | Missing Adapter Import | Add `import _ "github.com/usace/goquery/v3/adapters/..."` |
| `"no rows in result set"` | Query returned nothing | Check input or handle `err != nil` |
| `"converting NULL to string..."` | Nullable DB column | Use pointer types (e.g., `*string`) in struct |
| `"too many open connections"` | Connection leak | Ensure `.FetchRows()` is followed by `defer rows.Close()` |

### Quick Debugging
*   **SQL Visibility:** Append `.LogSql(true)` to any fluent chain.
*   **Manual Inspection:** Use `store.(*goquery.RdbmsDataStore).RdbmsDb().InsertStmt(&ds)` to see generated SQL.

## OUTPUT FORMATS

| Method | Purpose | Key Config Options |
| :--- | :--- | :--- |
| `.OutputJson(w)` | Stream to `io.Writer` | `.IsJsonArray(bool)`, `.CamelCase(bool)`, `.OmitNull(bool)` |
| `.FetchCSV()` | Return CSV string | `.CamelCase(bool)`, `.DateFormat(string)` |

**Example: Direct HTTP JSON Response**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    err := store.Select("SELECT * FROM users").
        OutputJson(w).
        IsJsonArray(true).
        CamelCase(true).
        Fetch()
    if err != nil { http.Error(w, err.Error(), 500) }
}
```
