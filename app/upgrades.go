package nillionapp

import (
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"fmt"
	"github.com/NillionNetwork/nilchain/app/upgrades"
)

var Upgrades = []upgrades.Upgrade{
	upgrades.Upgrade_0_2_1,
	upgrades.Upgrade_0_2_4,
}

func (app NillionApp) setupUpgradeHandlers() {
	for _, upgrade := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(
				*app.ModuleManager,
				app.configurator,
			),
		)
	}
}

func (app NillionApp) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrade := range Upgrades {
		if upgradeInfo.Name == upgrade.UpgradeName {
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &upgrade.StoreUpgrades))
		}
	}
}
