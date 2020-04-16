package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	ibcclient "github.com/persistenceOne/comdexCrust/modules/ibc/02-client"
	connection "github.com/persistenceOne/comdexCrust/modules/ibc/03-connection"
	channel "github.com/persistenceOne/comdexCrust/modules/ibc/04-channel"
	transfer "github.com/persistenceOne/comdexCrust/modules/ibc/20-transfer/client/cli"
	"github.com/persistenceOne/comdexCrust/modules/ibc/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	ibcTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "IBC transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ibcTxCmd.AddCommand(client.PostCommands(
		ibcclient.GetTxCmd(cdc, storeKey),
		connection.GetTxCmd(cdc, storeKey),
		channel.GetTxCmd(cdc, storeKey),
		transfer.GetTxCmd(cdc),
	)...)
	return ibcTxCmd
}

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group ibc queries under a subcommand
	ibcQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the IBC module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ibcQueryCmd.AddCommand(client.GetCommands(
		ibcclient.GetQueryCmd(cdc, queryRoute),
		connection.GetQueryCmd(cdc, queryRoute),
		channel.GetQueryCmd(cdc, queryRoute),
		transfer.GetQueryCmd(cdc, queryRoute),
	)...)
	return ibcQueryCmd
}
