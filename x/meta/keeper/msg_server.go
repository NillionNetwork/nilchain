package keeper

import (
	"context"
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
	//TODO implement me
	panic("implement me")
}
