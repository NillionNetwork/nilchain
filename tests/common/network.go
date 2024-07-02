package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NillionNetwork/nilchain/params"
	json2 "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
	testifysuite "github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"time"
)

const (
	chainRepository = "ghcr.io/nillionnetwork/nilliond"
	haltHeight      = int64(50)
	emptyLogs       = "[]"
	authority       = "nillion10d07y265gmmuvt4z0w9aw880jnsr700jpzdkas"
)

type NetworkTestSuite struct {
	testifysuite.Suite

	Chain *cosmos.CosmosChain
}

// InitChain binary included because version 0.1.1 was a different binary name
func (s *NetworkTestSuite) InitChain(version string, binary string) {
	t := s.T()

	client, network := interchaintest.DockerSetup(t)
	factory := interchaintest.NewBuiltinChainFactory(
		zaptest.NewLogger(t),
		[]*interchaintest.ChainSpec{
			{
				ChainConfig: ibc.ChainConfig{
					Type:    "cosmos",
					Name:    "nilliond",
					ChainID: "nillion-1",
					Images: []ibc.DockerImage{
						{
							Repository: chainRepository,
							Version:    version,
							UidGid:     "1025:1025",
						},
					},
					Bin:            binary,
					Bech32Prefix:   "nillion",
					Denom:          "unillion",
					GasPrices:      "0.00unillion",
					GasAdjustment:  1.3,
					TrustingPeriod: "508h",
					NoHostMount:    false,
					ModifyGenesis: func(config ibc.ChainConfig, bytes []byte) ([]byte, error) {
						genDoc, err := types.GenesisDocFromJSON(bytes)
						if err != nil {
							return nil, fmt.Errorf("failed to unmarshal genesis bytes into genesis doc: %w", err)
						}

						var appState genutiltypes.AppMap
						if err := json.Unmarshal(genDoc.AppState, &appState); err != nil {
							return nil, fmt.Errorf("failed to unmarshal genesis bytes into app state: %w", err)
						}

						cfg := params.MakeTestEncodingConfig()
						govv1beta1.RegisterInterfaces(cfg.InterfaceRegistry)
						cdc := codec.NewProtoCodec(cfg.InterfaceRegistry)

						govGenesisState := &govv1beta1.GenesisState{}
						if err := cdc.UnmarshalJSON(appState[govtypes.ModuleName], govGenesisState); err != nil {
							return nil, fmt.Errorf("failed to unmarshal genesis bytes into gov genesis state: %w", err)
						}

						// set correct minimum deposit using configured denom
						govGenesisState.VotingParams.VotingPeriod = VotingPeriod

						govGenBz, err := cdc.MarshalJSON(govGenesisState)
						if err != nil {
							return nil, fmt.Errorf("failed to marshal gov genesis state: %w", err)
						}

						appState[govtypes.ModuleName] = govGenBz

						genDoc.AppState, err = json.Marshal(appState)
						if err != nil {
							return nil, err
						}

						bz, err := json2.MarshalIndent(genDoc, "", "  ")
						if err != nil {
							return nil, err
						}

						return bz, nil
					},
				},
			},
		},
	)

	chains, err := factory.Chains(t.Name())
	s.Require().NoError(err, "error creating chains")

	ic := interchaintest.NewInterchain().
		AddChain(chains[0])

	require.NoError(t, ic.Build(context.Background(), nil, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),

		SkipPathCreation: false,
	}))

	s.Chain = chains[0].(*cosmos.CosmosChain)
}

// UpgradeChain
func (s *NetworkTestSuite) UpgradeChain(ctx context.Context, chain *cosmos.CosmosChain, wallet ibc.Wallet, upgradeVersion string) {
	planName := "upgrade-test"
	s.ExecuteGovProposal(ctx, chain, wallet, planName, upgradeVersion)

	height, err := chain.Height(ctx)
	s.Require().NoError(err, "error fetching height before upgrade")

	timeoutCtx, timeoutCtxCancel := context.WithTimeout(ctx, time.Minute*2)
	defer timeoutCtxCancel()

	err = testutil.WaitForBlocks(timeoutCtx, int(haltHeight-height)+1, chain)
	s.Require().Error(err, "chain did not halt at halt height")

	//err = chain.StopAllNodes(ctx)
	//s.Require().NoError(err, "error stopping node(s)")
	//
	//chain.UpgradeVersion(ctx, s.DockerClient, upgradeVersion)
	//
	//err = chain.StartAllNodes(ctx)
	//s.Require().NoError(err, "error starting upgraded node(s)")
	//
	//timeoutCtx, timeoutCtxCancel = context.WithTimeout(ctx, time.Minute*2)
	//defer timeoutCtxCancel()
	//
	//err = test.WaitForBlocks(timeoutCtx, int(blocksAfterUpgrade), chain)
	//s.Require().NoError(err, "chain did not produce blocks after upgrade")
	//
	//height, err = chain.Height(ctx)
	//s.Require().NoError(err, "error fetching height after upgrade")
	//
	//s.Require().GreaterOrEqual(height, haltHeight+blocksAfterUpgrade, "height did not increment enough after upgrade")
}

// AssertValidTxResponse verifies that an sdk.TxResponse
// has non-empty values.
func (s *NetworkTestSuite) AssertValidTxResponse(resp sdk.TxResponse) {
	respLogsMsg := resp.Logs.String()
	if respLogsMsg == emptyLogs {
		respLogsMsg = resp.RawLog
	}
	s.Require().NotEqual(int64(0), resp.GasUsed, respLogsMsg)
	s.Require().NotEqual(int64(0), resp.GasWanted, respLogsMsg)
	s.Require().NotEmpty(resp.Events, respLogsMsg)
	s.Require().NotEmpty(resp.Data, respLogsMsg)
}
