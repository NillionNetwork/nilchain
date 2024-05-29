package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	metatypes "github.com/NillionNetwork/nilchain/x/meta/types"
)

func TestImportExportGenesis(t *testing.T) {
	f := initFixture(t)

	genesis := &metatypes.GenesisState{
		Resources: []*metatypes.Resource{
			{
				FromAddress: "cosmos1kc068s88tkyjcc0lkx67x95dwc7hrfm44u8k55",
				Metadata:    []byte("resource1"),
			},
			{
				FromAddress: "cosmos1kc068s88tkyjcc0lkx67x95dwc7hrfm44u8k55",
				Metadata:    []byte("resource2"),
			},
		},
	}

	f.mockedAccountKeeper.EXPECT().GetModuleAccount(f.ctx, metatypes.ModuleName).Return(nil).Times(1)

	f.keeper.InitGenesis(f.ctx, f.mockedAccountKeeper, genesis)

	exists := f.keeper.ResourceExists(
		f.ctx,
		sdk.MustAccAddressFromBech32(genesis.Resources[0].FromAddress),
		genesis.Resources[0].Metadata,
	)
	require.True(t, exists)

	exists = f.keeper.ResourceExists(
		f.ctx,
		sdk.MustAccAddressFromBech32(genesis.Resources[1].FromAddress),
		genesis.Resources[1].Metadata,
	)
	require.True(t, exists)

	exported := f.keeper.ExportGenesis(f.ctx)
	require.Equal(t, genesis, exported)
}
