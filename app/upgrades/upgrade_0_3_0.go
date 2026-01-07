package upgrades

import (
	"context"
	"cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
    migrationmngrtypes "github.com/evstack/ev-abci/modules/migrationmngr/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var Upgrade_0_3_0 = Upgrade{
	UpgradeName: "v0.3.0",
	CreateUpgradeHandler: func(mm module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
		return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			return mm.RunMigrations(ctx, configurator, fromVM)
		}
	},
	StoreUpgrades: types.StoreUpgrades{
		Added: []string{
			migrationmngrtypes.ModuleName,
		},
	},
}
