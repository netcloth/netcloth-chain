package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommitID(t *testing.T) {
	t.Parallel()
	require.True(t, CommitID{}.IsZero())
	require.False(t, CommitID{Version: int64(1)}.IsZero())
	require.False(t, CommitID{Hash: []byte("x")}.IsZero())
	require.Equal(t, "CommitID{[120 120 120 120]:64}", CommitID{Version: int64(100), Hash: []byte("xxxx")}.String())
}

func TestKVStoreKey(t *testing.T) {
	t.Parallel()
	key := NewKVStoreKey("test")
	require.Equal(t, "test", key.name)
	require.Equal(t, key.name, key.Name())
	require.Equal(t, fmt.Sprintf("KVStoreKey{%p, test}", key), key.String())
}

func TestTransientStoreKey(t *testing.T) {
	t.Parallel()
	key := NewTransientStoreKey("test")
	require.Equal(t, "test", key.name)
	require.Equal(t, key.name, key.Name())
	require.Equal(t, fmt.Sprintf("TransientStoreKey{%p, test}", key), key.String())
}
