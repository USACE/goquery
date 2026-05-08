### Pattern 3: Dynamic Query Generation with Loops

Build queries dynamically when you need to use custom functions for complex logic:

```go
executor := goquery.NewParallelExecutor()

// Store results for each permission
results := make(map[string][]Record)
var mu sync.Mutex

// Add queries based on user permissions
for _, permission := range user.Permissions {
    perm := permission  // Capture loop variable
    
    // Use func() error when you need the closure
    executor.Add(func() error {
        var data []Record
        query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", perm.Table)
        err := store.Select(query).Params(user.ID).Dest(&data).Fetch()
        if err != nil {
            return err
        }
        
        // Store results safely
        mu.Lock()
        results[perm.Table] = data
        mu.Unlock()
        
        return nil
    })
}

err := executor.Run(ctx)

// Now results map contains data from all tables
for table, records := range results {
    fmt.Printf("%s: %d records\n", table, len(records))
}
```

**⚠️ Important:** Always capture loop variables when using closures in loops. Use a mutex when writing to shared data structures.