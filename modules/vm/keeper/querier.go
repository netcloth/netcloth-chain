package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParameters(ctx, k)
		case types.QueryContractCode:
			return queryCode(ctx, req, k)
		case types.QueryContractState:
			return queryContractState(ctx, req, k)
		case types.QueryStorage:
			return queryStorage(ctx, path, k)
		case types.QueryTxLogs:
			return queryTxLogs(ctx, path, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown vm query endpoint")
		}
	}
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return res, nil
}

func queryCode(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	if len(req.Data) != 20 {
		return nil, sdk.ErrInvalidAddress("address invalid")
	}

	accAddr := sdk.AccAddress(req.Data)
	code := k.GetCode(ctx, accAddr)

	return code, nil
}

func queryContractState(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var p types.QueryContractStateParams
	codec.Cdc.UnmarshalJSON(req.Data, p)

	//st := vm.StateTransition{
	//	Sender:    p.Addr,
	//	Recipient: p.Addr,
	//	Price:     sdk.NewInt(1000000),
	//	GasLimit:  10000000,
	//	Payload:   p.Payload,
	//	StateDB:   k.StateDB.WithContext(ctx),
	//}
	//
	//evmCtx := vm.Context{
	//	CanTransfer: st.CanTransfer,
	//	Transfer:    st.Transfer,
	//	GetHash:     st.GetHash,
	//	Origin:      st.Sender,
	//	GasPrice:    st.Price.BigInt(),
	//	CoinBase:    ctx.BlockHeader().ProposerAddress,
	//	GasLimit:    st.GasLimit,
	//	BlockNumber: sdk.NewInt(ctx.BlockHeader().Height).BigInt(),
	//}
	//
	//cfg := vm.Config{}
	//
	//evm := vm.NewEVM(evmCtx, st.StateDB, cfg)
	//
	//var (
	//	ret         []byte
	//	leftOverGas uint64
	//	err         sdk.Error
	//)
	//
	//ret, leftOverGas, err = evm.Call(st.Sender, st.Recipient, st.Payload, 1000000000, st.Amount.BigInt())
	//fmt.Print(leftOverGas)
	//if err != nil {
	//	return nil, err
	//}

	//return ret, nil

	return nil, nil
}

func queryStorage(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr, _ := sdk.AccAddressFromBech32(path[1])
	key := sdk.HexToHash(path[2])
	val := keeper.GetState(ctx, addr, key)
	bRes := types.QueryResStorage{Value: val.Bytes()}
	res, err := codec.MarshalJSONIndent(keeper.cdc, bRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}
	return res, nil
}

func queryTxLogs(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	txHash := sdk.HexToHash(path[1])
	logs := keeper.GetLogs(ctx, txHash)

	bRes := types.QueryLogs{Logs: logs}
	res, err := codec.MarshalJSONIndent(keeper.cdc, bRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}

	return res, nil
}
