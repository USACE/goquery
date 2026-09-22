# Overview

**goquery** is a lightweight Go library that simplifies database operations through a fluent, type-safe API. It provides a unified interface for multiple database backends while maintaining performance and safety.

### Key Features

- ✅ **Fluent API** - Chainable, readable query building
- ✅ **Multi-Database Support** - PostgreSQL (pgx), SQLite (native/CGO), Oracle, DuckDB
- ✅ **Type-Safe Mapping** - Automatic struct-to-row mapping via tags
- ✅ **Transaction Support** - Automatic rollback on panic, commit on success
- ✅ **Batch Operations** - High-performance bulk inserts (pgx)
- ✅ **Multiple Output Formats** - Structs, JSON, CSV
- ✅ **Connection Pooling** - Configurable pool settings
- ✅ **SQL Generation** - Auto-generate INSERT/SELECT from structs
- ✅ **Security First** - Parameterized queries prevent SQL injection
- ✅ **OnConnect Hooks** - Initialize connections with custom logic
- ✅ **Modular Architecture** - Import only the database drivers you need

### Architecture

```
┌─────────────────┐
│   Your Code     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   DataStore     │  ← Unified interface
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌───────┐  ┌───────┐
│  PgxDb│  │SqlxDb │  ← Driver implementations
└───┬───┘  └───┬───┘
    │          │
    ▼          ▼
┌───────┐  ┌───────┐
│  pgx  │  │ sqlx  │  ← Underlying libraries
└───────┘  └───────┘
```
