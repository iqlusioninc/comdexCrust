package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commitHub/commitBlockchain/codec"

	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/auth/client/utils"

	"github.com/commitHub/commitBlockchain/types"

	fiatFactoryTypes "github.com/commitHub/commitBlockchain/modules/fiatFactory/internal/types"
)

func IssueFiatCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Initializes fiat with the given details and issues to the given address",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			toStr := viper.GetString(FlagTo)
			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := types.GetFiatPegHashHex(pegHashStr)
			transactionIDStr := viper.GetString(FlagTransactionID)
			transactionAmountInt64 := viper.GetInt64(FlagTransactionAmount)

			fiatPeg := types.BaseFiatPeg{
				PegHash:           pegHashHex,
				TransactionID:     transactionIDStr,
				TransactionAmount: transactionAmountInt64,
			}
			fiatPegI := types.ToFiatPeg(fiatPeg)

			msg := fiatFactoryTypes.BuildIssueFiatMsg(cliCtx.GetFromAddress(), to, fiatPegI)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsTransactionID)
	cmd.Flags().AddFlagSet(fsTransactionAmount)

	return cmd
}
