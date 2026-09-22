# API Reference

### FluentSelect Methods

| Method | Description | Example |
| :--- | :--- | :--- |
| `DataSet(ds DataSet)` | Use a DataSet | `.DataSet(&usersDS)` |
| `StatementKey(key string)` | Use named statement from DataSet | `.StatementKey("get-all")` |
| `Params(params ...interface{})` | Bind parameters | `.Params(123, "john")` |
| `Apply(vals ...interface{})` | Apply formatting (internal only) | `.Apply("users", "id")` |
| `Suffix(suffix string)` | Append to query | `.Suffix("LIMIT 10")` |
| `Dest(dest interface{})` | Set destination | `.Dest(&users)` |
| `Tx(tx *Tx)` | Use transaction | `.Tx(&tx)` |
| `LogSql(log bool)` | Log SQL statement | `.LogSql(true)` |
| `PanicOnErr(panic bool)` | Panic on error | `.PanicOnErr(true)` |
| `OutputJson(w io.Writer)` | Output as JSON | `.OutputJson(writer)` |
| `OutputCsv(w io.Writer)` | Output as CSV | `.OutputCsv(writer)` |
| `IsJsonArray(array bool)` | Wrap JSON in array | `.IsJsonArray(true)` |
| `CamelCase(camel bool)` | Convert to camelCase | `.CamelCase(true)` |
| `OmitNull(omit bool)` | Omit null fields | `.OmitNull(true)` |
| `DateFormat(format string)` | Custom date format | `.DateFormat("2006-01-02")` |
| `ForEachRow(fn RowFunction)` | Iterate rows | `.ForEachRow(func(r Rows) error {...})` |
| `Fetch()` | Execute query | `.Fetch()` |
| `FetchRows()` | Get row cursor | `.FetchRows()` |

### FluentInsert Methods

| Method | Description | Example |
| :--- | :--- | :--- |
| `Records(recs interface{})` | Set records to insert | `.Records(&user)` or `.Records(&users)` |
| `Tx(tx *Tx)` | Use transaction | `.Tx(&tx)` |
| `Batch(batch bool)` | Use batch mode | `.Batch(true)` |
| `BatchSize(size int)` | Set batch size | `.BatchSize(500)` |
| `PanicOnErr(panic bool)` | Panic on error | `.PanicOnErr(true)` |
| `Execute()` | Execute insert | `.Execute()` |

### Transaction Methods

| Method | Description |
| :--- | :--- |
| `Commit()` | Commit transaction |
| `Rollback()` | Rollback transaction |
| `PgxTx()` | Get underlying pgx transaction |
| `SqlXTx()` | Get underlying sqlx transaction |
| `SqlTx()` | Get underlying sql.Tx |

### Rows Methods

| Method | Description |
| :--- | :--- |
| `Next()` | Advance to next row |
| `Scan(...interface{})` | Scan into variables |
| `ScanStruct(interface{})` | Scan into struct |
| `Columns()` | Get column names |
| `ColumnTypes()` | Get column types |
| `Close()` | Close cursor |
