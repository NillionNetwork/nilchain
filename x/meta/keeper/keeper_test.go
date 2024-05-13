package keeper_test

import (
	"context"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"testing"

	"github.com/NillionNetwork/nillion-chain/x/meta"
	"github.com/NillionNetwork/nillion-chain/x/meta/keeper"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
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

func TestGetAuthority(t *testing.T) {
	t.Parallel()
}
