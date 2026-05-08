# FluentUpdate Implementation Summary

## ✅ Implementation Complete

`fluent_update.go` has been implemented following the exact same pattern as `FluentSelect` and `FluentInsert`.

## Key Design Principles

### 1. DataSet-Based (Like FluentSelect)
UPDATE statements are defined in TableDataSet's Statements map, not built dynamically.

```go
var usersDS = goquery.TableDataSet{
    Name: "users",
    Statements: goquery.Statements{
        "update-status": "UPDATE users SET status = $1 WHERE id = $2",
        "update-email":  "UPDATE users SET email = $1, updated_at = NOW() WHERE id = $2",
    },
    TableFields: User{},
}
```

### 2. StatementKey-Based (Like FluentSelect)
Reference statements by key, not by building SQL.

```go
err := store.Update().
    DataSet(&usersDS).
    StatementKey("update-status").
    Params("active", 123).
    Execute()
```

### 3. Params() for Values (Like FluentSelect)
All data values go through `.Params()`, not through builder methods.

```go
// ✅ CORRECT - Follows goquery pattern
store.Update().
    DataSet(&usersDS).
    StatementKey("update-email").
    Params("newemail@example.com", userID).
    Execute()

// ❌ WRONG - Not the goquery way
store.Update().
    Table("users").
    Set("email", "newemail@example.com").
    Where("id = $1", userID).
    Execute()
```

## API Consistency

All three fluent APIs follow the same pattern:

| FluentSelect | FluentInsert | FluentUpdate |
|--------------|--------------|--------------|
| `DataSet()` | `DataSet()` | `DataSet()` |
| `StatementKey()` | (not used) | `StatementKey()` |
| `Params()` | `Records()` | `Params()` |
| `Apply()` | (not used) | `Apply()` |
| `Suffix()` | (not used) | `Suffix()` |
| `Tx()` | `Tx()` | `Tx()` |
| `PanicOnErr()` | `PanicOnErr()` | `PanicOnErr()` |
| `LogSql()` | (not used) | `LogSql()` |
| `Fetch()` | `Execute()` | `Execute()` |
| (not used) | (not used) | `Execr()` |

## Files Created

1. **`fluent_update.go`** - Implementation (100% consistent with existing patterns)
2. **`FLUENT_UPDATE_USAGE.md`** - Complete usage guide with examples
3. **`FLUENT_UPDATE_NOTES.md`** - Implementation notes (from earlier, now outdated)

## Example Usage

```go
// Define DataSet with UPDATE statements
var usersDS = goquery.TableDataSet{
    Name: "users",
    Statements: goquery.Statements{
        "update-status": "UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2",
        "update-login":  "UPDATE users SET last_login = $1, login_count = login_count + 1 WHERE id = $2",
    },
    TableFields: User{},
}

// Execute update
err := store.Update().
    DataSet(&usersDS).
    StatementKey("update-status").
    Params("active", 123).
    Execute()

// With rows affected
result, err := store.Update().
    DataSet(&usersDS).
    StatementKey("update-login").
    Params(time.Now(), userID).
    Execr()

fmt.Printf("Updated %d rows\n", result.RowsAffected())
```

## Works With ParallelExecutor

```go
executor := goquery.NewParallelExecutor().
    Add(
        store.Update().
            DataSet(&usersDS).
            StatementKey("update-status").
            Params("active", 123),
        
        store.Update().
            DataSet(&ordersDS).
            StatementKey("ship").
            Params(time.Now(), orderID),
    )

ctx := context.Background()
executor.Run(ctx)
```

## Implementation Notes

### What's Implemented
- ✅ Follows exact same pattern as FluentSelect
- ✅ Uses DataSet and StatementKey
- ✅ Params() for bind parameters
- ✅ Apply() for sprintf templates (with warning)
- ✅ Suffix() for RETURNING clauses
- ✅ Tx() for transactions
- ✅ PanicOnErr() for transaction safety
- ✅ LogSql() for debugging
- ✅ Execute() and Execr() methods
- ✅ Compatible with ParallelExecutor
- ✅ Compiles successfully

### What's Not in fluent_update.go
- ⏸️ `store.Update()` method (needs DataStore interface change)
- ⏸️ Tests (should be added separately)

### To Enable Full Integration

Add to DataStore interface and implementations:

```go
// In datastore.go
type DataStore interface {
    // ...
    Update() *FluentUpdate
}

// In rdbms_datastore.go
func (sds *RdbmsDataStore) Update() *FluentUpdate {
    return &FluentUpdate{store: sds}
}
```

## Why This Design?

This matches goquery's philosophy:
1. **SQL statements live in DataSets** - Organized, versioned, testable
2. **No dynamic SQL building** - Security through simplicity
3. **Consistent API** - Same pattern across Select/Insert/Update
4. **Explicit is better than implicit** - No magic, clear intent

## Comparison

### Other ORMs (Not goquery's way)
```go
db.Model(&User{}).Where("id = ?", 123).Update("status", "active")
```

### goquery's Way
```go
var usersDS = goquery.TableDataSet{
    Statements: goquery.Statements{
        "update-status": "UPDATE users SET status = $1 WHERE id = $2",
    },
}

store.Update().DataSet(&usersDS).StatementKey("update-status").Params("active", 123).Execute()
```

Benefits of goquery's approach:
- SQL is visible and maintainable
- Can be tested independently
- Can be reviewed by DBAs
- No surprises in generated SQL
- Works with complex SQL (not limited to simple updates)
