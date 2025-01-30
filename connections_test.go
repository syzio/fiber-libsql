package libsql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	_ "github.com/tursodatabase/go-libsql"
)

func Test_Local_Db(t *testing.T) {
	conn := Local{Database: ":memory:"}
	db, err := conn.Db()

	require.NoError(t, err)
	require.NoError(t, db.Ping())
}

func Test_Remote_Db(t *testing.T) {
	t.Skip("Requires actual remote database connection")

	url := "libsql://example.com"
	token := "super-secret-token"

	conn := Remote{
		Database:  url,
		AuthToken: token,
	}
	db, err := conn.Db()

	require.NoError(t, err)
	require.NoError(t, db.Ping())
}

func Test_EmbeddedReplica_Db(t *testing.T) {
	t.Skip("Requires actual remote database connection")

	conn := EmbeddedReplica{
		Database:      ":memory:",
		PrimaryURL:    "libsql://example.com",
		AuthToken:     "super-secret-token",
		EncryptionKey: "encryption-key",
		SyncInterval:  time.Second,
	}

	db, err := conn.Db()
	require.NoError(t, err)
	require.NoError(t, db.Ping())
}
