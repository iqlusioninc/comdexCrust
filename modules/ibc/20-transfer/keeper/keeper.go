package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/comdexCrust/modules/acl"
	"github.com/persistenceOne/comdexCrust/modules/assetFactory"
	"github.com/persistenceOne/comdexCrust/modules/ibc/20-transfer/types"
	ibctypes "github.com/persistenceOne/comdexCrust/modules/ibc/types"
	supplyexported "github.com/persistenceOne/comdexCrust/modules/supply/exported"
)

// DefaultPacketTimeout is the default packet timeout relative to the current block height
const (
	DefaultPacketTimeout = 1000 // NOTE: in blocks
)

// Keeper defines the IBC transfer keeper
type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *codec.Codec
	codespace sdk.CodespaceType
	prefix    []byte // prefix bytes for accessing the store

	clientKeeper     types.ClientKeeper
	connectionKeeper types.ConnectionKeeper
	channelKeeper    types.ChannelKeeper
	bankKeeper       types.BankKeeper
	supplyKeeper     types.SupplyKeeper
	assetKeeper      assetFactory.Keeper
	aclKeeper        acl.Keeper
}

// NewKeeper creates a new IBC transfer Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType,
	clientKeeper types.ClientKeeper, connnectionKeeper types.ConnectionKeeper,
	channelKeeper types.ChannelKeeper, bankKeeper types.BankKeeper,
	supplyKeeper types.SupplyKeeper, assetKeeper assetFactory.Keeper, aclKeeper acl.Keeper,
) Keeper {

	// ensure ibc transfer module account is set
	if addr := supplyKeeper.GetModuleAddress(types.GetModuleAccountName()); addr == nil {
		panic("the IBC transfer module account has not been set")
	}

	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		codespace:        sdk.CodespaceType(fmt.Sprintf("%s/%s", codespace, types.DefaultCodespace)), // "ibc/transfer",
		prefix:           []byte(types.SubModuleName + "/"),                                          // "transfer/"
		clientKeeper:     clientKeeper,
		connectionKeeper: connnectionKeeper,
		channelKeeper:    channelKeeper,
		bankKeeper:       bankKeeper,
		supplyKeeper:     supplyKeeper,
		assetKeeper:      assetKeeper,
		aclKeeper:        aclKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s/%s", ibctypes.ModuleName, types.SubModuleName))
}

// GetTransferAccount returns the ICS20 - transfers ModuleAccount
func (k Keeper) GetTransferAccount(ctx sdk.Context) supplyexported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.GetModuleAccountName())
}
