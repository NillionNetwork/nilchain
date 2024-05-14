package keeper

import (
	"context"
	"crypto/sha256"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	metatypes "github.com/NillionNetwork/nillion-chain/x/meta/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService

	bankKeeper types.BankKeeper

	Schema collections.Schema

	// Resources contains the resources mapped by account and resource
	Resources collections.Map[collections.Pair[sdk.AccAddress, []byte], []byte]
}

func NewKeeper(cdc codec.BinaryCodec, storeService store.KVStoreService) Keeper {
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

func forgeResourceKey(fromAcc sdk.AccAddress, resource []byte) collections.Pair[sdk.AccAddress, []byte] {
	hashResource := sha256.Sum256(resource)
	return collections.Join(fromAcc, hashResource[:])
}
