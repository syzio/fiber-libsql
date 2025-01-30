package libsql

import "time"

// Config defines the config for storage
type Config struct {
	Connection Connection

	Table      string
	Reset      bool
	GCInterval time.Duration
}

var ConfigDefault = Config{
	Connection: Local{
		Database: "file:./fiber.db",
	},

	Table:      "fiber_storage",
	Reset:      false,
	GCInterval: 10 * time.Second,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Connection == nil {
		cfg.Connection = ConfigDefault.Connection
	}

	switch conn := cfg.Connection.(type) {
	case Local:
		if conn.Database == "" {
			panic("libsql: Local connection requires Database")
		}
	case Remote:
		if conn.Database == "" || conn.AuthToken == "" {
			panic("libsql: Remote connection requires both Database URL and AuthToken")
		}
	case EmbeddedReplica:
		if conn.Database == "" || conn.PrimaryURL == "" || conn.AuthToken == "" {
			panic("libsql: EmbeddedReplica requires Database, PrimaryURL and AuthToken")
		}
	}

	if cfg.Table == "" {
		cfg.Table = ConfigDefault.Table
	}
	if int(cfg.GCInterval.Seconds()) <= 0 {
		cfg.GCInterval = ConfigDefault.GCInterval
	}

	return cfg
}
