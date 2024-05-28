package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/NillionNetwork/nilliond/x/meta/keeper"
	"github.com/NillionNetwork/nilliond/x/meta/types"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServer_PayFor(t *testing.T) {
	f := initFixture(t)

	msg := &types.MsgPayFor{
		Resource:    []byte("resource1"),
		FromAddress: "cosmos1kc068s88tkyjcc0lkx67x95dwc7hrfm44u8k55",
		Amount: []types2.Coin{
			{
				Denom:  "coin1",
				Amount: math.NewInt(100),
			},
		},
	}

	fromAddress, err := types2.AccAddressFromBech32(msg.FromAddress)
	require.NoError(t, err)

	f.mockedBankKeeper.EXPECT().
		SendCoinsFromAccountToModule(f.ctx, fromAddress, types.ModuleName, msg.Amount)

	f.mockedBankKeeper.EXPECT().
		BurnCoins(f.ctx, types.ModuleName, msg.Amount)

	_, err = keeper.NewMsgServerImpl(f.keeper).PayFor(f.ctx, msg)
	require.NoError(t, err)

	// Check if the resource was saved
	exists := f.keeper.ResourceExists(f.ctx, fromAddress, msg.Resource)
	require.True(t, exists)
}
