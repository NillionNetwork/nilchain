package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NillionNetwork/nillion-chain/x/meta/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	keeper Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

func (m msgServer) PayFor(ctx context.Context, payFor *types.MsgPayFor) (*types.MsgPayForResponse, error) {
	addr, err := sdk.AccAddressFromBech32(payFor.FromAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	err = m.keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, payFor.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to send coins: %w", err)
	}

	err = m.keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, payFor.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to burn coins: %w", err)
	}

	err = m.keeper.SaveResource(ctx, addr, payFor.Resource)
	if err != nil {
		return nil, fmt.Errorf("failed to save resource: %w", err)
	}

	return &types.MsgPayForResponse{}, nil
}
