package vm

import (
	"crypto/sha256"
	"math/big"

	"github.com/netcloth/netcloth-chain/modules/vm/common"

	sdk "github.com/netcloth/netcloth-chain/types"
	"golang.org/x/crypto/ripemd160"
)

// Common big integers often used
var (
	big0      = big.NewInt(0)
	big1      = big.NewInt(1)
	big2      = big.NewInt(2)
	big3      = big.NewInt(3)
	big4      = big.NewInt(4)
	big8      = big.NewInt(8)
	big16     = big.NewInt(16)
	big32     = big.NewInt(32)
	big64     = big.NewInt(64)
	big96     = big.NewInt(96)
	big256    = big.NewInt(256)
	big257    = big.NewInt(257)
	big480    = big.NewInt(480)
	big1024   = big.NewInt(1024)
	big3072   = big.NewInt(3072)
	big199680 = big.NewInt(199680)
)

// PrecompiledContract is the basic interface for native Go contracts. The implementation
// requires a deterministic gas count based on the input size of the Run method of the
// contract.
type PrecompiledContract interface {
	RequiredGas(input []byte) uint64  // RequiredPrice calculates the contract gas use
	Run(input []byte) ([]byte, error) // Run runs the precompiled contract
}

// PrecompiledContractsIstanbul contains the default set of pre-compiled Ethereum
// contracts used in the Istanbul release.
var PrecompiledContractsIstanbul = map[string]PrecompiledContract{
	(sdk.BytesToAddress([]byte{1})).String(): &ecrecover{},
	(sdk.BytesToAddress([]byte{2})).String(): &sha256hash{},
	(sdk.BytesToAddress([]byte{3})).String(): &ripemd160hash{},
	(sdk.BytesToAddress([]byte{4})).String(): &dataCopy{},
	//(sdk.BytesToAddress([]byte{5})).String(): &bigModExp{},
	//(sdk.BytesToAddress([]byte{6})).String(): &bn256AddIstanbul{},
	//(sdk.BytesToAddress([]byte{7})).String(): &bn256ScalarMulIstanbul{},
	//(sdk.BytesToAddress([]byte{8})).String(): &bn256PairingIstanbul{},
	//(sdk.BytesToAddress([]byte{9})).String(): &blake2F{},
}

// ECRECOVER implemented as a native contract.
type ecrecover struct{}

func (c *ecrecover) RequiredGas(input []byte) uint64 {
	return EcrecoverGas
}

func (c *ecrecover) Run(input []byte) ([]byte, error) {
	return nil, nil
}

// SHA256 implemented as a native contract.
type sha256hash struct{}

func (c *sha256hash) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*Sha256PerWordGas + Sha256BaseGas
}

func (c *sha256hash) Run(input []byte) ([]byte, error) {
	h := sha256.Sum256(input)
	return h[:], nil
}

// RIPEMD160 implemented as a native contract.
type ripemd160hash struct{}

func (c *ripemd160hash) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*Ripemd160PerWordGas + Ripemd160BaseGas
}

func (c *ripemd160hash) Run(input []byte) ([]byte, error) {
	ripemd := ripemd160.New()
	ripemd.Write(input)
	return common.LeftPadBytes(ripemd.Sum(nil), 32), nil
}

// data copy implemented as a native contract.
type dataCopy struct{}

func (c *dataCopy) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*IdentityPerWordGas + IdentityBaseGas
}

func (c *dataCopy) Run(in []byte) ([]byte, error) {
	return in, nil
}

// bigModExp implements a native big integer exponential modular operation.
type bigModExp struct{}
