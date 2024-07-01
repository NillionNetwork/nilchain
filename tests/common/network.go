package common

import (
	"context"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/require"
	testifysuite "github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

const (
	chainRepository = "ghcr.io/nillionnetwork/nilliond"
)

type NetworkTestSuite struct {
	testifysuite.Suite
}

func (s *NetworkTestSuite) InitChain(version string) {
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

					Bin:            "nilchaind",
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
}
