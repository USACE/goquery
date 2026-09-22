# goquery - Comprehensive Documentation

**Version:** 3.0  
**License:** MIT  
**Go Version:** 1.24+

---

## What's New in v3

Version 3.0 introduces significant architectural improvements and new database support:

- 🆕 **Module Versioning** - Import path is now `github.com/usace/goquery/v3`
- 🔌 **OnConnect Hook** - Execute initialization code when connections are established
- 🦆 **DuckDB Support** - Full support for DuckDB with spatial extensions
- 🗄️ **Dual SQLite Modes** - Choose between native Go (`sqlite`) or CGO (`sqlite3`) drivers
- 🔗 **Driver Connectors** - Direct `driver.Connector` support for advanced connection management
- 📦 **Modular Adapters** - Database adapters are now separate modules to reduce dependencies

See the [v3 Migration Guide](docs/migration-guide.md) for upgrade instructions.

---

## Table of Contents

1. [Overview](docs/overview.md)
2. [Installation](docs/installation.md)
3. [Quick Start](docs/quickstart.md)
4. [Configuration](docs/configuration.md)
5. [Core Concepts](docs/core-concepts.md)
6. [DataStore Operations](docs/operations.md)
7. [Transactions](docs/transactions.md)
8. [Batch Operations](docs/batching.md)
9. [Output Formats](docs/output-formats.md)
10. [Security Best Practices](docs/security.md)
11. [Advanced Usage](docs/advanced-usage.md)
12. [Troubleshooting](docs/troubleshooting.md)
13. [API Reference](docs/api-reference.md)
14. [v3 Migration Guide](docs/migration-guide.md)

---

## License
MIT License - see LICENSE file for details

## Support

- **Issues:** [https://github.com/usace/goquery/issues](https://github.com/usace/goquery/issues)
- **Documentation:** [https://github.com/usace/goquery](https://github.com/usace/goquery)
- **Email:** support@usace.army.mil

**Version:** 3.0  
**Last Updated:** 2024  
**Maintained by:** U.S. Army Corps of Engineers
