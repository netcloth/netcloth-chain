package types

import (
	"encoding/hex"
	"fmt"
)

const (
	// HashLength is the expected length of the hash
	HashLength = 32
)

// sha256 hash for contract
type Hash [HashLength]byte

// Bytes gets the byte representation of the underlying hash
func (h Hash) Bytes() []byte { return h[:] }

// Hex converts a hash to a hex string.
func (h Hash) Hex() string {
	b := h[:]
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h Hash) String() string {
	return h.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (h Hash) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), h[:])
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}
