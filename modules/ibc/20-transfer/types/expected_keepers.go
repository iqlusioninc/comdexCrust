package types

import (
	clientexported "github.com/commitHub/commitBlockchain/modules/ibc/02-client/exported"
	connection "github.com/commitHub/commitBlockchain/modules/ibc/03-connection"
	channel "github.com/commitHub/commitBlockchain/modules/ibc/04-channel"
	channelexported "github.com/commitHub/commitBlockchain/modules/ibc/04-channel/exported"
	commitment "github.com/commitHub/commitBlockchain/modules/ibc/23-commitment"
	supplyexported "github.com/commitHub/commitBlockchain/modules/supply/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitTypes "github.com/commitHub/commitBlockchain/types"
)

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
}

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channel.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	SendPacket(ctx sdk.Context, packet channelexported.PacketI, portCapability commitTypes.CapabilityKey) error
	RecvPacket(ctx sdk.Context, packet channelexported.PacketI, proof commitment.ProofI, proofHeight uint64, acknowledgement []byte, portCapability commitTypes.CapabilityKey) (channelexported.PacketI, error)
}

// ClientKeeper defines the expected IBC client keeper
type ClientKeeper interface {
	GetConsensusState(ctx sdk.Context, clientID string) (connection clientexported.ConsensusState, found bool)
}

// ConnectionKeeper defines the expected IBC connection keeper
type ConnectionKeeper interface {
	GetConnection(ctx sdk.Context, connectionID string) (connection connection.ConnectionEnd, found bool)
}

// SupplyKeeper expected supply keeper
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
}
