# Output Formats

### JSON Output

#### Stream to Writer (Recommended for large datasets)

```go
ctx := context.Background()
var buf bytes.Buffer
writer := bufio.NewWriter(&buf)

err := store.Select("SELECT * FROM users").
    Context(ctx).
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
ctx := context.Background()
err := store.Select("SELECT * FROM users").
    OutputJson(writer).
    IsJsonArray(true).  // Wrap results in []
    Fetch()

// Output: [{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]
```

#### Single Object

```go
ctx := context.Background()
err := store.Select("SELECT * FROM users WHERE id = $1").
    Params(1).
    Context(ctx).
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
ctx := context.Background()
jsonBytes, err := store.Select("SELECT * FROM users LIMIT 10").
    Context(ctx).
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
