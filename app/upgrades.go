package nillionapp

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	metatypes "github.com/NillionNetwork/nilchain/x/meta/types"
)

// this upgrade adds the metatypes.StoreKey to the store
const (
	UpgradeName   = "v020"
	upgradeHeight = 10
)

func (app NillionApp) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		},
	)

	if !app.UpgradeKeeper.IsSkipHeight(upgradeHeight) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				metatypes.StoreKey,
			},
		}
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeHeight, &storeUpgrades))
	}

}
