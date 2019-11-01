package keeper

import (
	"github.com/NetCloth/netcloth-chain/modules/ipal/types"
	"github.com/NetCloth/netcloth-chain/modules/params"
)

const (
	DefaultParamspace = types.ModuleName
)

func ParamKeyTable() params.KeyTable { //useless if module has no params
	return params.NewKeyTable()
}
