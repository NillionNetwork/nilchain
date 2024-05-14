package keeper_test

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types"
	"testing"

	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"

	"github.com/NillionNetwork/nillion-chain/x/meta"
	"github.com/NillionNetwork/nillion-chain/x/meta/keeper"
)

type fixture struct {
	ctx    context.Context
	keeper keeper.Keeper
}

func initFixture(t *testing.T) *fixture {
	encCfg := moduletestutil.MakeTestEncodingConfig(meta.AppModuleBasic{})

	mockStoreKey := storetypes.NewKVStoreKey("test")
	storeService := runtime.NewKVStoreService(mockStoreKey)
	k := keeper.NewKeeper(encCfg.Codec, storeService)

	return &fixture{
		ctx:    testutil.DefaultContextWithDB(t, mockStoreKey, storetypes.NewTransientStoreKey("transient_test")).Ctx.WithBlockHeader(cmtproto.Header{}),
		keeper: k,
	}
}

func TestSetResource(t *testing.T) {
	t.Parallel()

	f := initFixture(t)

	addr, err := types.AccAddressFromBech32("cosmos1kc068s88tkyjcc0lkx67x95dwc7hrfm44u8k55")
	require.NoError(t, err)

	err = f.keeper.SaveResource(f.ctx, addr, []byte("resource1"))
	require.NoError(t, err)
}

func TestResourceExists(t *testing.T) {
	t.Parallel()

	f := initFixture(t)

	addr, err := types.AccAddressFromBech32("cosmos1kc068s88tkyjcc0lkx67x95dwc7hrfm44u8k55")
	require.NoError(t, err)

	exists := f.keeper.ResourceExists(f.ctx, addr, []byte("resource1"))
	require.False(t, exists)

	err = f.keeper.SaveResource(f.ctx, addr, []byte("resource1"))
	require.NoError(t, err)

	exists = f.keeper.ResourceExists(f.ctx, addr, []byte("resource1"))
	require.True(t, exists)
}
