package types

import "cosmossdk.io/collections"

const (
	// ModuleName is the name of the module
	ModuleName = "meta"

	// StoreKey is the store key string for meta
	StoreKey = ModuleName

	// RouterKey is the message route for meta
	RouterKey = ModuleName
)

// KVStore keys
var (
	ResourcesPrefix = collections.NewPrefix(1)
)
