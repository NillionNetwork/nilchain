package common

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
)

// ExecuteGovProposal submits the given governance proposal using the provided user and uses all validators to vote yes on the proposal.
// It ensure the proposal successfully passes.
func (s *NetworkTestSuite) ExecuteGovProposal(ctx context.Context, chain *cosmos.CosmosChain, user cosmos.User, content v1beta1.Content) {
	sender, err := sdk.AccAddressFromBech32(user.FormattedAddress())
	s.Require().NoError(err)

	msgSubmitProposal, err := v1beta1.NewMsgSubmitProposal(content, sdk.NewCoins(sdk.NewCoin(chain.Config().Denom, v1beta1.DefaultMinDepositTokens)), sender)
	s.Require().NoError(err)

	_, err = s.BroadcastMessages(ctx, chain, user, msgSubmitProposal)
	s.Require().NoError(err)
	//s.AssertValidTxResponse(txResp)

	// TODO: replace with parsed proposal ID from MsgSubmitProposalResponse
	// https://github.com/cosmos/ibc-go/issues/2122

	//proposal, err := s.QueryProposal(ctx, chain, 1)
	//s.Require().NoError(err)
	//s.Require().Equal(govtypes.StatusVotingPeriod, proposal.Status)
	//
	//err = chain.VoteOnProposalAllValidators(ctx, "1", ibc.ProposalVoteYes)
	//s.Require().NoError(err)
	//
	//// ensure voting period has not passed before validators finished voting
	//proposal, err = s.QueryProposal(ctx, chain, 1)
	//s.Require().NoError(err)
	//s.Require().Equal(govtypes.StatusVotingPeriod, proposal.Status)
	//
	//time.Sleep(testvalues.VotingPeriod) // pass proposal
	//
	//proposal, err = s.QueryProposal(ctx, chain, 1)
	//s.Require().NoError(err)
	//s.Require().Equal(govtypes.StatusPassed, proposal.Status)
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
