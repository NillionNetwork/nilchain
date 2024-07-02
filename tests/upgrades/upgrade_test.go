package upgrades

import (
	"context"
	"cosmossdk.io/math"
	"github.com/NillionNetwork/nilchain/tests/common"
	"github.com/strangelove-ventures/interchaintest/v8"
	"testing"

	testifysuite "github.com/stretchr/testify/suite"
)

func TestUpgradeTestSuite(t *testing.T) {
	testifysuite.Run(t, new(UpgradeTestSuite))
}

type UpgradeTestSuite struct {
	common.NetworkTestSuite
}

// TestUpgrade0_2_1 tests the upgrade from 0.1.0 to 0.2.1
func (s *UpgradeTestSuite) TestUpgrade0_2_1() {
	oldVersion := "v0.1.1"
	//newVersion := "0.2.1"

	s.InitChain(oldVersion, "nilliond")

	ctx := context.Background()

	users := interchaintest.GetAndFundTestUsers(s.T(), ctx, "nillion", math.NewInt(1000000000))

	s.UpgradeChain(ctx, s.Chain, &users[0], "upgrade-test", "v0.2.1")
}
