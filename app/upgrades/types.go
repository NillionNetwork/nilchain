package upgrades

import (
	"cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

type Upgrade struct {
	UpgradeName          string
	CreateUpgradeHandler func(mm module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler
	StoreUpgrades        types.StoreUpgrades
}
