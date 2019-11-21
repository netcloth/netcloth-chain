package vm

import (
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/netcloth/netcloth-chain/types"
	"math/big"
)

// ContractRef is a interface to the contract's backing object
type ContractRef interface {
	Address() sdk.AccAddress
}

// AccountRef implements ContractRef
type AccountRef sdk.AccAddress

// Address casts AccountRef to a Address
func (ar AccountRef) Address() sdk.AccAddress {
	return (sdk.AccAddress)(ar)
}

// Contract implements ContractRef
type Contract struct {
	CallerAddress sdk.AccAddress
	caller        ContractRef
	self          ContractRef

	//jumpdests map[common.Hash]bitvec // Aggregated result of JUMPDEST analysis.
	//analysis  bitvec                 // Locally cached result of JUMPDEST analysis

	Code     []byte
	CodeHash sdk.Hash
	CodeAddr *sdk.AccAddress
	Input    []byte

	Gas   uint64
	value *big.Int
}

func NetContract(caller ContractRef, object ContractRef, value *big.Int, gas uint64) *ContractRef {
	c := &Contract{CallerAddress: caller.Address(), caller: caller, self: object}

	// TODO

	c.Gas = gas
	c.value = value
}

func (c *Contract) validJumpdest(dest *big.Int) bool {
	return false
}

// GetByte returns the n'th byte in the contract's byte array
func (c *Contract) GetByte(n uint64) byte {
	if n < uint64(len(c.Code)) {
		return c.Code[n]
	}
	return 0
}

// Caller returns the caller of the contract.
//
// Caller will recursively call caller when the contract is a delegate
// call, including that of caller's caller.
func (c *Contract) Caller() sdk.AccAddress {
	return c.CallerAddress
}

// UseGas attempts the use gas and subtracts it and returns true on success
func (c *Contract) UseGas(gas uint64) (ok bool) {
	if c.Gas < gas {
		return false
	}
	c.Gas -= gas
	return true
}

// Address returns the contract address
func (c *Contract) Address() sdk.AccAddress {
	return c.self.Address()
}

// Value returns the contract value (sent to it from it's caller)
func (c *Contract) Value() *big.Int {
	return c.value
}

// SetCallCode sets the code of the contract and address of the backing data
// object
func (c *Contract) SetCallCode(addr *sdk.AccAddress, hash sdk.Hash, code []byte) {
	c.Code = code
	c.CodeHash = hash
	c.CodeAddr = addr
}

// SetCodeOptionalHash can be used to provide code, but it's optional to provide hash.
// In case hash is not provided, the jumpdest analysis will not be saved to the parent context
func (c *Contract) SetCodeOptionalHash(addr *sdk.AccAddress, codeAndHash *codeAndHash) {
	c.Code = codeAndHash.code
	c.CodeHash = codeAndHash.hash
	c.CodeAddr = addr
}
