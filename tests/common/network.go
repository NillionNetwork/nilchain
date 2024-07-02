package common

import (
	"context"
	"time"

	"github.com/docker/docker/client"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
	testifysuite "github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

const (
	chainRepository    = "ghcr.io/nillionnetwork/nilliond"
	haltHeight         = int64(50)
	authority          = "nillion10d07y265gmmuvt4z0w9aw880jnsr700jpzdkas"
	blocksAfterUpgrade = int64(10)
)

type NetworkTestSuite struct {
	testifysuite.Suite

	Chain        *cosmos.CosmosChain
	DockerClient *client.Client
}

// InitChain binary included because version 0.1.1 was a different binary name
func (s *NetworkTestSuite) InitChain(version string, binary string) {
	t := s.T()

	if version == "v0.1.1" {
		binary = "nilliond"
	}
	binary = "nilchaind"

	client, network := interchaintest.DockerSetup(t)
	s.DockerClient = client
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
					ModifyGenesis: cosmos.ModifyGenesis(
						[]cosmos.GenesisKV{
							cosmos.NewGenesisKV("app_state.gov.params.expedited_voting_period", "30s"),
						},
					),
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

	err = chain.StopAllNodes(ctx)
	s.Require().NoError(err, "error stopping node(s)")

	chain.UpgradeVersion(ctx, s.DockerClient, chainRepository, upgradeVersion)

	err = chain.StartAllNodes(ctx)
	s.Require().NoError(err, "error starting upgraded node(s)")

	timeoutCtx, timeoutCtxCancel = context.WithTimeout(ctx, time.Minute*2)
	defer timeoutCtxCancel()

	err = testutil.WaitForBlocks(timeoutCtx, int(blocksAfterUpgrade), chain)
	s.Require().NoError(err, "chain did not produce blocks after upgrade")

	height, err = chain.Height(ctx)
	s.Require().NoError(err, "error fetching height after upgrade")

	s.Require().GreaterOrEqual(height, haltHeight+blocksAfterUpgrade, "height did not increment enough after upgrade")
}
