package meta

import (
	"context"
	"cosmossdk.io/client/v2/autocli"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"time"

	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/NillionNetwork/nillion-chain/x/meta/client/cli"
	"github.com/NillionNetwork/nillion-chain/x/meta/keeper"
	metatypes "github.com/NillionNetwork/nillion-chain/x/meta/types"
)

var (
	_ module.AppModuleBasic      = AppModule{}
	_ module.HasGenesis          = AppModule{}
	_ module.HasServices         = AppModule{}
	_ module.HasConsensusVersion = AppModule{}
	_ autocli.HasAutoCLIConfig   = AppModule{}

	_ appmodule.AppModule = AppModule{}
)

// ConsensusVersion defines the current meta module consensus version.
const ConsensusVersion = 1

type AppModuleBasic struct {
}

func (a AppModuleBasic) Name() string {
	return metatypes.ModuleName
}

func (a AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

func (a AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	metatypes.RegisterInterfaces(registry)
}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := metatypes.RegisterQueryHandlerClient(context.Background(), mux, metatypes.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (a AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(metatypes.DefaultGenesisState())
}

func (a AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	var genesisState metatypes.GenesisState
	if err := cdc.UnmarshalJSON(message, &genesisState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", metatypes.ModuleName, err)
	}

	return genesisState.Validate()
}

func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.TxCmd()
}

type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
}

func (a AppModule) ConsensusVersion() uint64 {
	return ConsensusVersion
}

func (a AppModule) RegisterServices(configurator module.Configurator) {
	metatypes.RegisterMsgServer(configurator.MsgServer(), keeper.NewMsgServerImpl(a.keeper))
	metatypes.RegisterQueryServer(configurator.MsgServer(), keeper.NewQueryServer(a.keeper))
}

func (a AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	start := time.Now()
	var genesisState metatypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	telemetry.MeasureSince(start, "InitGenesis", "meta", "unmarshal")

	a.keeper.InitGenesis(ctx, &genesisState)
}

func (a AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := a.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

func (a AppModule) IsOnePerModuleType() {}

func (a AppModule) IsAppModule() {}
