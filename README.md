# Fiber LibSQL Storage

[![Tests LibSQL](https://github.com/syzio/fiber-libsql/actions/workflows/test-libsql.yml/badge.svg)](https://github.com/syzio/fiber-libsql/actions/workflows/test-libsql.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/syzio/fiber-libsql.svg)](https://pkg.go.dev/github.com/syzio/fiber-libsql)
[![Go Report Card](https://goreportcard.com/badge/github.com/syzio/fiber-libsql)](https://goreportcard.com/report/github.com/syzio/fiber-libsql)

LibSQL storage implementation for [Fiber](https://github.com/gofiber/fiber) framework using [go-libsql](https://github.com/tursodatabase/go-libsql).

### Requirements
- Go 1.23+

### Install
```bash
go get github.com/syzio/fiber-libsql
```

### Examples

#### Local Database
```go
storage := libsql.New(libsql.Config{
    Connection: libsql.Local{
        Database: "file:./fiber.db",
    },
})
```

#### Remote Database (Turso)
```go
storage := libsql.New(libsql.Config{
    Connection: libsql.Remote{
        Database:  "libsql://your-database.turso.io",
        AuthToken: "your-token",
    },
})
```

#### Embedded Replica
```go
storage := libsql.New(libsql.Config{
    Connection: libsql.EmbeddedReplica{
        Database:      "file:./local.db",
        PrimaryURL:    "libsql://your-database.turso.io",
        AuthToken:     "your-token",
        EncryptionKey: "optional-key",
        SyncInterval:  time.Second,
    },
})
```

### Configuration

```go
type Config struct {
    Connection Connection  // Required: database connection configuration
    Table      string     // Optional: table name (default: fiber_storage)
    Reset      bool       // Optional: drop table if exists (default: false)
    GCInterval time.Duration // Optional: GC interval (default: 10s)
}
```

### License
MIT License - see [LICENSE](LICENSE) for more details.
