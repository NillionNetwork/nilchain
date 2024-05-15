package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	metatypes "github.com/NillionNetwork/nillion-chain/x/meta/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, ak metatypes.AccountKeeper, _ *metatypes.GenesisState) {
	fmt.Printf("init the genesis")
	ak.GetModuleAccount(ctx, metatypes.ModuleName)
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *metatypes.GenesisState {
	return &metatypes.GenesisState{}
}
