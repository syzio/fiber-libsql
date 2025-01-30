package libsql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testStorage = New(Config{
	Connection: Local{
		Database: ":memory:",
	},

	Table:      "fiber_storage",
	Reset:      true,
	GCInterval: 10 * time.Second,
})

func Test_LibSQL_Set(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	err := testStorage.Set(key, val, 0)
	require.NoError(t, err)
}

func Test_LibSQL_Get(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	err := testStorage.Set(key, val, 0)
	require.NoError(t, err)

	data, err := testStorage.Get(key)
	require.NoError(t, err)
	require.Equal(t, val, data)
}

func Test_LibSQL_Expiration(t *testing.T) {
	var (
		key = "temp"
		val = []byte("temporary")
	)

	err := testStorage.Set(key, val, 1*time.Second)
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	data, err := testStorage.Get(key)
	require.NoError(t, err)
	require.Nil(t, data)
}

func Test_LibSQL_Delete(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	err := testStorage.Set(key, val, 0)
	require.NoError(t, err)

	err = testStorage.Delete(key)
	require.NoError(t, err)

	data, err := testStorage.Get(key)
	require.NoError(t, err)
	require.Nil(t, data)
}

func Test_LibSQL_Close(t *testing.T) {
	require.Nil(t, testStorage.Close())
}

func Test_LibSQL_Conn(t *testing.T) {
	require.True(t, testStorage.Conn() != nil)
}
