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
