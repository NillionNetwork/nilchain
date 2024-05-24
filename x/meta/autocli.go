package meta

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
)

func (a AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Tx: &autocliv1.ServiceCommandDescriptor{
			EnhanceCustomCommand: true,
		},
	}
}
