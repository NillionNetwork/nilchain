package cli_test

import (
	"fmt"
	"github.com/NillionNetwork/nillion-chain/x/meta/client/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/stretchr/testify/require"
	"io"
	"testing"

	"github.com/NillionNetwork/nillion-chain/x/meta"
	abci "github.com/cometbft/cometbft/abci/types"
	rpcclientmock "github.com/cometbft/cometbft/rpc/client/mock"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testutilmod "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/suite"
)

type CLITestSuite struct {
	suite.Suite

	kr        keyring.Keyring
	encCfg    testutilmod.TestEncodingConfig
	baseCtx   client.Context
	clientCtx client.Context
}

func TestCLITestSuite(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}

func (s *CLITestSuite) SetupSuite() {
	s.encCfg = testutilmod.MakeTestEncodingConfig(meta.AppModule{})
	s.kr = keyring.NewInMemory(s.encCfg.Codec)
	s.baseCtx = client.Context{}.
		WithKeyring(s.kr).
		WithTxConfig(s.encCfg.TxConfig).
		WithCodec(s.encCfg.Codec).
		WithClient(clitestutil.MockCometRPC{Client: rpcclientmock.Client{}}).
		WithAccountRetriever(client.MockAccountRetriever{}).
		WithOutput(io.Discard).
		WithChainID("test-chain")

	ctxGen := func() client.Context {
		bz, _ := s.encCfg.Codec.Marshal(&sdk.TxResponse{})
		c := clitestutil.NewMockCometRPC(abci.ResponseQuery{
			Value: bz,
		})
		return s.baseCtx.WithClient(c)
	}
	s.clientCtx = ctxGen()
}

func (s *CLITestSuite) TestCmdPayFor() {
	val := testutil.CreateKeyringAccounts(s.T(), s.kr, 1)

	someResource := `{
		"resource": "some resource"
	}`
	someResourceFile := testutil.WriteToNewTempFile(s.T(), someResource)
	defer someResourceFile.Close()

	cmd := cli.CmdPayFor()
	out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, []string{
		val[0].Address.String(),
		"1000infinity",
		someResourceFile.Name(),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	})
	require.NoError(s.T(), err)
	msg := &sdk.TxResponse{}
	s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), msg), out.String())
}
