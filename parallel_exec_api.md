# Parallel Query Executor - Simplified API

## Summary

The `ParallelExecutor.Add()` method now accepts multiple types directly without needing `func() error` wrappers:

✅ **FluentSelect** - Pass directly  
✅ **FluentInsert** - Pass directly  
✅ **DeferredExec** - For UPDATE/DELETE statements  
✅ **func() error** - Still supported for custom logic  
✅ **Executable interface** - Custom types  

## Quick Examples

### Before (Verbose)
```go
var users []User

executor := goquery.NewParallelExecutor().
    Add(func() error {
        return store.Select("SELECT * FROM users").Dest(&users).Fetch()
    })
```

### After (Clean)
```go
var users []User

executor := goquery.NewParallelExecutor().
    Add(store.Select("SELECT * FROM users").Dest(&users))
```

## Complete Example

```go
package main

import (
    "context"
    "github.com/usace/goquery/v3"
)

func main() {
    var store goquery.DataStore
    
    // Declare results outside
    var users []User
    var products []Product
    var newOrder Order
    
    executor := goquery.NewParallelExecutor().
        MaxConcurrency(5).
        // SELECT queries - pass FluentSelect directly
        Add(
            store.Select("SELECT * FROM users WHERE active = $1").
                Params(true).
                Dest(&users),
        ).
        Add(
            store.Select("SELECT * FROM products").
                Dest(&products),
        ).
        // INSERT - pass FluentInsert directly
        Add(
            store.Insert(&ordersDS).Records(&newOrder),
        ).
        // UPDATE/DELETE - use DeferredExec
        Add(
            goquery.NewDeferredExec(store, goquery.NoTx,
                "UPDATE stats SET views = views + 1"),
        ).
        // Custom logic - use func() error
        Add(func() error {
            // Complex logic here
            return customOperation()
        })
    
    ctx := context.Background()
    if err := executor.Run(ctx); err != nil {
        log.Fatal(err)
    }
    
    // Access results
    fmt.Printf("Loaded %d users, %d products\n", len(users), len(products))
}
```

## Supported Types

### 1. FluentSelect (Most Common)
```go
var data []Record

executor.Add(
    store.Select("SELECT * FROM table").Dest(&data),
)
// Automatically calls .Fetch()
```

### 2. FluentInsert
```go
executor.Add(
    store.Insert(&dataset).Records(&record),
)
// Automatically calls .Execute()
```

### 3. DeferredExec (For Direct SQL)
```go
executor.Add(
    goquery.NewDeferredExec(store, goquery.NoTx,
        "UPDATE users SET last_login = NOW() WHERE id = $1", userID),
    goquery.NewDeferredExec(store, goquery.NoTx,
        "DELETE FROM sessions WHERE expired = true"),
)
```

**Why DeferredExec?**

`store.Exec()` executes immediately, so you can't pass it to `Add()`. `DeferredExec` wraps it for deferred execution.

### 4. Custom Functions (For Complex Logic)
```go
executor.Add(func() error {
    log.Println("Starting operation...")
    var data []Record
    err := store.Select("SELECT * FROM records").Dest(&data).Fetch()
    if err != nil {
        return err
    }
    return processData(data)
})
```

### 5. Executable Interface (Advanced)
```go
type Executable interface {
    Execute() error
}

type MyQuery struct {
    store DataStore
}

func (q *MyQuery) Execute() error {
    return q.store.Exec(goquery.NoTx, "UPDATE ...")
}

executor.Add(&MyQuery{store: store})
```

## Key Benefits

1. **Less Boilerplate** - No need to wrap everything in `func() error`
2. **More Readable** - Code intent is clearer
3. **Flexible** - Mix and match different types
4. **Backward Compatible** - `func() error` still works
5. **Type Safe** - Panics on unsupported types

## Important Notes

### Declare Variables Outside
```go
// ✅ CORRECT
var users []User
executor.Add(store.Select("SELECT * FROM users").Dest(&users))
// users is accessible here

// ❌ WRONG
executor.Add(func() error {
    var users []User  // Lost after function returns!
    return store.Select("SELECT * FROM users").Dest(&users).Fetch()
})
// users is NOT accessible here
```

### Loop Variable Capture
```go
for _, id := range ids {
    id := id  // Capture loop variable!
    executor.Add(func() error {
        return fetchDataForID(id)
    })
}
```

## API Reference

### NewDeferredExec
```go
func NewDeferredExec(store DataStore, tx *Tx, stmt string, params ...interface{}) *DeferredExec
```

Creates a deferred Exec call.

**Parameters:**
- `store` - The DataStore to execute against
- `tx` - Transaction (use `goquery.NoTx` for no transaction)
- `stmt` - SQL statement
- `params` - Statement parameters

**Returns:**
- `*DeferredExec` - A deferred execution wrapper

**Example:**
```go
exec := goquery.NewDeferredExec(
    store,
    goquery.NoTx,
    "UPDATE users SET active = $1 WHERE id = $2",
    false,
    123,
)

executor.Add(exec)
```

## Complete Comparison

### Old Way (Still Works)
```go
var users []User
var products []Product

executor := goquery.NewParallelExecutor().
    Add(func() error {
        return store.Select("SELECT * FROM users").Dest(&users).Fetch()
    }).
    Add(func() error {
        return store.Select("SELECT * FROM products").Dest(&products).Fetch()
    }).
    Add(func() error {
        return store.Exec(goquery.NoTx, "UPDATE stats SET count = count + 1")
    })

ctx := context.Background()
executor.Run(ctx)
```

### New Way (Recommended)
```go
var users []User
var products []Product

executor := goquery.NewParallelExecutor().
    Add(
        store.Select("SELECT * FROM users").Dest(&users),
        store.Select("SELECT * FROM products").Dest(&products),
        goquery.NewDeferredExec(store, goquery.NoTx, 
            "UPDATE stats SET count = count + 1"),
    )

ctx := context.Background()
executor.Run(ctx)
```

## Additional Resources

- **Full Documentation**: See original `parallel_exec.md` for complete details
- **goquery Basics**: See `readme.md` for DataStore operations and query patterns
- **Examples**: See `parallel_exec_example.go` for more usage patterns

## Implementation Notes

The `Add()` method uses type switching to handle different input types:

```go
func (p *ParallelExecutor) Add(queries ...interface{}) *ParallelExecutor {
    for _, q := range queries {
        switch v := q.(type) {
        case func() error:
            p.queries = append(p.queries, v)
        case *FluentSelect:
            p.queries = append(p.queries, func() error { return v.Fetch() })
        case *FluentInsert:
            p.queries = append(p.queries, func() error { return v.Execute() })
        case *DeferredExec:
            p.queries = append(p.queries, func() error { return v.Execute() })
        case Executable:
            p.queries = append(p.queries, func() error { return v.Execute() })
        default:
            panic(fmt.Sprintf("unsupported type: %T", q))
        }
    }
    return p
}
```

All types are internally converted to `func() error` for uniform execution.
