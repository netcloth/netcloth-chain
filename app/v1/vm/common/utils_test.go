package common

import (
	"bytes"
	"testing"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestCreateAddress2(t *testing.T) {
	type testcase struct {
		origin   string
		salt     string
		code     string
		expected string
	}

	for i, tt := range []testcase{
		{
			origin:   "nch1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq7hadyk",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0x00",
			expected: "nch1l2l35x4c06nzdemz4uy3grc0vk2gzphcdugdch",
		},
		{
			origin:   "nch1l2l35x4c06nzdemz4uy3grc0vk2gzphcdugdch",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0x00",
			expected: "nch17tj3hmt5fdae3e7j2p3lw2xnt84aanr79f20uw",
		},
		{
			origin:   "nch17tj3hmt5fdae3e7j2p3lw2xnt84aanr79f20uw",
			salt:     "0xfeed000000000000000000000000000000000000",
			code:     "0x00",
			expected: "nch1y3rq2v0yy2ugcpkjyfsxh7p3jpfu08zphw82ze",
		},
		{
			origin:   "nch1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq7hadyk",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0xdeadbeef",
			expected: "nch1lxz5m0z5hyj4hgagthm5n4esghpazz8e2ajxxt",
		},
		{
			origin:   "nch1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq7hadyk",
			salt:     "0xcafebabe",
			code:     "0xdeadbeef",
			expected: "nch1jh9ny6pmtm42aa2myl2k0mnq0dmrtag06fhm3t",
		},
		{
			origin:   "nch1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq7hadyk",
			salt:     "0xcafebabe",
			code:     "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
			expected: "nch1ed399kjapfxqqsea6dkk2uxrzepq8cjw5fxmql",
		},
		{
			origin:   "nch1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq7hadyk",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0x",
			expected: "nch1sv0mfghrpsp0usxcpc9eax8u2d9t4mx8qeglq8",
		},
	} {

		origin, _ := sdk.AccAddressFromBech32(tt.origin)
		salt := sdk.BytesToHash(FromHex(tt.salt))
		codeHash := crypto.Sha256(FromHex(tt.code))
		address := CreateAddress2(origin, salt, codeHash)

		expected, _ := sdk.AccAddressFromBech32(tt.expected)
		if !bytes.Equal(expected.Bytes(), address.Bytes()) {
			t.Errorf("test %d: expected %s, got %s", i, expected.String(), address.String())
		}

	}
}
