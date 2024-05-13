package keeper

import (
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
}

func NewKeeper(cdc codec.BinaryCodec, storeService store.KVStoreService) Keeper {
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
	}
}
