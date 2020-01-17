package types

import (
	"github.com/persistenceOne/persistenceSDK/types"
)

type GenesisState struct {
	AssetPegs []types.AssetPeg `json:"assetPegs"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
