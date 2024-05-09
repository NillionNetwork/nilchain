package keeper

import (
	metatypes "github.com/NillionNetwork/nillion-chain/x/meta/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState *metatypes.GenesisState) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *metatypes.GenesisState {
	//TODO implement me
	panic("implement me")
}
