package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/modules/auth/exported"
)

type AccountKeeper interface {
	GetAccount(ctx cTypes.Context, addr cTypes.AccAddress) exported.Account
}
