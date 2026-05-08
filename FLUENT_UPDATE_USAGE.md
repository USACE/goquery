# FluentUpdate - Consistent with goquery Fluent API

## Overview

`FluentUpdate` provides a chainable, fluent API for UPDATE operations that follows the exact same pattern as `FluentSelect` and `FluentInsert`. It uses DataSets with predefined statements and parameterized queries.

## Pattern Consistency

FluentUpdate follows the established goquery pattern:
- **DataSet** - Container with named SQL statements
- **StatementKey** - Reference to a statement in the DataSet
- **Params()** - Bind parameters for the statement
- **Apply()** - Template values (use with caution)

## Basic Usage

### Step 1: Define a DataSet with UPDATE statements

```go
var usersDS = goquery.TableDataSet{
    Name: "users",
    Schema: "public",
    Statements: goquery.Statements{
        "update-status":     "UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2",
        "update-login":      "UPDATE users SET last_login = $1, login_count = login_count + 1 WHERE id = $2",
        "deactivate-old":    "UPDATE users SET status = 'inactive' WHERE last_login < $1",
        "update-by-email":   "UPDATE users SET %s WHERE email = $1",
    },
    TableFields: User{},
}
```

### Step 2: Execute updates using the fluent API

```go
// Simple update with parameters
err := store.Update().
    DataSet(&usersDS).
    StatementKey("update-status").
    Params("active", 123).
    Execute()

// With transaction
tx, _ := store.NewTransaction()
err := store.Update().
    DataSet(&usersDS).
    StatementKey("update-login").
    Params(time.Now(), userID).
    Tx(&tx).
    Execute()
tx.Commit()
```

## Complete Examples

### Example 1: Basic Update

```go
package main

import (
    "github.com/usace/goquery/v3"
    "log"
    "time"
)

var usersDS = goquery.TableDataSet{
    Name: "users",
    Statements: goquery.Statements{
        "update-login": "UPDATE users SET last_login = $1 WHERE id = $2",
    },
    TableFields: User{},
}

func main() {
    var store goquery.DataStore
    
    err := store.Update().
        DataSet(&usersDS).
        StatementKey("update-login").
        Params(time.Now(), 123).
        Execute()
    
    if err != nil {
        log.Fatal(err)
    }
}
```

### Example 2: Update with Rows Affected

```go
var productsDS = goquery.TableDataSet{
    Name: "products",
    Statements: goquery.Statements{
        "discount": "UPDATE products SET price = price * $1 WHERE category = $2",
    },
    TableFields: Product{},
}

result, err := store.Update().
    DataSet(&productsDS).
    StatementKey("discount").
    Params(0.9, "electronics").  // 10% discount
    Execr()

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Discounted %d products\n", result.RowsAffected())
```

### Example 3: Using Apply() for Dynamic Columns

```go
var usersDS = goquery.TableDataSet{
    Name: "users",
    Statements: goquery.Statements{
        // %s placeholder for Apply()
        "update-field": "UPDATE users SET %s = $1 WHERE id = $2",
    },
    TableFields: User{},
}

// Update different fields dynamically
err := store.Update().
    DataSet(&usersDS).
    StatementKey("update-field").
    Apply("email").  // Fills in the %s
    Params("newemail@example.com", 123).
    Execute()

// ⚠️ WARNING: Never use Apply() with user input (SQL injection risk)
```

### Example 4: With Transaction and Panic Mode

```go
var ordersDS = goquery.TableDataSet{
    Name: "orders",
    Statements: goquery.Statements{
        "ship": "UPDATE orders SET status = 'shipped', shipped_at = $1 WHERE id = $2",
    },
    TableFields: Order{},
}

var inventoryDS = goquery.TableDataSet{
    Name: "inventory",
    Statements: goquery.Statements{
        "decrement": "UPDATE inventory SET quantity = quantity - $1 WHERE product_id = $2",
    },
    TableFields: Inventory{},
}

err := store.Transaction(func(tx goquery.Tx) {
    // Update order status
    store.Update().
        DataSet(&ordersDS).
        StatementKey("ship").
        Params(time.Now(), orderID).
        Tx(&tx).
        PanicOnErr(true).
        Execute()
    
    // Update inventory
    store.Update().
        DataSet(&inventoryDS).
        StatementKey("decrement").
        Params(quantity, productID).
        Tx(&tx).
        PanicOnErr(true).
        Execute()
})

if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

### Example 5: With Suffix (RETURNING clause)

```go
var usersDS = goquery.TableDataSet{
    Name: "users",
    Statements: goquery.Statements{
        "increment-version": "UPDATE users SET version = version + 1 WHERE id = $1",
    },
    TableFields: User{},
}

var newVersion int
err := store.Select().  // Use Select for RETURNING
    DataSet(&usersDS).
    StatementKey("increment-version").
    Params(123).
    Suffix("RETURNING version").
    Dest(&newVersion).
    Fetch()

fmt.Printf("New version: %d\n", newVersion)
```

### Example 6: With ParallelExecutor

```go
var usersDS = goquery.TableDataSet{
    Name: "users",
    Statements: goquery.Statements{
        "update-last-seen": "UPDATE users SET last_seen = $1 WHERE id = $2",
    },
    TableFields: User{},
}

var statsDS = goquery.TableDataSet{
    Name: "stats",
    Statements: goquery.Statements{
        "increment-views": "UPDATE stats SET views = views + 1 WHERE page = $1",
    },
    TableFields: Stats{},
}

executor := goquery.NewParallelExecutor().
    Add(
        store.Update().
            DataSet(&usersDS).
            StatementKey("update-last-seen").
            Params(time.Now(), userID),
        
        store.Update().
            DataSet(&statsDS).
            StatementKey("increment-views").
            Params("/dashboard"),
    )

ctx := context.Background()
if err := executor.Run(ctx); err != nil {
    log.Fatal(err)
}
```

### Example 7: Bulk Update with Dynamic Filters

```go
var usersDS = goquery.TableDataSet{
    Name: "users",
    Statements: goquery.Statements{
        "bulk-update": "UPDATE users SET status = $1 WHERE %s",
    },
    TableFields: User{},
}

// Update users based on different criteria
whereClause := "last_login < $2 AND role = $3"

err := store.Update().
    DataSet(&usersDS).
    StatementKey("bulk-update").
    Apply(whereClause).
    Params("inactive", "2023-01-01", "user").
    LogSql(true).
    Execute()
```

## API Reference

### DataSet(ds DataSet) *FluentUpdate
Sets the DataSet containing UPDATE statements.

### StatementKey(key string) *FluentUpdate
Specifies which statement to use from the DataSet's Statements map.

### Params(params ...interface{}) *FluentUpdate
Sets the bind parameters ($1, $2, etc.) for the UPDATE statement.

### Apply(vals ...interface{}) *FluentUpdate
Adds values to be inserted into the statement using fmt.Sprintf (%s, %d, etc.).
**⚠️ WARNING:** Only use with internal/trusted values. Never use with user input (SQL injection risk).

### Suffix(suffix string) *FluentUpdate
Adds additional SQL to append to the statement (e.g., "RETURNING id").

### Tx(tx *Tx) *FluentUpdate
Sets the transaction to use for this update.

### PanicOnErr(panicOnErr bool) *FluentUpdate
Enables panicking on errors instead of returning them. Useful in transaction blocks.

### LogSql(logSql bool) *FluentUpdate
Enables SQL statement logging for debugging.

### Execute() error
Executes the UPDATE statement and returns an error if one occurs.

### Execr() (ExecResult, error)
Executes the UPDATE statement and returns the result with rows affected count.

## Comparison with FluentSelect

The pattern is identical:

```go
// FluentSelect
var users []User
err := store.Select().
    DataSet(&usersDS).
    StatementKey("get-by-status").
    Params("active").
    Dest(&users).
    Fetch()

// FluentUpdate
err := store.Update().
    DataSet(&usersDS).
    StatementKey("update-status").
    Params("inactive", 123).
    Execute()
```

## Best Practices

### 1. Define UPDATE statements in DataSets

```go
// ✅ GOOD - Organized and reusable
var usersDS = goquery.TableDataSet{
    Name: "users",
    Statements: goquery.Statements{
        "activate":   "UPDATE users SET status = 'active' WHERE id = $1",
        "deactivate": "UPDATE users SET status = 'inactive' WHERE id = $1",
        "update-email": "UPDATE users SET email = $1 WHERE id = $2",
    },
    TableFields: User{},
}

// ❌ BAD - Hard to maintain
err := store.Exec(goquery.NoTx, "UPDATE users SET status = 'active' WHERE id = $1", id)
```

### 2. Use Params() for all data values

```go
// ✅ SAFE
store.Update().
    DataSet(&usersDS).
    StatementKey("update-email").
    Params(email, userID).
    Execute()

// ❌ DANGEROUS - SQL injection risk
store.Update().
    DataSet(&usersDS).
    StatementKey("update-email").
    Apply(email, userID).  // NEVER DO THIS
    Execute()
```

### 3. Check rows affected for critical updates

```go
result, err := store.Update().
    DataSet(&ordersDS).
    StatementKey("cancel").
    Params(orderID).
    Execr()

if err != nil {
    return err
}

if result.RowsAffected() == 0 {
    return fmt.Errorf("order %d not found", orderID)
}
```

### 4. Use PanicOnErr in transactions

```go
store.Transaction(func(tx goquery.Tx) {
    store.Update().
        DataSet(&ds).
        StatementKey("update1").
        Params(val1).
        Tx(&tx).
        PanicOnErr(true).  // Panic triggers rollback
        Execute()
    
    store.Update().
        DataSet(&ds).
        StatementKey("update2").
        Params(val2).
        Tx(&tx).
        PanicOnErr(true).
        Execute()
})
```

### 5. Enable LogSql during development

```go
err := store.Update().
    DataSet(&usersDS).
    StatementKey("update-login").
    Params(time.Now(), 123).
    LogSql(true).  // See generated SQL
    Execute()
```

## Migration from Previous API

If you were using the old non-DataSet API, here's how to migrate:

### Old (Non-standard)
```go
err := store.Update().
    Table("users").
    Set("name", "John").
    Where("id = $1", 123).
    Execute()
```

### New (Standard goquery pattern)
```go
var usersDS = goquery.TableDataSet{
    Name: "users",
    Statements: goquery.Statements{
        "update-name": "UPDATE users SET name = $1 WHERE id = $2",
    },
    TableFields: User{},
}

err := store.Update().
    DataSet(&usersDS).
    StatementKey("update-name").
    Params("John", 123).
    Execute()
```

## When to Use FluentUpdate vs store.Exec()

### Use FluentUpdate when:
- ✅ Updates are part of your core domain logic
- ✅ You want organized, reusable UPDATE statements
- ✅ You need to use with ParallelExecutor
- ✅ You want consistent API with Select/Insert

### Use store.Exec() when:
- ✅ One-off administrative updates
- ✅ Complex SQL with multiple JOINs
- ✅ Dynamic SQL that can't be templated
- ✅ Performance-critical code where builder overhead matters

## Integration Note

To enable the `store.Update()` method, add these to the codebase:

**In `datastore.go` interface:**
```go
Update() *FluentUpdate
```

**In `rdbms_datastore.go` implementation:**
```go
func (sds *RdbmsDataStore) Update() *FluentUpdate {
    fu := FluentUpdate{
        store: sds,
    }
    return &fu
}
```

This matches the existing pattern for `Select()` and `Insert()`.
