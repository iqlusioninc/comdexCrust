package types

import (
	"fmt"

	ibctypes "github.com/commitHub/commitBlockchain/modules/ibc/types"
)

// IBC transfer events
const (
	AttributeKeyReceiver = "receiver"
)

// IBC transfer events vars
var (
	AttributeValueCategory = fmt.Sprintf("%s_%s", ibctypes.ModuleName, SubModuleName)
)
