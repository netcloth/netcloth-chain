package vm

import (
	"github.com/netcloth/netcloth-chain/server/config"
	"github.com/tendermint/tendermint/crypto"
	"math/big"

	sdk "github.com/netcloth/netcloth-chain/types"
)

type (
	// CanTransferFunc is the signature of a transfer guard function
	CanTransferFunc func(StateDB, sdk.AccAddress, *big.Int) bool
	// TransferFunc is the signature of a transfer function
	TransferFunc func(StateDB, sdk.AccAddress, sdk.AccAddress, *big.Int)
	// GetHashFunc returns the nth block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) sdk.Hash
)

type codeAndHash struct {
	code []byte
	hash sdk.Hash
}

func (c *codeAndHash) Hash() sdk.Hash {
	if c.hash == (sdk.Hash{}) {
		copy(c.hash[:], crypto.Sha256(c.code))
	}
	return c.hash
}

// Context provides the VM with auxiliary information.
// Once provided it shouldn't be modified
type Context struct {
	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers ether from one account to the other
	Transfer TransferFunc
	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc

	// Msg information
	Origin   sdk.AccAddress
	GasPrice *big.Int

	// Block information
	CoinBase    sdk.AccAddress
	GasLimit    uint64
	BlockNumber *big.Int
	Time        *big.Int
}

type VM struct {
	Context

	// depth is the current call stack
	depth int

	// virtual machine configuration options used to initialise the vm
	vmConfig Config
}
