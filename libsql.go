package libsql

import (
	"database/sql"
	"fmt"
	"time"
)

// Storage interface implementation
type Storage struct {
	db         *sql.DB
	done       chan bool
	gcInterval time.Duration

	sqlGet    string
	sqlSet    string
	sqlDelete string
	sqlReset  string
	sqlGC     string
}

var (
	dropQuery   = `DROP TABLE IF EXISTS %s;`
	initQueries = [...]string{
		`CREATE TABLE IF NOT EXISTS %s (
			k VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT '',
			v BLOB NOT NULL,
			e BIGINT NOT NULL DEFAULT '0'
		);`,
		`CREATE INDEX IF NOT EXISTS e ON %s (e);`,
	}

	getQuery    = "SELECT v, e FROM %s WHERE k=?;"
	setQuery    = "INSERT OR REPLACE INTO %s (k, v, e) VALUES (?,?,?)"
	deleteQuery = "DELETE FROM %s WHERE k=?"
	resetQuery  = "DELETE FROM %s;"
	gcQuery     = "DELETE FROM %s WHERE e <= ? AND e != 0"
)

// New creates a storage instance
func New(config ...Config) *Storage {
	cfg := configDefault(config...)

	db, err := cfg.Connection.Db()
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	if cfg.Reset {
		if _, err := db.Exec(fmt.Sprintf(dropQuery, cfg.Table)); err != nil {
			_ = db.Close()
			panic(err)
		}
	}

	for _, initQuery := range initQueries {
		if _, err := db.Exec(fmt.Sprintf(initQuery, cfg.Table)); err != nil {
			_ = db.Close()
			panic(err)
		}
	}

	storage := &Storage{
		db:         db,
		done:       make(chan bool),
		gcInterval: cfg.GCInterval,

		sqlGet:    fmt.Sprintf(getQuery, cfg.Table),
		sqlSet:    fmt.Sprintf(setQuery, cfg.Table),
		sqlDelete: fmt.Sprintf(deleteQuery, cfg.Table),
		sqlReset:  fmt.Sprintf(resetQuery, cfg.Table),
		sqlGC:     fmt.Sprintf(gcQuery, cfg.Table),
	}

	go storage.gcTicker()

	return storage
}

// Get gets the value for the given key
func (s *Storage) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}
	row := s.db.QueryRow(s.sqlGet, key)
	var (
		data       = []byte{}
		exp  int64 = 0
	)
	if err := row.Scan(&data, &exp); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if exp != 0 && exp <= time.Now().Unix() {
		return nil, nil
	}
	return data, nil
}

// Set stores the value for the given key
func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}
	var expiry int64
	if exp != 0 {
		expiry = time.Now().Add(exp).Unix()
	}
	_, err := s.db.Exec(s.sqlSet, key, val, expiry)
	return err
}

// Delete deletes the value for the given key
func (s *Storage) Delete(key string) error {
	if len(key) <= 0 {
		return nil
	}
	_, err := s.db.Exec(s.sqlDelete, key)
	return err
}

// Reset resets the storage
func (s *Storage) Reset() error {
	_, err := s.db.Exec(s.sqlReset)
	return err
}

// Close closes the storage, stopping the GC and closing the database connection
func (s *Storage) Close() error {
	s.done <- true
	close(s.done)
	return s.db.Close()
}

// Conn returns the underlying database connection
func (s *Storage) Conn() *sql.DB {
	return s.db
}

func (s *Storage) gcTicker() {
	ticker := time.NewTicker(s.gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if _, err := s.db.Exec(s.sqlGC, time.Now().Unix()); err != nil {
				panic(err)
			}
		case <-s.done:
			return
		}
	}
}
