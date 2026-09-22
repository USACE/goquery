# Transactions

### Automatic Transactions

The `Transaction()` method handles commit/rollback automatically:

```go
// All operations within this function are part of the transaction
err := store.Transaction(func(tx goquery.Tx) {
    ctx := context.Background()
    
    // Insert a user
    store.Insert(&usersDS).
        Records(&user).
        Tx(&tx).
        Context(ctx).
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
ctx := context.Background()
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
ctx := context.Background()
err := store.Transaction(func(tx goquery.Tx) {
    // Option 1: PanicOnErr with fluent API
    store.Insert(&usersDS).
        Records(&user).
        Tx(&tx).
        Context(ctx).
        PanicOnErr(true).
        Execute()
    
    // Option 2: MustExec (panics on error)
    store.MustExec(&tx, "UPDATE accounts SET balance = balance + $1 WHERE user_id = $2", 0, user.ID)
    
    // Option 3: Manual panic
    result, err := store.Execr(&tx, "DELETE FROM temp_data")
    if err != nil {
        panic(err)  // Triggers rollback
    }
})
```
