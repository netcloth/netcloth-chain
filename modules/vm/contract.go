package vm

import (
	"math/big"

	sdk "github.com/netcloth/netcloth-chain/types"
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
	// CallerAddress is the result of the caller which initialised this
	// contract. However when the "call method" is delegated this value
	// needs to be initialised to that of the caller's caller.
	CallerAddress sdk.AccAddress
	caller        ContractRef
	self          ContractRef

	jumpdests map[sdk.Hash]bitvec // Aggregated result of JUMPDEST analysis.
	analysis  bitvec              // Locally cached result of JUMPDEST analysis

	Code     []byte
	CodeHash sdk.Hash
	CodeAddr *sdk.AccAddress
	Input    []byte

	Gas   uint64
	value *big.Int
}

func NewContract(caller ContractRef, object ContractRef, value *big.Int, gas uint64) *Contract {
	c := &Contract{CallerAddress: caller.Address(), caller: caller, self: object}

	if parent, ok := caller.(*Contract); ok {
		c.jumpdests = parent.jumpdests
	} else {
		c.jumpdests = make(map[sdk.Hash]bitvec)
	}

	c.Gas = gas
	c.value = value

	return c
}

func (c *Contract) validJumpdest(dest *big.Int) bool {
	udest := dest.Uint64()
	// PC cannot go beyond len(code) and certainly can't be bigger than 63bits.
	// Don't bother checking for JUMPDEST in that case.
	if dest.BitLen() >= 63 || udest >= uint64(len(c.Code)) {
		return false
	}
	// Only JUMPDESTs allowed for destinations
	if OpCode(c.Code[udest]) != JUMPDEST {
		return false
	}
	// Do we have a contract hash already ?
	if c.CodeHash != (sdk.Hash{}) {
		analysis, exist := c.jumpdests[c.CodeHash]
		if !exist {
			analysis = codeBitmap(c.Code)
			c.jumpdests[c.CodeHash] = analysis
		}
		return analysis.codeSegment(udest)
	}

	// We don't have the code hash, most likely a piece of initcode not already
	// in state trie. In that case, we do an analysis, and save it locally, so
	// we don't have to recalculate it for every JUMP instruction in the execution
	// However, we don't save it within the parent context
	if c.analysis == nil {
		c.analysis = codeBitmap(c.Code)
	}
	return c.analysis.codeSegment(udest)
}

func (c *Contract) AsDelegate() *Contract {
	// NOTE: caller must, at all times be a contract.
	parent := c.caller.(*Contract)
	c.CallerAddress = parent.CallerAddress
	c.value = parent.value

	return c
}

// GetOp returns the n'th element in the contract's byte array
func (c *Contract) GetOp(n uint64) OpCode {
	return OpCode(c.GetByte(n))
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
