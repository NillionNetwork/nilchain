package keeper

import (
	"context"
	metatypes "github.com/NillionNetwork/nillion-chain/x/meta/types"
)

var _ metatypes.QueryServer = QueryServer{}

type QueryServer struct {
	keeper Keeper
}

func NewQueryServer(keeper Keeper) metatypes.QueryServer {
	return &QueryServer{keeper: keeper}
}

func (q QueryServer) MetaData(ctx context.Context, request *metatypes.QueryMetadataRequest) (*metatypes.QueryMetadataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q QueryServer) AllMetaData(ctx context.Context, request *metatypes.QueryAllMetadataRequest) (*metatypes.QueryAllMetadataResponse, error) {
	//TODO implement me
	panic("implement me")
}
