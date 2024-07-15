package upgrades

import (
	"context"
	"cosmossdk.io/math"
	nillionapp "github.com/NillionNetwork/nilchain/app"
	"github.com/NillionNetwork/nilchain/app/upgrades"
	"github.com/NillionNetwork/nilchain/tests/common"
	"github.com/strangelove-ventures/interchaintest/v8"
	"testing"

	testifysuite "github.com/stretchr/testify/suite"
)

func init() {
	nillionapp.SetBech32AddressPrefixes()
}

func TestUpgradeTestSuite(t *testing.T) {
	testifysuite.Run(t, new(UpgradeTestSuite))
}

type UpgradeTestSuite struct {
	common.NetworkTestSuite
}

// TestUpgrade0_2_1 tests the upgrade from 0.2.1 to 0.2.4
func (s *UpgradeTestSuite) TestUpgrade0_2_4() {
	oldVersion := "v0.2.1-test-only"
	newVersion := upgrades.Upgrade_0_2_4.UpgradeName

	s.InitChain(oldVersion, "nilliond")

	ctx := context.Background()

	users := interchaintest.GetAndFundTestUsers(s.T(), ctx, "nillion", math.NewInt(1000000000), s.Chain)

	s.UpgradeChain(ctx, s.Chain, users[0], newVersion)
}
