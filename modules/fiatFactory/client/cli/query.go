package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/persistenceSDK/modules/fiatFactory/internal/keeper"
	fiatFactoryTypes "github.com/persistenceOne/persistenceSDK/modules/fiatFactory/internal/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

func QueryFiatCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pegHash [pegHash]",
		Short: "Query fiat peg",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			pegHash := args[0]

			ctx := context.NewCLIContext()
			pegHashHex, err := types.GetFiatPegHashHex(pegHash)
			if err != nil {
				return err
			}

			res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", fiatFactoryTypes.QuerierRoute, keeper.QueryFiat, fiatFactoryTypes.PegHashKey), pegHashHex)
			if err != nil {
				return err
			}

			if res == nil {
				return cTypes.ErrUnknownAddress("No fiat with pegHash " + pegHash +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}
			fmt.Println(res)
			return nil
		},
	}
}
