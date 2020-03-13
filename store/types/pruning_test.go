package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPruningOptions_FlushVersion(t *testing.T) {
	t.Parallel()
	require.True(t, PruneEverything.FlushVersion(-1))
	require.True(t, PruneEverything.FlushVersion(0))
	require.True(t, PruneEverything.FlushVersion(1))
	require.True(t, PruneEverything.FlushVersion(2))

	require.True(t, PruneNothing.FlushVersion(-1))
	require.True(t, PruneNothing.FlushVersion(0))
	require.True(t, PruneNothing.FlushVersion(1))
	require.True(t, PruneNothing.FlushVersion(2))

	require.False(t, PruneSyncable.FlushVersion(-1))
	require.True(t, PruneSyncable.FlushVersion(0))
	require.False(t, PruneSyncable.FlushVersion(1))
	require.True(t, PruneSyncable.FlushVersion(100))
	require.False(t, PruneSyncable.FlushVersion(101))
}

func TestPruningOptions_SnapshotVersion(t *testing.T) {
	t.Parallel()
	require.False(t, PruneEverything.SnapshotVersion(-1))
	require.False(t, PruneEverything.SnapshotVersion(0))
	require.False(t, PruneEverything.SnapshotVersion(1))
	require.False(t, PruneEverything.SnapshotVersion(2))

	require.True(t, PruneNothing.SnapshotVersion(-1))
	require.True(t, PruneNothing.SnapshotVersion(0))
	require.True(t, PruneNothing.SnapshotVersion(1))
	require.True(t, PruneNothing.SnapshotVersion(2))

	require.False(t, PruneSyncable.SnapshotVersion(-1))
	require.True(t, PruneSyncable.SnapshotVersion(0))
	require.False(t, PruneSyncable.SnapshotVersion(1))
	require.True(t, PruneSyncable.SnapshotVersion(10000))
	require.False(t, PruneSyncable.SnapshotVersion(10001))
}

func TestPruningOptions_IsValid(t *testing.T) {
	t.Parallel()
	type fields struct {
		KeepEvery     int64
		SnapshotEvery int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"PruneEverything", fields{PruneEverything.KeepEvery, PruneEverything.SnapshotEvery}, true},
		{"PruneNothing", fields{PruneNothing.KeepEvery, PruneNothing.SnapshotEvery}, true},
		{"PruneSyncable", fields{PruneSyncable.KeepEvery, PruneSyncable.SnapshotEvery}, true},
		{"KeepEvery=0", fields{0, 0}, false},
		{"KeepEvery<0", fields{-1, 0}, false},
		{"SnapshotEvery<0", fields{1, -1}, false},
		{"SnapshotEvery%KeepEvery!=0", fields{15, 30}, true},
		{"SnapshotEvery%KeepEvery!=0", fields{15, 20}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			po := PruningOptions{
				KeepEvery:     tt.fields.KeepEvery,
				SnapshotEvery: tt.fields.SnapshotEvery,
			}
			require.Equal(t, tt.want, po.IsValid(), "IsValid() = %v, want %v", po.IsValid(), tt.want)
		})
	}
}
