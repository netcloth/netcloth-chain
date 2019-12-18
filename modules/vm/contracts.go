package vm

import (
	"crypto/sha256"
	"math/big"

	"golang.org/x/crypto/ripemd160"

	"github.com/netcloth/netcloth-chain/modules/vm/common"
	"github.com/netcloth/netcloth-chain/modules/vm/common/math"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// PrecompiledContract is the basic interface for native Go contracts. The implementation
// requires a deterministic gas count based on the input size of the Run method of the
// contract.
type PrecompiledContract interface {
	RequiredGas(input []byte) uint64      // RequiredPrice calculates the contract gas use
	Run(input []byte) ([]byte, sdk.Error) // Run runs the precompiled contract
}

// PrecompiledContracts contains the default set of pre-compiled contracts used in the Istanbul release.
var PrecompiledContracts = map[string]PrecompiledContract{
	(sdk.BytesToAddress([]byte{1})).String(): &ecrecover{},
	(sdk.BytesToAddress([]byte{2})).String(): &sha256hash{},
	(sdk.BytesToAddress([]byte{3})).String(): &ripemd160hash{},
	(sdk.BytesToAddress([]byte{4})).String(): &dataCopy{},
	(sdk.BytesToAddress([]byte{5})).String(): &bigModExp{},
	//(sdk.BytesToAddress([]byte{6})).String(): &bn256Add{},
	//(sdk.BytesToAddress([]byte{7})).String(): &bn256ScalarMul{},
	//(sdk.BytesToAddress([]byte{8})).String(): &bn256Pairing{},
	//(sdk.BytesToAddress([]byte{9})).String(): &blake2F{},
}

var (
	// true32Byte is returned if the bn256 pairing check succeeds.
	true32Byte = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	// false32Byte is returned if the bn256 pairing check fails.
	false32Byte = make([]byte, 32)
)

// RunPrecompiledContract runs and evaluates the output of a precompiled contract.
func RunPrecompiledContract(p PrecompiledContract, input []byte, contract *Contract) (ret []byte, err sdk.Error) {
	gas := p.RequiredGas(input)
	if contract.UseGas(gas) {
		return p.Run(input)
	}
	return nil, ErrOutOfGas()
}

// ECRECOVER implemented as a native contract.
type ecrecover struct{}

func (c *ecrecover) RequiredGas(input []byte) uint64 {
	return EcrecoverGas
}

func (c *ecrecover) Run(input []byte) ([]byte, sdk.Error) {
	const ecRecoverInputLength = 128

	// TODO
	//input = common.RightPadBytes(input, ecRecoverInputLength)
	//r := new(big.Int).SetBytes(input[64:96])
	//s := new(big.Int).SetBytes(input[96:128])
	//v := input[63] - 27
	//
	//if !allZero(input[32:64]) || !common.ValidateSignatureValues(v, r, s) {
	//	return nil, nil
	//}
	//
	//sig := make([]byte, 65)
	//copy(sig, input[64:128])
	//sig[64] = v

	return nil, nil
}

// SHA256 implemented as a native contract.
type sha256hash struct{}

func (c *sha256hash) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*Sha256PerWordGas + Sha256BaseGas
}

func (c *sha256hash) Run(input []byte) ([]byte, sdk.Error) {
	h := sha256.Sum256(input)
	return h[:], nil
}

// RIPEMD160 implemented as a native contract.
type ripemd160hash struct{}

func (c *ripemd160hash) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*Ripemd160PerWordGas + Ripemd160BaseGas
}

func (c *ripemd160hash) Run(input []byte) ([]byte, sdk.Error) {
	ripemd := ripemd160.New()
	ripemd.Write(input)
	return common.LeftPadBytes(ripemd.Sum(nil), 32), nil
}

// data copy implemented as a native contract.
type dataCopy struct{}

func (c *dataCopy) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*IdentityPerWordGas + IdentityBaseGas
}

func (c *dataCopy) Run(in []byte) ([]byte, sdk.Error) {
	return in, nil
}

// bigModExp implements a native big integer exponential modular operation.
type bigModExp struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bigModExp) RequiredGas(input []byte) uint64 {
	var (
		baseLen = new(big.Int).SetBytes(getData(input, 0, 32))
		expLen  = new(big.Int).SetBytes(getData(input, 32, 32))
		modLen  = new(big.Int).SetBytes(getData(input, 64, 32))
	)
	if len(input) > 96 {
		input = input[96:]
	} else {
		input = input[:0]
	}
	// Retrieve the head 32 bytes of exp for the adjusted exponent length
	var expHead *big.Int
	if big.NewInt(int64(len(input))).Cmp(baseLen) <= 0 {
		expHead = new(big.Int)
	} else {
		if expLen.Cmp(common.Big32) > 0 {
			expHead = new(big.Int).SetBytes(getData(input, baseLen.Uint64(), 32))
		} else {
			expHead = new(big.Int).SetBytes(getData(input, baseLen.Uint64(), expLen.Uint64()))
		}
	}
	// Calculate the adjusted exponent length
	var msb int
	if bitlen := expHead.BitLen(); bitlen > 0 {
		msb = bitlen - 1
	}
	adjExpLen := new(big.Int)
	if expLen.Cmp(common.Big32) > 0 {
		adjExpLen.Sub(expLen, common.Big32)
		adjExpLen.Mul(common.Big8, adjExpLen)
	}
	adjExpLen.Add(adjExpLen, big.NewInt(int64(msb)))

	// Calculate the gas cost of the operation
	gas := new(big.Int).Set(math.BigMax(modLen, baseLen))
	switch {
	case gas.Cmp(common.Big64) <= 0:
		gas.Mul(gas, gas)
	case gas.Cmp(common.Big1024) <= 0:
		gas = new(big.Int).Add(
			new(big.Int).Div(new(big.Int).Mul(gas, gas), common.Big4),
			new(big.Int).Sub(new(big.Int).Mul(common.Big96, gas), common.Big3072),
		)
	default:
		gas = new(big.Int).Add(
			new(big.Int).Div(new(big.Int).Mul(gas, gas), common.Big16),
			new(big.Int).Sub(new(big.Int).Mul(common.Big480, gas), common.Big199680),
		)
	}
	gas.Mul(gas, math.BigMax(adjExpLen, common.Big1))
	gas.Div(gas, new(big.Int).SetUint64(ModExpQuadCoeffDiv))

	if gas.BitLen() > 64 {
		return math.MaxUint64
	}
	return gas.Uint64()
}

func (c *bigModExp) Run(input []byte) ([]byte, sdk.Error) {
	var (
		baseLen = new(big.Int).SetBytes(getData(input, 0, 32)).Uint64()
		expLen  = new(big.Int).SetBytes(getData(input, 32, 32)).Uint64()
		modLen  = new(big.Int).SetBytes(getData(input, 64, 32)).Uint64()
	)
	if len(input) > 96 {
		input = input[96:]
	} else {
		input = input[:0]
	}
	// Handle a special case when both the base and mod length is zero
	if baseLen == 0 && modLen == 0 {
		return []byte{}, nil
	}
	// Retrieve the operands and execute the exponentiation
	var (
		base = new(big.Int).SetBytes(getData(input, 0, baseLen))
		exp  = new(big.Int).SetBytes(getData(input, baseLen, expLen))
		mod  = new(big.Int).SetBytes(getData(input, baseLen+expLen, modLen))
	)
	if mod.BitLen() == 0 {
		// Modulo 0 is undefined, return zero
		return common.LeftPadBytes([]byte{}, int(modLen)), nil
	}
	return common.LeftPadBytes(base.Exp(base, exp, mod).Bytes(), int(modLen)), nil
}
