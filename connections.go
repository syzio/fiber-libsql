package libsql

import (
	"database/sql"
	"fmt"
	"time"

	driver "github.com/tursodatabase/go-libsql"
)

// Connection defines the interface for different database connection types
type Connection interface {
	Db() (*sql.DB, error)
}

// Local represents a local LibSQL database connection
type Local struct {
	Database string
}

func (c Local) Db() (*sql.DB, error) {
	return sql.Open("libsql", c.Database)
}

// Remote represents a remote LibSQL database connection
type Remote struct {
	Database  string
	AuthToken string
}

func (c Remote) Db() (*sql.DB, error) {
	return sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", c.Database, c.AuthToken))
}

// EmbeddedReplica represents a Turso embedded replica connection
type EmbeddedReplica struct {
	Database   string
	PrimaryURL string
	AuthToken  string

	EncryptionKey string

	SyncInterval time.Duration
}

func (c EmbeddedReplica) Db() (*sql.DB, error) {
	var options = []driver.Option{
		driver.WithAuthToken(c.AuthToken),
		driver.WithSyncInterval(c.SyncInterval),
	}

	if c.EncryptionKey != "" {
		options = append(options, driver.WithEncryption(c.EncryptionKey))
	}

	if conn, err := driver.NewEmbeddedReplicaConnector(c.Database, c.PrimaryURL, options...); err != nil {
		return nil, err
	} else {
		return sql.OpenDB(conn), nil
	}
}
