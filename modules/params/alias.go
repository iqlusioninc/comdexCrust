// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/persistenceOne/comdexCrust/modules/params/subspace
// ALIASGEN: github.com/persistenceOne/comdexCrust/modules/params/types
package params

import (
	"github.com/persistenceOne/comdexCrust/modules/params/subspace"
	"github.com/persistenceOne/comdexCrust/modules/params/types"
)

const (
	StoreKey             = subspace.StoreKey
	TStoreKey            = subspace.TStoreKey
	DefaultCodespace     = types.DefaultCodespace
	CodeUnknownSubspace  = types.CodeUnknownSubspace
	CodeSettingParameter = types.CodeSettingParameter
	CodeEmptyData        = types.CodeEmptyData
	ModuleName           = types.ModuleName
	RouterKey            = types.RouterKey
	ProposalTypeChange   = types.ProposalTypeChange
)

var (
	// functions aliases
	NewSubspace                = subspace.NewSubspace
	NewKeyTable                = subspace.NewKeyTable
	RegisterCodec              = types.RegisterCodec
	ErrUnknownSubspace         = types.ErrUnknownSubspace
	ErrSettingParameter        = types.ErrSettingParameter
	ErrEmptyChanges            = types.ErrEmptyChanges
	ErrEmptySubspace           = types.ErrEmptySubspace
	ErrEmptyKey                = types.ErrEmptyKey
	ErrEmptyValue              = types.ErrEmptyValue
	NewParameterChangeProposal = types.NewParameterChangeProposal
	NewParamChange             = types.NewParamChange
	NewParamChangeWithSubkey   = types.NewParamChangeWithSubkey
	ValidateChanges            = types.ValidateChanges
)

type (
	ParamSetPair            = subspace.ParamSetPair
	ParamSetPairs           = subspace.ParamSetPairs
	ParamSet                = subspace.ParamSet
	Subspace                = subspace.Subspace
	ReadOnlySubspace        = subspace.ReadOnlySubspace
	KeyTable                = subspace.KeyTable
	ParameterChangeProposal = types.ParameterChangeProposal
	ParamChange             = types.ParamChange
)
