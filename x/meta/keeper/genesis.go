package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	metatypes "github.com/NillionNetwork/nilliond/x/meta/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, ak metatypes.AccountKeeper, genState *metatypes.GenesisState) {
	ak.GetModuleAccount(ctx, metatypes.ModuleName)

	for _, resource := range genState.Resources {
		addr := sdk.MustAccAddressFromBech32(resource.FromAddress)
		err := k.SaveResource(ctx, addr, resource.Metadata)
		if err != nil {
			panic(err)
		}
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *metatypes.GenesisState {
	var resources []*metatypes.Resource
	k.IterateResources(ctx, func(addr sdk.AccAddress, metadata []byte) {
		resources = append(resources, &metatypes.Resource{
			FromAddress: addr.String(),
			Metadata:    metadata,
		})
	})

	return &metatypes.GenesisState{
		Resources: resources,
	}
}
