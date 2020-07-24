package common

import (
	"testing"

	"github.com/netcloth/netcloth-chain/client/context"
	"github.com/netcloth/netcloth-chain/codec"

	"github.com/stretchr/testify/require"
)

func TestQueryDelegationRewardsAddrValidation(t *testing.T) {
	cdc := codec.New()
	ctx := context.NewCLIContext().WithCodec(cdc)
	type args struct {
		delAddr string
		valAddr string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"invalid delegator address", args{"invalid", ""}, nil, true},
		{"empty delegator address", args{"", ""}, nil, true},
		{"invalid validator address", args{"nch169j6jw5lv9l76zdqxqcgzz84yrtsu2r7w5s2dg", "invalid"}, nil, true},
		{"empty validator address", args{"nch169j6jw5lv9l76zdqxqcgzz84yrtsu2r7w5s2dg", ""}, nil, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := QueryDelegationRewards(ctx, "", tt.args.delAddr, tt.args.valAddr)
			require.True(t, err != nil, tt.wantErr)
		})
	}
}
