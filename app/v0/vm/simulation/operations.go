package simulation

import (
	"encoding/hex"
	"math/big"
	"math/rand"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/netcloth/netcloth-chain/app/simapp/helpers"
	"github.com/netcloth/netcloth-chain/app/v0/simulation"
	"github.com/netcloth/netcloth-chain/app/v0/vm/keeper"
	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
	"github.com/netcloth/netcloth-chain/baseapp"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
	simtypes "github.com/netcloth/netcloth-chain/types/simulation"
)

func WeightedOperations(appParams simtypes.AppParams, cdc *codec.Codec, ak types.AccountKeeper, k keeper.Keeper) simulation.WeightedOperations {
	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			10,
			SimulateMsgContractCreate(ak, k),
		),
		simulation.NewWeightedOperation(
			100,
			SimulateMsgContractCall(ak, k),
		),
	}
}

const (
	// contract info for pay.sol: https://github.com/netcloth/contracts/blob/master/payment/pay.sol
	codeStr = "608060405261271060005534801561001657600080fd5b506040516107ec3803806107ec8339818101604052602081101561003957600080fd5b8101908080519060200190929190505050610059816100ad60201b60201c565b61005f57fe5b33600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600181905550506100bb565b600080548211159050919050565b610722806100ca6000396000f3fe60806040526004361061007b5760003560e01c80639c75ed6c1161004e5780639c75ed6c1461019c578063d9c90cff146101c7578063de8c50c81461021a578063e662bd25146102455761007b565b80632e1a7d4d146100805780633a6e3d98146100bb57806345596e2e1461010a5780638da5cb5b14610145575b600080fd5b34801561008c57600080fd5b506100b9600480360360208110156100a357600080fd5b8101908080359060200190929190505050610289565b005b3480156100c757600080fd5b506100f4600480360360208110156100de57600080fd5b810190808035906020019092919050505061034c565b6040518082815260200191505060405180910390f35b34801561011657600080fd5b506101436004803603602081101561012d57600080fd5b810190808035906020019092919050505061037e565b005b34801561015157600080fd5b5061015a6103ee565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156101a857600080fd5b506101b1610414565b6040518082815260200191505060405180910390f35b3480156101d357600080fd5b50610200600480360360208110156101ea57600080fd5b810190808035906020019092919050505061041a565b604051808215151515815260200191505060405180910390f35b34801561022657600080fd5b5061022f610428565b6040518082815260200191505060405180910390f35b6102876004803603602081101561025b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061042e565b005b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102e057fe5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610348573d6000803e3d6000fd5b5050565b60006103776000546103696001548561053590919063ffffffff16565b6105bb90919063ffffffff16565b9050919050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146103d557fe5b6103de8161041a565b6103e457fe5b8060018190555050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60015481565b600080548211159050919050565b60005481565b60006104393461034c565b9050600081340390508273ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610488573d6000803e3d6000fd5b507f9ed053bb818ff08b8353cd46f78db1f0799f31c9e4458fdb425c10eccd2efc4433843484604051808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182815260200194505050505060405180910390a1505050565b60008083141561054857600090506105b5565b600082840290508284828161055957fe5b04146105b0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806106cc6021913960400191505060405180910390fd5b809150505b92915050565b60006105fd83836040518060400160405280601a81526020017f536166654d6174683a206469766973696f6e206279207a65726f000000000000815250610605565b905092915050565b600080831182906106b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561067657808201518184015260208101905061065b565b50505050905090810190601f1680156106a35780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5060008385816106bd57fe5b04905080915050939250505056fe536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77a26469706673582212204222d9732198380684d25bb3c9e975cd4ae4a7664c6d2a4d81ac457a590d9d7e64736f6c63430006000033"
	abiStr  = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_feeRateE4\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_actual_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"E4\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"calcCommission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address payable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"doTransfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeRateE4\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_feeRateE4\",\"type\":\"uint256\"}],\"name\":\"feeRateValid\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address payable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_feeRateE4\",\"type\":\"uint256\"}],\"name\":\"setFeeRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
)

func SimulateMsgContractCreate(ak types.AccountKeeper, k keeper.Keeper) simtypes.Operation {

	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		a := baseapp.DereferenceBaseApp(app)
		if a == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgContractCreate, "app invalid"), nil, nil
		}

		acc, _ := simtypes.RandomAcc(r, accs)
		accountObj := ak.GetAccount(ctx, acc.Address)

		code, err := hex.DecodeString(codeStr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgContractCreate, "wrong contract code"), nil, err
		}

		abiObj, err := abi.JSON(strings.NewReader(abiStr))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgContractCreate, "wrong contract abi"), nil, err
		}

		args := []interface{}{big.NewInt(10)}

		payload, err := abiObj.Constructor.Inputs.PackValues(args)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgContractCreate, "gen payload failed"), nil, err
		}

		code = append(code, payload...)

		msg := types.NewMsgContract(acc.Address, nil, code, sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0)))

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.NewCoins(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(1000000))),
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{accountObj.GetAccountNumber()},
			[]uint64{accountObj.GetSequence()},
			acc.PrivKey,
		)

		_, _, err = a.Deliver(tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateMsgContractCall(ak types.AccountKeeper, k keeper.Keeper) simtypes.Operation {

	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		a := baseapp.DereferenceBaseApp(app)
		if a == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgContractCall, "app invalid"), nil, nil
		}

		acc, _ := simtypes.RandomAcc(r, accs)
		accountObj := ak.GetAccount(ctx, acc.Address)

		abiObj, err := abi.JSON(strings.NewReader(abiStr))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgContractCall, "wrong contract abi"), nil, err
		}

		m, ok := abiObj.Methods["doTransfer"]
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgContractCall, "method not exist in abi"), nil, err
		}

		var Address [20]byte
		copy(Address[:], acc.Address)

		args := []interface{}{Address}
		payload := m.ID
		argsBin, err := m.Inputs.PackValues(args)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgContractCall, "gen payload failed"), nil, err
		}
		payload = append(payload, argsBin...)

		contractAddrs := k.GetAllHostContractAddresses(ctx)
		if len(contractAddrs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgContractCall, "no contract"), nil, nil
		}

		contractIndex := r.Intn(len(contractAddrs))
		msg := types.NewMsgContract(acc.Address, contractAddrs[contractIndex], payload, sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(10000)))

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.NewCoins(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(1000000))),
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{accountObj.GetAccountNumber()},
			[]uint64{accountObj.GetSequence()},
			acc.PrivKey,
		)

		_, _, err = a.Deliver(tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
