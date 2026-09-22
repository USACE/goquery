# Core Concepts

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
