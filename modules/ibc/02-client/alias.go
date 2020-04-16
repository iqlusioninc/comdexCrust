package client

// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/persistenceOne/comdexCrust/modules/ibc/02-client/keeper
// ALIASGEN: github.com/persistenceOne/comdexCrust/modules/ibc/02-client/types

import (
	"github.com/persistenceOne/comdexCrust/modules/ibc/02-client/keeper"
	"github.com/persistenceOne/comdexCrust/modules/ibc/02-client/types"
)

const (
	DefaultCodespace           = types.DefaultCodespace
	CodeClientExists           = types.CodeClientExists
	CodeClientNotFound         = types.CodeClientNotFound
	CodeClientFrozen           = types.CodeClientFrozen
	CodeConsensusStateNotFound = types.CodeConsensusStateNotFound
	CodeInvalidConsensusState  = types.CodeInvalidConsensusState
	CodeClientTypeNotFound     = types.CodeClientTypeNotFound
	CodeInvalidClientType      = types.CodeInvalidClientType
	CodeRootNotFound           = types.CodeRootNotFound
	CodeInvalidHeader          = types.CodeInvalidHeader
	CodeInvalidEvidence        = types.CodeInvalidEvidence
	AttributeKeyClientID       = types.AttributeKeyClientID
	SubModuleName              = types.SubModuleName
	StoreKey                   = types.StoreKey
	RouterKey                  = types.RouterKey
	QuerierRoute               = types.QuerierRoute
	QueryClientState           = types.QueryClientState
	QueryConsensusState        = types.QueryConsensusState
	QueryVerifiedRoot          = types.QueryVerifiedRoot
)

var (
	// functions aliases
	NewKeeper                    = keeper.NewKeeper
	QuerierClientState           = keeper.QuerierClientState
	QuerierConsensusState        = keeper.QuerierConsensusState
	QuerierVerifiedRoot          = keeper.QuerierVerifiedRoot
	RegisterCodec                = types.RegisterCodec
	ErrClientExists              = types.ErrClientExists
	ErrClientNotFound            = types.ErrClientNotFound
	ErrClientFrozen              = types.ErrClientFrozen
	ErrConsensusStateNotFound    = types.ErrConsensusStateNotFound
	ErrInvalidConsensus          = types.ErrInvalidConsensus
	ErrClientTypeNotFound        = types.ErrClientTypeNotFound
	ErrInvalidClientType         = types.ErrInvalidClientType
	ErrRootNotFound              = types.ErrRootNotFound
	ErrInvalidHeader             = types.ErrInvalidHeader
	ErrInvalidEvidence           = types.ErrInvalidEvidence
	ClientStatePath              = types.ClientStatePath
	ClientTypePath               = types.ClientTypePath
	ConsensusStatePath           = types.ConsensusStatePath
	RootPath                     = types.RootPath
	KeyClientState               = types.KeyClientState
	KeyClientType                = types.KeyClientType
	KeyConsensusState            = types.KeyConsensusState
	KeyRoot                      = types.KeyRoot
	NewMsgCreateClient           = types.NewMsgCreateClient
	NewMsgUpdateClient           = types.NewMsgUpdateClient
	NewMsgSubmitMisbehaviour     = types.NewMsgSubmitMisbehaviour
	NewQueryClientStateParams    = types.NewQueryClientStateParams
	NewQueryCommitmentRootParams = types.NewQueryCommitmentRootParams
	NewClientState               = types.NewClientState

	// variable aliases
	SubModuleCdc                = types.SubModuleCdc
	EventTypeCreateClient       = types.EventTypeCreateClient
	EventTypeUpdateClient       = types.EventTypeUpdateClient
	EventTypeSubmitMisbehaviour = types.EventTypeSubmitMisbehaviour
	AttributeValueCategory      = types.AttributeValueCategory
)

type (
	Keeper                    = keeper.Keeper
	MsgCreateClient           = types.MsgCreateClient
	MsgUpdateClient           = types.MsgUpdateClient
	MsgSubmitMisbehaviour     = types.MsgSubmitMisbehaviour
	QueryClientStateParams    = types.QueryClientStateParams
	QueryCommitmentRootParams = types.QueryCommitmentRootParams
	State                     = types.State
)
