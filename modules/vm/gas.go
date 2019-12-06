package vm

import (
	"math/big"

	sdk "github.com/netcloth/netcloth-chain/types"
)

// Gas costs
const (
	GasQuickStep   uint64 = 2
	GasFastestStep uint64 = 3
	GasFastStep    uint64 = 5
	GasMidStep     uint64 = 8
	GasSlowStep    uint64 = 10
	GasExtStep     uint64 = 20
)

// calcGas returns the actual gas cost of the call.
//
// The returned gas is gas - base * 63 / 64.
func callGas(availableGas, base uint64, callCost *big.Int) (uint64, sdk.Error) {
	availableGas = availableGas - base
	gas := availableGas - availableGas/64
	// If the bit length exceeds 64 bit we know that the newly calculated "gas" for EIP150
	// is smaller than the requested amount. Therefor we return the new gas instead
	// of returning an error.
	if !callCost.IsUint64() || gas < callCost.Uint64() {
		return gas, nil
	}

	if !callCost.IsUint64() {
		return 0, ErrGasUintOverflow()
	}

	return callCost.Uint64(), nil
}
