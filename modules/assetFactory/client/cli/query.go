package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/commitHub/commitBlockchain/types"

	assetFactoryTypes "github.com/commitHub/commitBlockchain/modules/assetFactory/internal/types"
)

func QueryAssetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "[pegHash]",
		Short: "Query asset transaction details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			pegHash := args[0]

			ctx := context.NewCLIContext()
			pegHashHex, err := types.GetAssetPegHashHex(pegHash)
			if err != nil {
				return err
			}

			res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/%s", assetFactoryTypes.QuerierRoute,
				assetFactoryTypes.PegHashKey), pegHashHex)
			if err != nil {
				return err
			}

			if res == nil {
				return cTypes.ErrUnknownAddress("No asset with pegHash " + pegHash +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}

			fmt.Println(res)
			return nil
		},
	}
}
