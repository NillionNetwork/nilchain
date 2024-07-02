package common

import (
	"context"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"time"
)

const VotingPeriod time.Duration = time.Second * 30

// ExecuteGovProposal submits the given governance proposal using the provided user and uses all validators to vote yes on the proposal.
// It ensure the proposal successfully passes.
func (s *NetworkTestSuite) ExecuteGovProposal(ctx context.Context, chain *cosmos.CosmosChain, wallet ibc.Wallet, planName string, upgradeVersion string) {
	address, err := chain.GetAddress(ctx, wallet.KeyName())
	s.Require().NoError(err)

	addrString, err := sdk.Bech32ifyAddressBytes(chain.Config().Bech32Prefix, address)
	s.Require().NoError(err)

	prevVersion := chain.Nodes()[0].Image.Version
	plan := upgradetypes.Plan{
		Name:   planName,
		Height: int64(haltHeight),
		Info:   fmt.Sprintf("upgrade version test from %s to %s", prevVersion, upgradeVersion),
	}
	upgrade := upgradetypes.MsgSoftwareUpgrade{
		Authority: authority,
		Plan:      plan,
	}

	proposal, err := chain.BuildProposal(
		[]cosmos.ProtoMessage{&upgrade},
		fmt.Sprintf("upgrade from %s to %s", prevVersion, upgradeVersion),
		"upgrade chain E2E test",
		"",
		"50000000unillion",
		addrString,
		true,
	)
	s.Require().NoError(err)

	_, err = chain.SubmitProposal(ctx, wallet.KeyName(), proposal)
	s.Require().NoError(err)

	prop, err := chain.GovQueryProposal(ctx, 1)
	s.Require().NoError(err)
	s.Require().Equal(v1beta1.StatusVotingPeriod, prop.Status)

	err = chain.VoteOnProposalAllValidators(ctx, 1, cosmos.ProposalVoteYes)
	s.Require().NoError(err)

	// ensure voting period has not passed before validators finished voting
	prop, err = chain.GovQueryProposal(ctx, 1)
	s.Require().NoError(err)
	s.Require().Equal(v1beta1.StatusVotingPeriod, prop.Status)

	time.Sleep(VotingPeriod) // pass proposal

	prop, err = chain.GovQueryProposal(ctx, 1)
	s.Require().NoError(err)
	s.Require().Equal(v1beta1.StatusPassed, prop.Status)
}

// BroadcastMessages broadcasts the provided messages to the given chain and signs them on behalf of the provided user.
// Once the broadcast response is returned, we wait for a few blocks to be created on both chain A and chain B.
func (s *NetworkTestSuite) BroadcastMessages(ctx context.Context, chain *cosmos.CosmosChain, user cosmos.User, msgs ...sdk.Msg) (sdk.TxResponse, error) {
	broadcaster := cosmos.NewBroadcaster(s.T(), chain)
	resp, err := cosmos.BroadcastTx(ctx, broadcaster, user, msgs...)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	//chainA, chainB := s.GetChains()
	//err = test.WaitForBlocks(ctx, 2, chainA, chainB)
	return resp, err
}
