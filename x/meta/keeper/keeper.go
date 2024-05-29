package keeper

import (
	"context"
	"crypto/sha256"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	metatypes "github.com/NillionNetwork/nilchain/x/meta/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService

	bankKeeper metatypes.BankKeeper

	Schema collections.Schema

	// Resources contains the resources mapped by account and resource
	Resources collections.Map[collections.Pair[sdk.AccAddress, []byte], []byte]
}

func NewKeeper(cdc codec.BinaryCodec, storeService store.KVStoreService, bankKeeper metatypes.BankKeeper) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		Resources: collections.NewMap(
			sb,
			metatypes.ResourcesPrefix,
			"resources",
			collections.PairKeyCodec(sdk.AccAddressKey, collections.BytesKey),
			collections.BytesValue,
		),
		bankKeeper: bankKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

func (k Keeper) SaveResource(ctx context.Context, fromAddress sdk.AccAddress, resource []byte) error {
	err := k.Resources.Set(ctx, forgeResourceKey(fromAddress, resource), resource)
	if err != nil {
		return fmt.Errorf("failed to set resource: %w", err)
	}

	return nil
}

func (k Keeper) ResourceExists(ctx context.Context, fromAddress sdk.AccAddress, resource []byte) bool {
	exists, err := k.Resources.Has(ctx, forgeResourceKey(fromAddress, resource))
	if err != nil {
		return false
	}

	return exists
}

func (k Keeper) IterateResources(ctx context.Context, cb func(fromAddress sdk.AccAddress, resource []byte)) {
	iterate, err := k.Resources.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	for ; iterate.Valid(); iterate.Next() {
		kv, err := iterate.KeyValue()
		if err != nil {
			panic(err)
		}

		cb(kv.Key.K1(), kv.Value)
	}
}

func forgeResourceKey(fromAcc sdk.AccAddress, resource []byte) collections.Pair[sdk.AccAddress, []byte] {
	hashResource := sha256.Sum256(resource)
	return collections.Join(fromAcc, hashResource[:])
}
