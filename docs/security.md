# Security Best Practices

### 1. Always Use Parameterized Queries

**✅ CORRECT:**

```go
ctx := context.Background()
userID := getUserInput()
err := store.Select("SELECT * FROM users WHERE id = $1").
    Params(userID).  // Safe - uses parameter binding
    Context(ctx).
    Dest(&user).
    Fetch()
```

**❌ WRONG:**

```go
ctx := context.Background()
userID := getUserInput()
sql := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID)  // SQL INJECTION!
err := store.Select(sql).Context(ctx).Dest(&user).Fetch()
```

### 2. Apply() is for Internal Use Only

**✅ CORRECT:**

```go
ctx := context.Background()
// Internal constants only
tableName := "users"
columnName := "id"
userID := getUserInput()

err := store.Select("SELECT * FROM %s WHERE %s = $1").
    Apply(tableName, columnName).  // Internal values
    Params(userID).                // User input
    Context(ctx).
    Dest(&user).
    Fetch()
```

**❌ WRONG:**

```go
ctx := context.Background()
// NEVER pass user input to Apply()
tableName := getUserInput()
err := store.Select("SELECT * FROM %s").
    Apply(tableName).  // SQL INJECTION!
    Dest(&results).
    Fetch()
```

### 3. SSL/TLS Configuration

**Development:**

```go
config := goquery.RdbmsConfig{
    // ... other settings ...
    DbDriverSettings: "sslmode=disable",  // OK for local development
}
```

**Production (Recommended):**

```go
config := goquery.RdbmsConfig{
    // ... other settings ...
    DbDriverSettings: "sslmode=verify-full sslrootcert=/path/to/ca.crt",
}
```

**Available SSL modes (PostgreSQL):**

| Mode | Encryption | Certificate Check | Use Case |
|------|------------|-------------------|----------|
| `disable` | ❌ No | ❌ No | Local dev only |
| `require` | ✅ Yes | ❌ No | Default, basic security |
| `verify-ca` | ✅ Yes | ✅ CA only | Verify server identity |
| `verify-full` | ✅ Yes | ✅ CA + hostname | **Production (recommended)** |

### 4. Credential Management

**❌ DON'T** hardcode credentials:

```go
ctx := context.Background()
config := goquery.RdbmsConfig{
    Dbuser: "admin",
    Dbpass: "password123",  // NEVER hardcode passwords!
}
```

**✅ DO** use environment variables:

```go
ctx := context.Background()
config := goquery.RdbmsConfigFromEnv()
// Or:
config := goquery.RdbmsConfig{
    Dbuser: os.Getenv("DB_USER"),
    Dbpass: os.Getenv("DB_PASS"),
    // ...
}
```

**✅ BETTER** - Use secret managers:

```go
// Example with AWS Secrets Manager
ctx := context.Background()
secret := getSecretFromAWS("prod/db/credentials")
config := goquery.RdbmsConfig{
    Dbuser: secret["username"],
    Dbpass: secret["password"],
    // ...
}
```

### 5. Never Log Connection Strings

Connection strings contain passwords in plaintext. Never log them:

```go
// ❌ WRONG
log.Printf("Connecting with config: %+v", config)  // Logs password!

// ✅ CORRECT
log.Printf("Connecting to %s@%s:%s/%s", config.Dbuser, config.Dbhost, config.Dbport, config.Dbname)
```

### 6. Input Validation

Always validate user input before using it in queries:

```go
ctx := context.Background()
func getUser(idStr string) (*User, error) {
    // Validate input
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }
    if id <= 0 {
        return nil, errors.New("user ID must be positive")
    }
    
    // Safe to use in query
    var user User
    err = store.Select("SELECT * FROM users WHERE id = $1").
        Params(id).
        Context(ctx).
        Dest(&user).
        Fetch()
    
    return &user, err
}
```

### 7. Least Privilege Principle

Use database users with minimal required permissions:

```sql
-- Create application user with limited permissions
CREATE USER app_user WITH PASSWORD 'secure_password';

-- Grant only what's needed
GRANT SELECT, INSERT, UPDATE ON users TO app_user;
GRANT SELECT ON products TO app_user;

-- Don't grant DELETE or DROP permissions unless necessary
```
