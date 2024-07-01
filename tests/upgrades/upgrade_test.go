package upgrades

import (
	"github.com/NillionNetwork/nilchain/tests/common"
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

	s.InitChain(oldVersion)
}
