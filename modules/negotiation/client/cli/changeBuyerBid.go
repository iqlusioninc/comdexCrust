package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/persistenceOne/comdexCrust/codec"
	"github.com/persistenceOne/comdexCrust/modules/auth"
	"github.com/persistenceOne/comdexCrust/modules/auth/client/utils"
	negotiationTypes "github.com/persistenceOne/comdexCrust/modules/negotiation/internal/types"
	"github.com/persistenceOne/comdexCrust/types"
)

func ChangeBuyerBidCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-buyer-bid",
		Short: "Change or create a buyer negotiation bid",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			toStr := viper.GetString(FlagTo)

			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			bid := viper.GetInt64(FlagBid)
			time := viper.GetInt64(FlagTime)
			hashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := types.GetAssetPegHashHex(hashStr)
			negotiationID := types.NegotiationID(append(append(cliCtx.GetFromAddress().Bytes(), to.Bytes()...), pegHashHex...))

			proposedNegotiation := &types.BaseNegotiation{
				NegotiationID: negotiationID,
				BuyerAddress:  cliCtx.GetFromAddress(),
				SellerAddress: to,
				PegHash:       pegHashHex,
				Bid:           bid,
				Time:          time,
			}

			msg := negotiationTypes.BuildMsgChangeBuyerBid(proposedNegotiation)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsBid)
	cmd.Flags().AddFlagSet(fsTime)
	return cmd
}
