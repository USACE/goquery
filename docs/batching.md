# Batch Operations

Batch operations dramatically improve performance for bulk inserts by reducing network round trips.

### When to Use Batches

- ✅ Inserting 100+ records
- ✅ High-throughput data pipelines
- ✅ ETL operations
- ✅ Importing large datasets

### Batch Insert Example

```go
// Generate test data
users := make([]User, 10000)
for i := range users {
    users[i] = User{
        Username: fmt.Sprintf("user_%d", i),
        Email:    fmt.Sprintf("user_%d@example.com", i),
    }
}

// Batch insert with error handling
ctx := context.Background()
// ... populate users ...

err := store.Insert(&usersDS).
    Records(&users).
    Context(ctx).
    Batch(true).
    BatchSize(500).  // 500 records per batch
    Execute()

if err != nil {
    log.Fatalf("Batch insert failed: %v", err)
}
```

### Batch Size Tuning

Choose batch size based on your data and network:

| Batch Size | Use Case | Notes |
|------------|----------|-------|
| 50-100 | Small records, slow network | Reduces packet size |
| 100-500 | Typical web applications | Good default |
| 500-1000 | Large records, fast network | Maximum throughput |
| 1000+ | Very large datasets, LAN | Diminishing returns |

### Performance Comparison

Example with 10,000 records:

| Method | Time | Network Trips |
|--------|------|---------------|
| Individual inserts | ~45s | 10,000 |
| Bulk insert (multi-tx) | ~25s | 10,000 |
| Batch insert (size 100) | ~2s | 100 |
| Batch insert (size 500) | ~1s | 20 |

### Batch Error Handling

Batch operations validate **each statement** in the batch:

```go
err := store.Insert(&usersDS).
    Records(&users).
    Batch(true).
    BatchSize(100).
    Execute()

if err != nil {
    // Error message includes which record failed
    // Example: "batch insert failed at record 342: duplicate key value"
    log.Printf("Batch failed: %v\n", err)
}
```

**Important:** If any record in a batch fails, the entire batch is rolled back (not the entire operation).
