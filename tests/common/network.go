package common

import (
	"context"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/require"
	testifysuite "github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

const (
	chainRepository = "ghcr.io/nillionnetwork/nilliond"
	haltHeight      = uint64(50)
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
					Type: "cosmos",
					Name: "ibc-go-simd",
					Images: []ibc.DockerImage{
						{
							Repository: chainRepository,
							Version:    version,
							UidGid:     "1025:1025",
						},
					},

					Bin:            binary,
					Bech32Prefix:   "nillion",
					Denom:          "nillion",
					GasPrices:      "0.00nillion",
					GasAdjustment:  1.3,
					TrustingPeriod: "508h",
					NoHostMount:    false,
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

func (s *NetworkTestSuite) UpgradeChain(ctx context.Context, chain *cosmos.CosmosChain, wallet *ibc.Wallet, planName, upgradeVersion string) {
	//prevVersion := chain.Nodes()[0].Image.Version
	//plan := upgradetypes.Plan{
	//	Name:   planName,
	//	Height: int64(haltHeight),
	//	Info:   fmt.Sprintf("upgrade version test from %s to %s", prevVersion, upgradeVersion),
	//}
	//upgradeProposal := upgradetypes.NewSoftwareUpgradeProposal(fmt.Sprintf("upgrade from %s to %s", prevVersion, upgradeVersion), "upgrade chain E2E test", plan)
	//s.ExecuteGovProposal(ctx, chain, wallet, upgradeProposal)
	//
	//height, err := chain.Height(ctx)
	//s.Require().NoError(err, "error fetching height before upgrade")
	//
	//fmt.Printf("height: %d\n", height)
	//
	//timeoutCtx, timeoutCtxCancel := context.WithTimeout(ctx, time.Minute*2)
	//defer timeoutCtxCancel()
	//
	//err = test.WaitForBlocks(timeoutCtx, int(haltHeight-height)+1, chain)
	//s.Require().Error(err, "chain did not halt at halt height")
	//
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
