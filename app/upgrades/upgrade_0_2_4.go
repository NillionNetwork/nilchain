package upgrades

import (
	"context"

	"cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var Upgrade_0_2_4 = Upgrade{
	UpgradeName: "v0.2.4-rc11",
	CreateUpgradeHandler: func(mm module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
		return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			return mm.RunMigrations(ctx, configurator, fromVM)
		}
	},
	StoreUpgrades: types.StoreUpgrades{},
}
