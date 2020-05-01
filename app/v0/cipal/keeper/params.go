package keeper

import (
	"github.com/netcloth/netcloth-chain/app/v0/cipal/types"
	"github.com/netcloth/netcloth-chain/app/v0/params"
)

const (
	DefaultParamspace = types.ModuleName
)

func ParamKeyTable() params.KeyTable { //useless if module has no params
	return params.NewKeyTable()
}
