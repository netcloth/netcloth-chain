package aipal

import (
    "github.com/NetCloth/netcloth-chain/modules/aipal/keeper"
    "github.com/NetCloth/netcloth-chain/modules/aipal/types"
)

const (
    ModuleName        = types.ModuleName
    StoreKey          = types.StoreKey
    RouterKey         = types.RouterKey
    QuerierRoute      = types.QuerierRoute
    DefaultCodespace  = types.DefaultCodespace
    DefaultParamspace = keeper.DefaultParamspace
)

var (
    NewKeeper              = keeper.NewKeeper
    NewQuerier             = keeper.NewQuerier
    RegisterCodec          = types.RegisterCodec
    NewServerNodeObject    = types.NewServiceNode
    NewMsgServiceNodeClaim = types.NewMsgServiceNodeClaim
    ErrEmptyInputs         = types.ErrEmptyInputs
    ModuleCdc              = types.ModuleCdc
    AttributeValueCategory = types.AttributeValueCategory
)

type (
    Keeper              = keeper.Keeper
    MsgServiceNodeClaim = types.MsgServiceNodeClaim
)
