package v0

import (
	"encoding/json"
	"log"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/netcloth/netcloth-chain/app/protocol"
	"github.com/netcloth/netcloth-chain/app/v0/slashing"
	"github.com/netcloth/netcloth-chain/app/v0/staking"
	"github.com/netcloth/netcloth-chain/codec"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func (p *ProtocolV0) ExportAppStateAndValidators(ctx sdk.Context, forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	if forZeroHeight {
		p.prepForZeroHeightGenesis(ctx, jailWhiteList)
	}

	genState := p.moduleManager.ExportGenesis(ctx)
	appState, err = codec.MarshalJSONIndent(p.GetCodec(), genState)
	if err != nil {
		return nil, nil, err
	}
	validators = staking.WriteValidators(ctx, p.StakingKeeper)
	return appState, validators, nil
}

func (p *ProtocolV0) prepForZeroHeightGenesis(ctx sdk.Context, jailWhiteList []string) {
	applyWhiteList := false

	if len(jailWhiteList) > 0 {
		applyWhiteList = true
	}

	whiteListMap := make(map[string]bool)

	for _, addr := range jailWhiteList {
		_, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			log.Fatal(err)
		}
		whiteListMap[addr] = true
	}

	p.StakingKeeper.IterateValidators(ctx, func(_ int64, val staking.ValidatorI) (stop bool) {
		_, _ = p.distrKeeper.WithdrawValidatorCommission(ctx, val.GetOperator())
		return false
	})

	dels := p.StakingKeeper.GetAllDelegations(ctx)
	for _, delegation := range dels {
		_, _ = p.distrKeeper.WithdrawDelegationRewards(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	}

	p.distrKeeper.DeleteAllValidatorSlashEvents(ctx)

	p.distrKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	p.StakingKeeper.IterateValidators(ctx, func(_ int64, val staking.ValidatorI) (stop bool) {
		scraps := p.distrKeeper.GetValidatorOutstandingRewards(ctx, val.GetOperator())
		feePool := p.distrKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps)
		p.distrKeeper.SetFeePool(ctx, feePool)

		p.distrKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		return false
	})

	for _, del := range dels {
		p.distrKeeper.Hooks().BeforeDelegationCreated(ctx, del.DelegatorAddress, del.ValidatorAddress)
		p.distrKeeper.Hooks().AfterDelegationModified(ctx, del.DelegatorAddress, del.ValidatorAddress)
	}

	ctx = ctx.WithBlockHeight(height)

	p.StakingKeeper.IterateRedelegations(ctx, func(_ int64, red staking.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		p.StakingKeeper.SetRedelegation(ctx, red)
		return false
	})

	p.StakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd staking.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		p.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
		return false
	})

	store := ctx.KVStore(protocol.Keys[protocol.StakingStoreKey])
	iter := sdk.KVStoreReversePrefixIterator(store, staking.ValidatorsKey)
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(iter.Key()[1:])
		validator, found := p.StakingKeeper.GetValidator(ctx, addr)
		if !found {
			panic("expected validator, not found")
		}

		validator.UnbondingHeight = 0
		if applyWhiteList && !whiteListMap[addr.String()] {
			validator.Jailed = true
		}

		p.StakingKeeper.SetValidator(ctx, validator)
		counter++
	}

	iter.Close()

	_ = p.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)

	p.slashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdk.ConsAddress, info slashing.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			p.slashingKeeper.SetValidatorSigningInfo(ctx, addr, info)
			return false
		},
	)
}
