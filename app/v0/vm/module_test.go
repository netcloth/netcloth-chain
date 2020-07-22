package vm_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/netcloth/netcloth-chain/app/v0/auth"
	"github.com/netcloth/netcloth-chain/app/v0/auth/exported"
	"github.com/netcloth/netcloth-chain/app/v0/vm"
	"github.com/netcloth/netcloth-chain/app/v0/vm/keeper"
	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
	"github.com/netcloth/netcloth-chain/hexutil"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/types/module"
)

type VMTestSuite struct {
	suite.Suite

	ctx      sdk.Context
	ak       auth.AccountKeeper
	vmKeeper vm.Keeper
	vmModule module.AppModule
	handler  sdk.Handler
	acc      exported.Account
	gs       json.RawMessage
}

func (st *VMTestSuite) SetupTest() {
	st.ctx, st.ak, st.vmKeeper, _ = keeper.CreateTestInput(st.T(), false, int64(1000000))
	st.vmModule = vm.NewAppModule(st.vmKeeper)
	st.handler = vm.NewHandler(st.vmKeeper)
	st.acc = st.ak.GetAccount(st.ctx, keeper.Addrs[0])
}

func (st *VMTestSuite) reset() {
	st.SetupTest()
}

// nolint
func getContractAddr(events sdk.Events) (addr sdk.AccAddress, err error) {
	for _, e := range events {
		if e.Type == types.EventTypeNewContract {
			for _, attr := range e.Attributes {
				if string(attr.Key) == types.AttributeKeyAddress {
					addr, err = sdk.AccAddressFromBech32(string(attr.Value))
					return
				}
			}
		}
		break
	}

	return
}
func (st *VMTestSuite) deployContract(acc sdk.AccAddress, code []byte) sdk.AccAddress {
	msg := types.NewMsgContract(acc, sdk.AccAddress{}, code, sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0)))

	result, err := st.handler(st.ctx, msg)
	st.Require().NoError(err, "failed to handle vm create msg")

	addr, err := getContractAddr(result.Events)
	st.Require().NoError(err, "failed to parse contract addr")

	return addr
}

func (st *VMTestSuite) callContract(accAddr, contractAddr sdk.AccAddress, payload []byte, amt sdk.Coin) {
	msg := types.NewMsgContract(accAddr, contractAddr, payload, amt)

	_, err := st.handler(st.ctx, msg)
	st.Require().NoError(err, "failed to handle vm call msg")
}

func (st *VMTestSuite) TestContractExportImport() {
	const (
		// pay.bc with constructor payload[0000000000000000000000000000000000000000000000000000000000000064]
		CodeWithConstructorPayloadHex = "0x608060405261271060005534801561001657600080fd5b506040516107ec3803806107ec8339818101604052602081101561003957600080fd5b8101908080519060200190929190505050610059816100ad60201b60201c565b61005f57fe5b33600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600181905550506100bb565b600080548211159050919050565b610722806100ca6000396000f3fe60806040526004361061007b5760003560e01c80639c75ed6c1161004e5780639c75ed6c1461019c578063d9c90cff146101c7578063de8c50c81461021a578063e662bd25146102455761007b565b80632e1a7d4d146100805780633a6e3d98146100bb57806345596e2e1461010a5780638da5cb5b14610145575b600080fd5b34801561008c57600080fd5b506100b9600480360360208110156100a357600080fd5b8101908080359060200190929190505050610289565b005b3480156100c757600080fd5b506100f4600480360360208110156100de57600080fd5b810190808035906020019092919050505061034c565b6040518082815260200191505060405180910390f35b34801561011657600080fd5b506101436004803603602081101561012d57600080fd5b810190808035906020019092919050505061037e565b005b34801561015157600080fd5b5061015a6103ee565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156101a857600080fd5b506101b1610414565b6040518082815260200191505060405180910390f35b3480156101d357600080fd5b50610200600480360360208110156101ea57600080fd5b810190808035906020019092919050505061041a565b604051808215151515815260200191505060405180910390f35b34801561022657600080fd5b5061022f610428565b6040518082815260200191505060405180910390f35b6102876004803603602081101561025b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061042e565b005b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102e057fe5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610348573d6000803e3d6000fd5b5050565b60006103776000546103696001548561053590919063ffffffff16565b6105bb90919063ffffffff16565b9050919050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146103d557fe5b6103de8161041a565b6103e457fe5b8060018190555050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60015481565b600080548211159050919050565b60005481565b60006104393461034c565b9050600081340390508273ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610488573d6000803e3d6000fd5b507f9ed053bb818ff08b8353cd46f78db1f0799f31c9e4458fdb425c10eccd2efc4433843484604051808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182815260200194505050505060405180910390a1505050565b60008083141561054857600090506105b5565b600082840290508284828161055957fe5b04146105b0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806106cc6021913960400191505060405180910390fd5b809150505b92915050565b60006105fd83836040518060400160405280601a81526020017f536166654d6174683a206469766973696f6e206279207a65726f000000000000815250610605565b905092915050565b600080831182906106b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561067657808201518184015260208101905061065b565b50505050905090810190601f1680156106a35780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5060008385816106bd57fe5b04905080915050939250505056fe536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77a26469706673582212204222d9732198380684d25bb3c9e975cd4ae4a7664c6d2a4d81ac457a590d9d7e64736f6c634300060000330000000000000000000000000000000000000000000000000000000000000064"
		// code after contract deployed by vm
		CodeInVMHex = "0x60806040526004361061007b5760003560e01c80639c75ed6c1161004e5780639c75ed6c1461019c578063d9c90cff146101c7578063de8c50c81461021a578063e662bd25146102455761007b565b80632e1a7d4d146100805780633a6e3d98146100bb57806345596e2e1461010a5780638da5cb5b14610145575b600080fd5b34801561008c57600080fd5b506100b9600480360360208110156100a357600080fd5b8101908080359060200190929190505050610289565b005b3480156100c757600080fd5b506100f4600480360360208110156100de57600080fd5b810190808035906020019092919050505061034c565b6040518082815260200191505060405180910390f35b34801561011657600080fd5b506101436004803603602081101561012d57600080fd5b810190808035906020019092919050505061037e565b005b34801561015157600080fd5b5061015a6103ee565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156101a857600080fd5b506101b1610414565b6040518082815260200191505060405180910390f35b3480156101d357600080fd5b50610200600480360360208110156101ea57600080fd5b810190808035906020019092919050505061041a565b604051808215151515815260200191505060405180910390f35b34801561022657600080fd5b5061022f610428565b6040518082815260200191505060405180910390f35b6102876004803603602081101561025b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061042e565b005b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102e057fe5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610348573d6000803e3d6000fd5b5050565b60006103776000546103696001548561053590919063ffffffff16565b6105bb90919063ffffffff16565b9050919050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146103d557fe5b6103de8161041a565b6103e457fe5b8060018190555050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60015481565b600080548211159050919050565b60005481565b60006104393461034c565b9050600081340390508273ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610488573d6000803e3d6000fd5b507f9ed053bb818ff08b8353cd46f78db1f0799f31c9e4458fdb425c10eccd2efc4433843484604051808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182815260200194505050505060405180910390a1505050565b60008083141561054857600090506105b5565b600082840290508284828161055957fe5b04146105b0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806106cc6021913960400191505060405180910390fd5b809150505b92915050565b60006105fd83836040518060400160405280601a81526020017f536166654d6174683a206469766973696f6e206279207a65726f000000000000815250610605565b905092915050565b600080831182906106b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561067657808201518184015260208101905061065b565b50505050905090810190601f1680156106a35780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5060008385816106bd57fe5b04905080915050939250505056fe536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77a26469706673582212204222d9732198380684d25bb3c9e975cd4ae4a7664c6d2a4d81ac457a590d9d7e64736f6c63430006000033"
		// payload for contract call: method=doTransfer args 'nch12p3yzf6w9xq3v3u48xykktqw7jsge6n5salffn'
		ContractCallPayload = "e662bd25000000000000000000000000506241274e298116479539896b2c0ef4a08cea74"
	)

	ensFactoryCode, err := hexutil.Decode(CodeWithConstructorPayloadHex)
	st.Require().True(err == nil, "failed to decode vm code")
	address := st.deployContract(st.acc.GetAddress(), ensFactoryCode)
	vm.EndBlocker(st.ctx, st.vmKeeper)

	var gs types.GenesisState
	st.Require().NotPanics(func() {
		gsJSON := st.vmModule.ExportGenesis(st.ctx)
		types.ModuleCdc.MustUnmarshalJSON(gsJSON, &gs)
	})

	// sanity check that contract was deployed
	deployedEnsFactoryCode, err := hexutil.Decode(CodeInVMHex)
	st.Require().True(err == nil, "failed to decode vm code")
	code := st.vmKeeper.GetCode(st.ctx, address)
	st.Require().Equal(deployedEnsFactoryCode, code)

	// clear keeper code and re-initialize
	st.vmKeeper.StateDB.SetCode(address, nil)
	gsJSON, err := types.ModuleCdc.MarshalJSON(gs)
	st.Require().NoError(err, "gs to json failed")
	st.vmModule.InitGenesis(st.ctx, gsJSON)
	vm.EndBlocker(st.ctx, st.vmKeeper)

	resCode := st.vmKeeper.StateDB.GetCode(address)
	st.Require().Equal(deployedEnsFactoryCode, resCode)

	// call contract
	callPayload, err := hexutil.Decode(ContractCallPayload)
	st.Require().NoError(err, "parser payload failed")
	st.callContract(st.acc.GetAddress(), address, callPayload, sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(1000)))
	vm.EndBlocker(st.ctx, st.vmKeeper)

	// export vm state after call contrat
	st.gs = st.vmModule.ExportGenesis(st.ctx)

	// reset db
	st.reset()

	// import state
	st.vmModule.InitGenesis(st.ctx, st.gs)
	vm.EndBlocker(st.ctx, st.vmKeeper)

	// export state
	newGS := st.vmModule.ExportGenesis(st.ctx)

	st.Require().True(sdk.JSONEqual(st.gs, newGS))
}

func TestStart(t *testing.T) {
	suite.Run(t, new(VMTestSuite))
}
