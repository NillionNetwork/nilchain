package cli

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/NillionNetwork/nilliond/x/meta/types"
)

func TxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Meta transactions subcommands",
	}

	cmd.AddCommand(
		CmdPayFor(),
	)

	return cmd
}

// CmdPayFor returns the command to pay for a resource.
func CmdPayFor() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pay-for [from_address] [amount] [resource-file-data]",
		Short: "Pay for a resource",
		Long: `
Pay for a resource by sending coins to the module account and burning them.
Usage:
$ nillion-chaind tx meta pay-for [from_address] [amount] [resource-file-data]

Where:
- [from_address] is the address of the sender.
- [amount] is the amount of coins to send.
- [resource] is the resource to pay for.
Example:
$ nillion-chaind tx meta pay-for infinity1kc068s88tkyjcc0lkx67x95dwc7hrfm44u8k55 1000infinity resource1.json

Where resource1.json contains the resource data.
`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Flags().Set(flags.FlagFrom, args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			metadata, err := parseMetadataFile(args[2])
			if err != nil {
				return err
			}

			msg := &types.MsgPayFor{
				FromAddress: clientCtx.GetFromAddress().String(),
				Amount:      coins,
				Resource:    metadata,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func parseMetadataFile(file string) ([]byte, error) {
	metadata, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}
