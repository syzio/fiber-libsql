package libsql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ConfigDefault(t *testing.T) {
	cfg := ConfigDefault

	require.IsType(t, Local{}, cfg.Connection)

	require.Equal(t, "file:./fiber.db", cfg.Connection.(Local).Database)
	require.Equal(t, "fiber_storage", cfg.Table)
	require.False(t, cfg.Reset)
	require.Equal(t, 10, int(cfg.GCInterval.Seconds()))
}
