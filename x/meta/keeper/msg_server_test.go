package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/NillionNetwork/nillion-chain/x/meta/keeper"
	"github.com/NillionNetwork/nillion-chain/x/meta/types"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgServer_PayFor(t *testing.T) {
	f := initFixture(t)

	msg := &types.MsgPayFor{
		Resource:    []byte("resource1"),
		FromAddress: "infinity",
		Amount: []types2.Coin{
			{
				Denom:  "coin1",
				Amount: math.NewInt(100),
			},
		},
	}

	_, err := keeper.NewMsgServerImpl(f.keeper).PayFor(f.ctx, msg)
	require.NoError(t, err)
}
