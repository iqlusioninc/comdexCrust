package types

import (
	"encoding/json"
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceSDK/modules/acl"
	"github.com/persistenceOne/persistenceSDK/types"
)

// RouterKey is they name of the bank module
const RouterKey = ModuleName

// MsgSend - high level transaction of the coin module
type MsgSend struct {
	FromAddress cTypes.AccAddress `json:"from_address" yaml:"from_address"`
	ToAddress   cTypes.AccAddress `json:"to_address" yaml:"to_address"`
	Amount      cTypes.Coins      `json:"amount" yaml:"amount"`
}

var _ cTypes.Msg = MsgSend{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(fromAddr, toAddr cTypes.AccAddress, amount cTypes.Coins) MsgSend {
	return MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: amount}
}

// Route Implements Msg.
func (msg MsgSend) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSend) Type() string { return "send" }

// ValidateBasic Implements Msg.
func (msg MsgSend) ValidateBasic() cTypes.Error {
	if msg.FromAddress.Empty() {
		return cTypes.ErrInvalidAddress("missing sender address")
	}
	if msg.ToAddress.Empty() {
		return cTypes.ErrInvalidAddress("missing recipient address")
	}
	if !msg.Amount.IsValid() {
		return cTypes.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}
	if !msg.Amount.IsAllPositive() {
		return cTypes.ErrInsufficientCoins("send amount must be positive")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
	return cTypes.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSend) GetSigners() []cTypes.AccAddress {
	return []cTypes.AccAddress{msg.FromAddress}
}

// MsgMultiSend - high level transaction of the coin module
type MsgMultiSend struct {
	Inputs  []Input  `json:"inputs" yaml:"inputs"`
	Outputs []Output `json:"outputs" yaml:"outputs"`
}

var _ cTypes.Msg = MsgMultiSend{}

// NewMsgMultiSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgMultiSend(in []Input, out []Output) MsgMultiSend {
	return MsgMultiSend{Inputs: in, Outputs: out}
}

// Route Implements Msg
func (msg MsgMultiSend) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgMultiSend) Type() string { return "multisend" }

// ValidateBasic Implements Msg.
func (msg MsgMultiSend) ValidateBasic() cTypes.Error {
	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside
	if len(msg.Inputs) == 0 {
		return ErrNoInputs(DefaultCodespace).TraceSDK("")
	}
	if len(msg.Outputs) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}

	return ValidateInputsOutputs(msg.Inputs, msg.Outputs)
}

// GetSignBytes Implements Msg.
func (msg MsgMultiSend) GetSignBytes() []byte {
	return cTypes.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgMultiSend) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	return addrs
}

// Input models transaction input
type Input struct {
	Address cTypes.AccAddress `json:"address" yaml:"address"`
	Coins   cTypes.Coins      `json:"coins" yaml:"coins"`
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() cTypes.Error {
	if len(in.Address) == 0 {
		return cTypes.ErrInvalidAddress(in.Address.String())
	}
	if !in.Coins.IsValid() {
		return cTypes.ErrInvalidCoins(in.Coins.String())
	}
	if !in.Coins.IsAllPositive() {
		return cTypes.ErrInvalidCoins(in.Coins.String())
	}
	return nil
}

// NewInput - create a transaction input, used with MsgMultiSend
func NewInput(addr cTypes.AccAddress, coins cTypes.Coins) Input {
	return Input{
		Address: addr,
		Coins:   coins,
	}
}

// Output models transaction outputs
type Output struct {
	Address cTypes.AccAddress `json:"address" yaml:"address"`
	Coins   cTypes.Coins      `json:"coins" yaml:"coins"`
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() cTypes.Error {
	if len(out.Address) == 0 {
		return cTypes.ErrInvalidAddress(out.Address.String())
	}
	if !out.Coins.IsValid() {
		return cTypes.ErrInvalidCoins(out.Coins.String())
	}
	if !out.Coins.IsAllPositive() {
		return cTypes.ErrInvalidCoins(out.Coins.String())
	}
	return nil
}

// NewOutput - create a transaction output, used with MsgMultiSend
func NewOutput(addr cTypes.AccAddress, coins cTypes.Coins) Output {
	return Output{
		Address: addr,
		Coins:   coins,
	}
}

// ValidateInputsOutputs validates that each respective input and output is
// valid and that the sum of inputs is equal to the sum of outputs.
func ValidateInputsOutputs(inputs []Input, outputs []Output) cTypes.Error {
	var totalIn, totalOut cTypes.Coins

	for _, in := range inputs {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
		totalIn = totalIn.Add(in.Coins)
	}

	for _, out := range outputs {
		if err := out.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
		totalOut = totalOut.Add(out.Coins)
	}

	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return ErrInputOutputMismatch(DefaultCodespace)
	}

	return nil
}

// *****Comdex

// *****IssueAsset

// IssueAsset - transaction input
type IssueAsset struct {
	IssuerAddress cTypes.AccAddress `json:"issuerAddress"`
	ToAddress     cTypes.AccAddress `json:"toAddress"`
	AssetPeg      types.AssetPeg    `json:"assetPeg"`
}

// NewIssueAsset : initializer
func NewIssueAsset(issuerAddress cTypes.AccAddress, toAddress cTypes.AccAddress, assetPeg types.AssetPeg) IssueAsset {
	return IssueAsset{issuerAddress, toAddress, assetPeg}
}

// GetSignBytes : get bytes to sign
func (in IssueAsset) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		IssuerAddress string         `json:"issuerAddress"`
		ToAddress     string         `json:"toAddress"`
		AssetPeg      types.AssetPeg `json:"assetPeg"`
	}{
		IssuerAddress: in.IssuerAddress.String(),
		ToAddress:     in.ToAddress.String(),
		AssetPeg:      in.AssetPeg,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

func (in IssueAsset) ValidateBasic() cTypes.Error {
	if len(in.IssuerAddress) == 0 {
		return cTypes.ErrInvalidAddress(fmt.Sprintf("invalid Issuer address %s", in.IssuerAddress.String()))
	} else if len(in.ToAddress) == 0 {
		return cTypes.ErrInvalidAddress(fmt.Sprintf("invalid To address %s", in.ToAddress.String()))
	} else if in.AssetPeg.GetAssetPrice() < 0 {
		return ErrNegativeAmount(DefaultCodespace, "Asset price should be grater than 0.")
	} else if in.AssetPeg.GetAssetQuantity() < 0 {
		return ErrNegativeAmount(DefaultCodespace, "Asset quantity should be grater than 0.")
	} else if in.AssetPeg.GetAssetType() == "" {
		return cTypes.ErrUnknownRequest("asset type should not be empty")
	} else if in.AssetPeg.GetDocumentHash() == "" {
		return cTypes.ErrUnknownRequest("DocumentHash should not be empty")
	}
	return nil
}

// #####IssueAsset

// *****MsgBankIssueAssets

// MsgBankIssueAssets : high level issuance of assets module
type MsgBankIssueAssets struct {
	IssueAssets []IssueAsset `json:"issueAssets"`
}

// NewMsgBankIssueAssets : initilizer
func NewMsgBankIssueAssets(issueAssets []IssueAsset) MsgBankIssueAssets {
	return MsgBankIssueAssets{issueAssets}
}

// ***** Implementing cTypes.Msg

var _ cTypes.Msg = MsgBankIssueAssets{}

// Type : implements msg
func (msg MsgBankIssueAssets) Type() string { return "bank" }

func (msg MsgBankIssueAssets) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgBankIssueAssets) ValidateBasic() cTypes.Error {
	if len(msg.IssueAssets) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.IssueAssets {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankIssueAssets) GetSignBytes() []byte {
	var issueAssets []json.RawMessage
	for _, issueAsset := range msg.IssueAssets {
		issueAssets = append(issueAssets, issueAsset.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		IssueAssets []json.RawMessage `json:"issueAssets"`
	}{
		IssueAssets: issueAssets,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgBankIssueAssets) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.IssueAssets))
	for i, in := range msg.IssueAssets {
		addrs[i] = in.IssuerAddress
	}
	return addrs
}

// ##### Implement cTypes.Msg

// #####MsgBankIssueAssets

// ****RedeemAsset

// RedeemAsset : transsction input
type RedeemAsset struct {
	IssuerAddress   cTypes.AccAddress `json:"issuerAddress"`
	RedeemerAddress cTypes.AccAddress `json:"redeemerAddress"`
	PegHash         types.PegHash     `json:"pegHash"`
}

// NewRedeemAsset : initializer
func NewRedeemAsset(issuerAddress cTypes.AccAddress, redeemerAddress cTypes.AccAddress, pegHash types.PegHash) RedeemAsset {
	return RedeemAsset{issuerAddress, redeemerAddress, pegHash}
}

// GetSignBytes : get bytes to sign
func (in RedeemAsset) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		IssuerAddress   string        `json:"issuerAddress"`
		RedeemerAddress string        `json:"redeemerAddress"`
		PegHash         types.PegHash `json:"pegHash"`
	}{
		IssuerAddress:   in.IssuerAddress.String(),
		RedeemerAddress: in.RedeemerAddress.String(),
		PegHash:         in.PegHash,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

func (in RedeemAsset) ValidateBasic() cTypes.Error {
	if len(in.IssuerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.IssuerAddress.String())
	} else if len(in.RedeemerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.RedeemerAddress.String())
	} else if len(in.PegHash) == 0 {
		return cTypes.ErrUnknownRequest("PegHash should not be empty.")
	}
	return nil
}

// #####RedeemAsset

// *****MsgBankRedeemAssets

// MsgBankRedeemAssets : Message to redeem issued assets
type MsgBankRedeemAssets struct {
	RedeemAssets []RedeemAsset `json:"redeemAssets"`
}

// NewMsgBankRedeemAssets : initializer
func NewMsgBankRedeemAssets(redeemAssets []RedeemAsset) MsgBankRedeemAssets {
	return MsgBankRedeemAssets{redeemAssets}
}

// *****Implementing cTypes.Msg

var _ cTypes.Msg = MsgBankRedeemAssets{}

// Type : implements msg
func (msg MsgBankRedeemAssets) Type() string { return "bank" }

func (msg MsgBankRedeemAssets) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgBankRedeemAssets) ValidateBasic() cTypes.Error {
	if len(msg.RedeemAssets) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.RedeemAssets {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankRedeemAssets) GetSignBytes() []byte {
	var redeemAssets []json.RawMessage
	for _, redeemAsset := range msg.RedeemAssets {
		redeemAssets = append(redeemAssets, redeemAsset.GetSignBytes())
	}

	bz, err := ModuleCdc.MarshalJSON(struct {
		RedeemAssets []json.RawMessage `json:"redeemAssets"`
	}{
		RedeemAssets: redeemAssets,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners : implements msg
func (msg MsgBankRedeemAssets) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.RedeemAssets))
	for i, in := range msg.RedeemAssets {
		addrs[i] = in.RedeemerAddress
	}
	return addrs
}

// ##### Implement cTypes.Msg

// ######MsgBankRedeemAssets

// *****IssueFiat

// IssueFiat - transaction input
type IssueFiat struct {
	IssuerAddress cTypes.AccAddress `json:"issuerAddress"`
	ToAddress     cTypes.AccAddress `json:"toAddress"`
	FiatPeg       types.FiatPeg     `json:"fiatPeg"`
}

// NewIssueFiat : initializer
func NewIssueFiat(issuerAddress cTypes.AccAddress, toAddress cTypes.AccAddress, fiatPeg types.FiatPeg) IssueFiat {
	return IssueFiat{issuerAddress, toAddress, fiatPeg}
}

// GetSignBytes : get bytes to sign
func (in IssueFiat) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		IssuerAddress string        `json:"issuerAddress"`
		ToAddress     string        `json:"toAddress"`
		FiatPeg       types.FiatPeg `json:"fiatPeg"`
	}{
		IssuerAddress: in.IssuerAddress.String(),
		ToAddress:     in.ToAddress.String(),
		FiatPeg:       in.FiatPeg,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

func (in IssueFiat) ValidateBasic() cTypes.Error {
	if len(in.IssuerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.IssuerAddress.String())
	} else if len(in.ToAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.ToAddress.String())
	} else if in.FiatPeg.GetTransactionAmount() < 0 {
		return ErrNegativeAmount(DefaultCodespace, "Transaction amount should be grater than 0.")
	} else if in.FiatPeg.GetTransactionID() == "" {
		return cTypes.ErrUnknownRequest("Transaction should not be empty")
	}
	return nil
}

// #####IssueFiat

// *****MsgBankIssueFiats

// MsgBankIssueFiats : high level issuance of fiats module
type MsgBankIssueFiats struct {
	IssueFiats []IssueFiat `json:"issueFiats"`
}

// NewMsgBankIssueFiats : initilizer
func NewMsgBankIssueFiats(issueFiats []IssueFiat) MsgBankIssueFiats {
	return MsgBankIssueFiats{issueFiats}
}

// ***** Implementing cTypes.Msg

var _ cTypes.Msg = MsgBankIssueFiats{}

// Type : implements msg
func (msg MsgBankIssueFiats) Type() string { return "bank" }

func (msg MsgBankIssueFiats) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgBankIssueFiats) ValidateBasic() cTypes.Error {
	if len(msg.IssueFiats) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.IssueFiats {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankIssueFiats) GetSignBytes() []byte {
	var issueFiats []json.RawMessage
	for _, issueFiat := range msg.IssueFiats {
		issueFiats = append(issueFiats, issueFiat.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		IssueFiats []json.RawMessage `json:"issueFiats"`
	}{
		IssueFiats: issueFiats,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgBankIssueFiats) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.IssueFiats))
	for i, in := range msg.IssueFiats {
		addrs[i] = in.IssuerAddress
	}
	return addrs
}

// ##### Implement cTypes.Msg

// #####MsgBankIssueFiats

// ****RedeemFiat

// RedeemFiat : transaction input
type RedeemFiat struct {
	RedeemerAddress cTypes.AccAddress `json:"redeemerAddress"`
	IssuerAddress   cTypes.AccAddress `json:"issuerAddress"`
	Amount          int64             `json:"amount"`
}

// NewRedeemFiat : initializer
func NewRedeemFiat(redeemerAddress cTypes.AccAddress, issuerAddress cTypes.AccAddress, amount int64) RedeemFiat {
	return RedeemFiat{redeemerAddress, issuerAddress, amount}
}

// GetSignBytes : get bytes to sign
func (in RedeemFiat) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		RedeemerAddress string `json:"redeemerAddress"`
		IssuerAddress   string `json:"issuerAddress"`
		Amount          int64  `json:"amount"`
	}{
		RedeemerAddress: in.RedeemerAddress.String(),
		IssuerAddress:   in.IssuerAddress.String(),
		Amount:          in.Amount,
	})
	if err != nil {
		panic(err)
	}
	return bin
}
func (in RedeemFiat) ValidateBasic() cTypes.Error {
	if len(in.IssuerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.IssuerAddress.String())
	} else if len(in.RedeemerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.RedeemerAddress.String())
	} else if in.Amount <= 0 {
		return cTypes.ErrUnknownRequest("Amount should be Positive")
	}
	return nil
}

// #####RedeemFiat

// *****MsgBankRedeemFiats

// MsgBankRedeemFiats : Message to redeem issued fiats
type MsgBankRedeemFiats struct {
	RedeemFiats []RedeemFiat `json:"redeemFiats"`
}

// NewMsgBankRedeemFiats : initializer
func NewMsgBankRedeemFiats(redeemFiats []RedeemFiat) MsgBankRedeemFiats {
	return MsgBankRedeemFiats{redeemFiats}
}

// *****Implementing cTypes.Msg

var _ cTypes.Msg = MsgBankRedeemFiats{}

// Type : implements msg
func (msg MsgBankRedeemFiats) Type() string { return "bank" }

func (msg MsgBankRedeemFiats) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgBankRedeemFiats) ValidateBasic() cTypes.Error {
	if len(msg.RedeemFiats) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.RedeemFiats {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankRedeemFiats) GetSignBytes() []byte {
	var redeemFiats []json.RawMessage
	for _, redeemFiat := range msg.RedeemFiats {
		redeemFiats = append(redeemFiats, redeemFiat.GetSignBytes())
	}

	bz, err := ModuleCdc.MarshalJSON(struct {
		RedeemFiats []json.RawMessage `json:"redeemFiats"`
	}{
		RedeemFiats: redeemFiats,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners : implements msg
func (msg MsgBankRedeemFiats) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.RedeemFiats))
	for i, in := range msg.RedeemFiats {
		addrs[i] = in.RedeemerAddress
	}
	return addrs
}

// ##### Implement cTypes.Msg

// ######MsgBankRedeemFiats

// *****SendAsset

// SendAsset - transaction input
type SendAsset struct {
	FromAddress cTypes.AccAddress `json:"fromAddress"`
	ToAddress   cTypes.AccAddress `json:"toAddress"`
	PegHash     types.PegHash     `json:"pegHash"`
}

// NewSendAsset : initializer
func NewSendAsset(fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress, pegHash types.PegHash) SendAsset {
	return SendAsset{fromAddress, toAddress, pegHash}
}

// GetSignBytes : get bytes to sign
func (in SendAsset) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		FromAddress string `json:"fromAddress"`
		ToAddress   string `json:"toAddress"`
		PegHash     string `json:"pegHash"`
	}{
		FromAddress: in.FromAddress.String(),
		ToAddress:   in.ToAddress.String(),
		PegHash:     in.PegHash.String(),
	})
	if err != nil {
		panic(err)
	}
	return bin
}

func (in SendAsset) ValidateBasic() cTypes.Error {
	if len(in.FromAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.FromAddress.String())
	} else if len(in.ToAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.ToAddress.String())
	} else if len(in.PegHash) == 0 {
		return cTypes.ErrUnknownRequest("PegHash is empty")
	}
	return nil
}

// #####SendAsset

// *****MsgBankSendAssets

// MsgBankSendAssets : high level issuance of assets module
type MsgBankSendAssets struct {
	SendAssets []SendAsset `json:"sendAssets"`
}

// NewMsgBankSendAssets : initilizer
func NewMsgBankSendAssets(sendAssets []SendAsset) MsgBankSendAssets {
	return MsgBankSendAssets{sendAssets}
}

// ***** Implementing cTypes.Msg

var _ cTypes.Msg = MsgBankSendAssets{}

// Type : implements msg
func (msg MsgBankSendAssets) Type() string { return "bank" }

func (msg MsgBankSendAssets) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgBankSendAssets) ValidateBasic() cTypes.Error {
	if len(msg.SendAssets) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.SendAssets {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankSendAssets) GetSignBytes() []byte {
	var sendAssets []json.RawMessage
	for _, sendAsset := range msg.SendAssets {
		sendAssets = append(sendAssets, sendAsset.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		SendAssets []json.RawMessage `json:"sendAssets"`
	}{
		SendAssets: sendAssets,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgBankSendAssets) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.SendAssets))
	for i, in := range msg.SendAssets {
		addrs[i] = in.FromAddress
	}
	return addrs
}

// ##### Implement cTypes.Msg

// #####MsgBankSendAssets

// *****SendFiat

// SendFiat - transaction input
type SendFiat struct {
	FromAddress cTypes.AccAddress `json:"fromAddress"`
	ToAddress   cTypes.AccAddress `json:"toAddress"`
	PegHash     types.PegHash     `json:"pegHash"`
	Amount      int64             `json:"amount"`
}

// NewSendFiat : initializer
func NewSendFiat(fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress, pegHash types.PegHash, amount int64) SendFiat {
	return SendFiat{fromAddress, toAddress, pegHash, amount}
}

// GetSignBytes : get bytes to sign
func (in SendFiat) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		FromAddress string `json:"fromAddress"`
		ToAddress   string `json:"toAddress"`
		PegHash     string `json:"pegHash"`
		Amount      int64  `json:"amount"`
	}{
		FromAddress: in.FromAddress.String(),
		ToAddress:   in.ToAddress.String(),
		PegHash:     in.PegHash.String(),
		Amount:      in.Amount,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

func (in SendFiat) ValidateBasic() cTypes.Error {
	if len(in.FromAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.FromAddress.String())
	} else if len(in.ToAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.ToAddress.String())
	} else if len(in.PegHash) == 0 {
		return cTypes.ErrUnknownRequest("PegHash is Empty")
	} else if in.Amount <= 0 {
		return ErrNegativeAmount(DefaultCodespace, "Amount should be positive")
	}
	return nil
}

// #####SendFiat

// *****MsgBankSendFiats

// MsgBankSendFiats : high level issuance of fiats module
type MsgBankSendFiats struct {
	SendFiats []SendFiat `json:"sendFiats"`
}

// NewMsgBankSendFiats : initilizer
func NewMsgBankSendFiats(sendFiats []SendFiat) MsgBankSendFiats {
	return MsgBankSendFiats{sendFiats}
}

// ***** Implementing cTypes.Msg

var _ cTypes.Msg = MsgBankSendFiats{}

// Type : implements msg
func (msg MsgBankSendFiats) Type() string { return "bank" }

func (msg MsgBankSendFiats) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgBankSendFiats) ValidateBasic() cTypes.Error {
	if len(msg.SendFiats) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.SendFiats {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankSendFiats) GetSignBytes() []byte {
	var sendFiats []json.RawMessage
	for _, sendFiat := range msg.SendFiats {
		sendFiats = append(sendFiats, sendFiat.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		SendFiats []json.RawMessage `json:"sendFiats"`
	}{
		SendFiats: sendFiats,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgBankSendFiats) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.SendFiats))
	for i, in := range msg.SendFiats {
		addrs[i] = in.FromAddress
	}
	return addrs
}

// ##### Implement cTypes.Msg

// #####MsgBankSendFiats

// *****BuyerExecuteOrder

// BuyerExecuteOrder - transaction input
type BuyerExecuteOrder struct {
	MediatorAddress cTypes.AccAddress `json:"mediatorAddress"`
	BuyerAddress    cTypes.AccAddress `json:"buyerAddress"`
	SellerAddress   cTypes.AccAddress `json:"sellerAddress"`
	PegHash         types.PegHash     `json:"pegHash"`
	FiatProofHash   string            `json:"fiatProofHash"`
}

// NewBuyerExecuteOrder : initializer
func NewBuyerExecuteOrder(mediatorAddress cTypes.AccAddress, buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress, pegHash types.PegHash, fiatProofHash string) BuyerExecuteOrder {
	return BuyerExecuteOrder{mediatorAddress, buyerAddress, sellerAddress, pegHash, fiatProofHash}
}

// GetSignBytes : get bytes to sign
func (in BuyerExecuteOrder) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		MediatorAddress string `json:"mediatorAddress"`
		BuyerAddress    string `json:"buyerAddress"`
		SellerAddress   string `json:"sellerAddress"`
		PegHash         string `json:"pegHash"`
		FiatProofHash   string `json:"fiatProofHash"`
	}{
		MediatorAddress: in.MediatorAddress.String(),
		BuyerAddress:    in.BuyerAddress.String(),
		SellerAddress:   in.SellerAddress.String(),
		PegHash:         in.PegHash.String(),
		FiatProofHash:   in.FiatProofHash,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

func (in BuyerExecuteOrder) ValidateBasic() cTypes.Error {
	if len(in.MediatorAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.MediatorAddress.String())
	} else if len(in.SellerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.SellerAddress.String())
	} else if len(in.BuyerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.BuyerAddress.String())
	} else if len(in.PegHash) == 0 {
		return cTypes.ErrUnknownRequest("PegHash is Empty")
	} else if in.FiatProofHash == "" {
		return cTypes.ErrUnknownRequest("FiatProofHash is Empty")
	}
	return nil
}

// #####BuyerExecuteOrder

// *****MsgBankBuyerExecuteOrders

// MsgBankBuyerExecuteOrders : high level issuance of fiats module
type MsgBankBuyerExecuteOrders struct {
	BuyerExecuteOrders []BuyerExecuteOrder `json:"buyerExecuteOrders"`
}

// NewMsgBankBuyerExecuteOrders : initilizer
func NewMsgBankBuyerExecuteOrders(buyerExecuteOrders []BuyerExecuteOrder) MsgBankBuyerExecuteOrders {
	return MsgBankBuyerExecuteOrders{buyerExecuteOrders}
}

// ***** Implementing cTypes.Msg

var _ cTypes.Msg = MsgBankBuyerExecuteOrders{}

// Type : implements msg
func (msg MsgBankBuyerExecuteOrders) Type() string { return "bank" }

func (msg MsgBankBuyerExecuteOrders) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgBankBuyerExecuteOrders) ValidateBasic() cTypes.Error {
	if len(msg.BuyerExecuteOrders) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.BuyerExecuteOrders {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankBuyerExecuteOrders) GetSignBytes() []byte {
	var buyerExecuteOrders []json.RawMessage
	for _, buyerExecuteOrder := range msg.BuyerExecuteOrders {
		buyerExecuteOrders = append(buyerExecuteOrders, buyerExecuteOrder.GetSignBytes())
	}
	b, err := ModuleCdc.MarshalJSON(struct {
		BuyerExecuteOrders []json.RawMessage `json:"buyerExecuteOrders"`
	}{
		BuyerExecuteOrders: buyerExecuteOrders,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgBankBuyerExecuteOrders) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.BuyerExecuteOrders))
	for i, in := range msg.BuyerExecuteOrders {
		addrs[i] = in.MediatorAddress
	}
	return addrs
}

// ##### Implement cTypes.Msg

// #####MsgBankBuyerExecuteOrders

// *****SellerExecuteOrder

// SellerExecuteOrder - transaction input
type SellerExecuteOrder struct {
	MediatorAddress cTypes.AccAddress `json:"mediatorAddress"`
	BuyerAddress    cTypes.AccAddress `json:"buyerAddress"`
	SellerAddress   cTypes.AccAddress `json:"sellerAddress"`
	PegHash         types.PegHash     `json:"pegHash"`
	AWBProofHash    string            `json:"awbProofHash"`
}

// NewSellerExecuteOrder : initializer
func NewSellerExecuteOrder(mediatorAddress cTypes.AccAddress, buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress, pegHash types.PegHash, awbProofHash string) SellerExecuteOrder {
	return SellerExecuteOrder{mediatorAddress, buyerAddress, sellerAddress, pegHash, awbProofHash}
}

// GetSignBytes : get bytes to sign
func (in SellerExecuteOrder) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		MediatorAddress string `json:"mediatorAddress"`
		BuyerAddress    string `json:"buyerAddress"`
		SellerAddress   string `json:"sellerAddress"`
		PegHash         string `json:"pegHash"`
		AWBProofHash    string `json:"awbProofHash"`
	}{
		MediatorAddress: in.MediatorAddress.String(),
		BuyerAddress:    in.BuyerAddress.String(),
		SellerAddress:   in.SellerAddress.String(),
		PegHash:         in.PegHash.String(),
		AWBProofHash:    in.AWBProofHash,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

func (in SellerExecuteOrder) ValidateBasic() cTypes.Error {
	if len(in.MediatorAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.MediatorAddress.String())
	} else if len(in.SellerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.SellerAddress.String())
	} else if len(in.BuyerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.BuyerAddress.String())
	} else if len(in.PegHash) == 0 {
		return cTypes.ErrUnknownRequest("PegHash is Empty")
	} else if in.AWBProofHash == "" {
		return cTypes.ErrUnknownRequest("ABAProofHash is Empty")
	}
	return nil
}

// #####SellerExecuteOrder

// *****MsgBankSellerExecuteOrders

// MsgBankSellerExecuteOrders : high level issuance of fiats module
type MsgBankSellerExecuteOrders struct {
	SellerExecuteOrders []SellerExecuteOrder `json:"sellerExecuteOrders"`
}

// NewMsgBankSellerExecuteOrders : initilizer
func NewMsgBankSellerExecuteOrders(sellerExecuteOrders []SellerExecuteOrder) MsgBankSellerExecuteOrders {
	return MsgBankSellerExecuteOrders{sellerExecuteOrders}
}

// ***** Implementing cTypes.Msg

var _ cTypes.Msg = MsgBankSellerExecuteOrders{}

// Type : implements msg
func (msg MsgBankSellerExecuteOrders) Type() string { return "bank" }

func (msg MsgBankSellerExecuteOrders) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgBankSellerExecuteOrders) ValidateBasic() cTypes.Error {
	if len(msg.SellerExecuteOrders) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.SellerExecuteOrders {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankSellerExecuteOrders) GetSignBytes() []byte {
	var sellerExecuteOrders []json.RawMessage
	for _, sellerExecuteOrder := range msg.SellerExecuteOrders {
		sellerExecuteOrders = append(sellerExecuteOrders, sellerExecuteOrder.GetSignBytes())
	}
	b, err := ModuleCdc.MarshalJSON(struct {
		SellerExecuteOrders []json.RawMessage `json:"sellerExecuteOrders"`
	}{
		SellerExecuteOrders: sellerExecuteOrders,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgBankSellerExecuteOrders) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.SellerExecuteOrders))
	for i, in := range msg.SellerExecuteOrders {
		addrs[i] = in.MediatorAddress
	}
	return addrs
}

// ##### Implement cTypes.Msg

// #####MsgBankSellerExecuteOrders

// *****ReleaseAsset

// ReleaseAsset - transaction input
type ReleaseAsset struct {
	ZoneAddress  cTypes.AccAddress `json:"zoneAddress"`
	OwnerAddress cTypes.AccAddress `json:"ownerAddress"`
	PegHash      types.PegHash     `json:"pegHash"`
}

// NewReleaseAsset : initializer
func NewReleaseAsset(zoneAddress cTypes.AccAddress, ownerAddress cTypes.AccAddress, pegHash types.PegHash) ReleaseAsset {
	return ReleaseAsset{zoneAddress, ownerAddress, pegHash}
}

// GetSignBytes : get bytes to sign
func (in ReleaseAsset) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		ZoneAddress  string `json:"zoneAddress"`
		OwnerAddress string `json:"ownerAddress"`
		PegHash      string `json:"pegHash"`
	}{
		ZoneAddress:  in.ZoneAddress.String(),
		OwnerAddress: in.OwnerAddress.String(),
		PegHash:      in.PegHash.String(),
	})
	if err != nil {
		panic(err)
	}
	return bin
}

func (in ReleaseAsset) ValidateBasic() cTypes.Error {
	if len(in.OwnerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.OwnerAddress.String())
	} else if len(in.ZoneAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.ZoneAddress.String())
	} else if len(in.PegHash) == 0 {
		return cTypes.ErrUnknownRequest("PegHash is Empty")
	}
	return nil
}

// #####ReleaseAsset

// *****MsgBankReleaseAssets

// MsgBankReleaseAssets : high level release of asset module
type MsgBankReleaseAssets struct {
	ReleaseAssets []ReleaseAsset `json:"releseAssets"`
}

// NewMsgBankReleaseAssets : initilizer
func NewMsgBankReleaseAssets(releseAsset []ReleaseAsset) MsgBankReleaseAssets {
	return MsgBankReleaseAssets{releseAsset}
}

// ***** Implementing cTypes.Msg

var _ cTypes.Msg = MsgBankReleaseAssets{}

// Type : implements msg
func (msg MsgBankReleaseAssets) Type() string { return "bank" }

func (msg MsgBankReleaseAssets) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgBankReleaseAssets) ValidateBasic() cTypes.Error {
	if len(msg.ReleaseAssets) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.ReleaseAssets {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankReleaseAssets) GetSignBytes() []byte {
	var releaseAssets []json.RawMessage
	for _, releaseAsset := range msg.ReleaseAssets {
		releaseAssets = append(releaseAssets, releaseAsset.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		ReleaseAssets []json.RawMessage `json:"releaseAssets"`
	}{
		ReleaseAssets: releaseAssets,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgBankReleaseAssets) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.ReleaseAssets))
	for i, in := range msg.ReleaseAssets {
		addrs[i] = in.ZoneAddress
	}
	return addrs
}

// ##### Implement cTypes.Msg

// #####MsgBankReleaseAssets

// DefineZone : singular define zone message
// *****ACL
type DefineZone struct {
	From   cTypes.AccAddress `json:"from"`
	To     cTypes.AccAddress `json:"to"`
	ZoneID acl.ZoneID        `json:"zoneID"`
}

// NewDefineZone : new define zone struct
func NewDefineZone(from cTypes.AccAddress, to cTypes.AccAddress, zoneID acl.ZoneID) DefineZone {
	return DefineZone{from, to, zoneID}
}

// GetSignBytes : get bytes to sign
func (in DefineZone) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		From   string `json:"from"`
		To     string `json:"to"`
		ZoneID string `json:"zoneID"`
	}{
		From:   in.From.String(),
		To:     in.To.String(),
		ZoneID: in.ZoneID.String(),
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// ValidateBasic : Validate Basic
func (in DefineZone) ValidateBasic() cTypes.Error {
	if len(in.From) == 0 {
		return cTypes.ErrInvalidAddress(in.From.String())
	} else if len(in.To) == 0 {
		return cTypes.ErrInvalidAddress(in.To.String())
	} else if len(in.ZoneID) == 0 {
		return cTypes.ErrInvalidAddress(in.ZoneID.String())
	}
	return nil
}

// MsgDefineZones : message define zones
type MsgDefineZones struct {
	DefineZones []DefineZone `json:"defineZones"`
}

// NewMsgDefineZones : new message define zones
func NewMsgDefineZones(defineZones []DefineZone) MsgDefineZones {
	return MsgDefineZones{defineZones}
}

var _ cTypes.Msg = MsgDefineZones{}

// Type : implements msg
func (msg MsgDefineZones) Type() string { return "bank" }

func (msg MsgDefineZones) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgDefineZones) ValidateBasic() cTypes.Error {
	if len(msg.DefineZones) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.DefineZones {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgDefineZones) GetSignBytes() []byte {
	var defineZones []json.RawMessage
	for _, defineZone := range msg.DefineZones {
		defineZones = append(defineZones, defineZone.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		DefineZones []json.RawMessage `json:"defineZones"`
	}{
		DefineZones: defineZones,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgDefineZones) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.DefineZones))
	for i, in := range msg.DefineZones {
		addrs[i] = in.From
	}
	return addrs
}

// BuildMsgDefineZones : build define zones message
func BuildMsgDefineZones(from cTypes.AccAddress, to cTypes.AccAddress, zoneID acl.ZoneID, msgs []DefineZone) []DefineZone {
	defineZone := NewDefineZone(from, to, zoneID)
	msgs = append(msgs, defineZone)
	return msgs
}

// BuildMsgDefineZoneWithDefineZones : build define zones message
func BuildMsgDefineZoneWithDefineZones(msgs []DefineZone) cTypes.Msg {
	return NewMsgDefineZones(msgs)
}

// BuildMsgDefineZone : build define zones message
func BuildMsgDefineZone(from cTypes.AccAddress, to cTypes.AccAddress, zoneID acl.ZoneID) cTypes.Msg {
	defineZone := NewDefineZone(from, to, zoneID)
	return NewMsgDefineZones([]DefineZone{defineZone})
}

// DefineOrganization : singular define organization message
type DefineOrganization struct {
	From           cTypes.AccAddress  `json:"from"`
	To             cTypes.AccAddress  `json:"to"`
	OrganizationID acl.OrganizationID `json:"organizationID"`
	ZoneID         acl.ZoneID         `json:"zoneID"`
}

// NewDefineOrganization : new define organization struct
func NewDefineOrganization(from cTypes.AccAddress, to cTypes.AccAddress, organizationID acl.OrganizationID, zoneID acl.ZoneID) DefineOrganization {
	return DefineOrganization{from, to, organizationID, zoneID}
}

// GetSignBytes : get bytes to sign
func (in DefineOrganization) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		From           string `json:"from"`
		To             string `json:"to"`
		OrganizationID string `json:"organizationID"`
		ZoneID         string `json:"zoneID"`
	}{
		From:           in.From.String(),
		To:             in.To.String(),
		OrganizationID: in.OrganizationID.String(),
		ZoneID:         in.ZoneID.String(),
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// ValidateBasic : Validate Basic
func (in DefineOrganization) ValidateBasic() cTypes.Error {
	if len(in.From) == 0 {
		return cTypes.ErrInvalidAddress(in.From.String())
	} else if len(in.To) == 0 {
		return cTypes.ErrInvalidAddress(in.To.String())
	} else if len(in.OrganizationID) == 0 {
		return cTypes.ErrInvalidAddress(in.OrganizationID.String())
	} else if len(in.ZoneID) == 0 {
		return cTypes.ErrInvalidAddress(in.ZoneID.String())
	}
	return nil
}

// MsgDefineOrganizations : message define organizations
type MsgDefineOrganizations struct {
	DefineOrganizations []DefineOrganization `json:"defineOrganizations"`
}

// NewMsgDefineOrganizations : new message define organizations
func NewMsgDefineOrganizations(defineOrganizations []DefineOrganization) MsgDefineOrganizations {
	return MsgDefineOrganizations{defineOrganizations}
}

var _ cTypes.Msg = MsgDefineOrganizations{}

// Type : implements msg
func (msg MsgDefineOrganizations) Type() string { return "bank" }

func (msg MsgDefineOrganizations) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgDefineOrganizations) ValidateBasic() cTypes.Error {
	if len(msg.DefineOrganizations) == 0 {
		return ErrNoInputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.DefineOrganizations {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgDefineOrganizations) GetSignBytes() []byte {
	var defineOrganizations []json.RawMessage
	for _, defineOrganization := range msg.DefineOrganizations {
		defineOrganizations = append(defineOrganizations, defineOrganization.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		DefineOrganizations []json.RawMessage `json:"defineOrganizations"`
	}{
		DefineOrganizations: defineOrganizations,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgDefineOrganizations) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.DefineOrganizations))
	for i, in := range msg.DefineOrganizations {
		addrs[i] = in.From
	}
	return addrs
}

// BuildMsgDefineOrganizations : build define organization message
func BuildMsgDefineOrganizations(from cTypes.AccAddress, to cTypes.AccAddress, organizationID acl.OrganizationID, zoneID acl.ZoneID, msgs []DefineOrganization) []DefineOrganization {
	defineOrganization := NewDefineOrganization(from, to, organizationID, zoneID)
	msgs = append(msgs, defineOrganization)
	return msgs
}

// BuildMsgDefineOrganizationWithMsgs : build define organization message
func BuildMsgDefineOrganizationWithMsgs(msgs []DefineOrganization) cTypes.Msg {
	return NewMsgDefineOrganizations(msgs)
}

// BuildMsgDefineOrganization : build define organization message
func BuildMsgDefineOrganization(from cTypes.AccAddress, to cTypes.AccAddress, organizationID acl.OrganizationID, zoneID acl.ZoneID) cTypes.Msg {
	defineOrganization := NewDefineOrganization(from, to, organizationID, zoneID)
	return NewMsgDefineOrganizations([]DefineOrganization{defineOrganization})
}

// DefineACL : indular define acl message
type DefineACL struct {
	From       cTypes.AccAddress `json:"from"`
	To         cTypes.AccAddress `json:"to"`
	ACLAccount acl.ACLAccount    `json:"aclAccount"`
}

// NewDefineACL : new define acl struct
func NewDefineACL(from cTypes.AccAddress, to cTypes.AccAddress, aclAccount acl.ACLAccount) DefineACL {
	return DefineACL{from, to, aclAccount}
}

// GetSignBytes : get bytes to sign
func (in DefineACL) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		From       string         `json:"from"`
		To         string         `json:"to"`
		ACLAccount acl.ACLAccount `json:"aclAccount"`
	}{
		From:       in.From.String(),
		To:         in.To.String(),
		ACLAccount: in.ACLAccount,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// ValidateBasic : Validate Basic
func (in DefineACL) ValidateBasic() cTypes.Error {
	if len(in.From) == 0 {
		return cTypes.ErrInvalidAddress(in.From.String())
	} else if len(in.To) == 0 {
		return cTypes.ErrInvalidAddress(in.To.String())
	}
	return nil
}

// MsgDefineACLs : message define acls
type MsgDefineACLs struct {
	DefineACLs []DefineACL `json:"defineACLs"`
}

// NewMsgDefineACLs : new message define acls
func NewMsgDefineACLs(defineACLs []DefineACL) MsgDefineACLs {
	return MsgDefineACLs{defineACLs}
}

var _ cTypes.Msg = MsgDefineACLs{}

// Type : implements msg
func (msg MsgDefineACLs) Type() string { return "bank" }

func (msg MsgDefineACLs) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgDefineACLs) ValidateBasic() cTypes.Error {
	if len(msg.DefineACLs) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.DefineACLs {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgDefineACLs) GetSignBytes() []byte {
	var defineACLs []json.RawMessage
	for _, defineACL := range msg.DefineACLs {
		defineACLs = append(defineACLs, defineACL.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		DefineACLs []json.RawMessage `json:"defineACLs"`
	}{
		DefineACLs: defineACLs,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgDefineACLs) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.DefineACLs))
	for i, in := range msg.DefineACLs {
		addrs[i] = in.From
	}
	return addrs
}

// BuildMsgDefineACLs : build define acls message
func BuildMsgDefineACLs(from cTypes.AccAddress, to cTypes.AccAddress, aclAccount acl.ACLAccount, msgs []DefineACL) []DefineACL {
	defineACL := NewDefineACL(from, to, aclAccount)
	msgs = append(msgs, defineACL)
	return msgs
}

// BuildMsgDefineACLWithACLs : build define acls message
func BuildMsgDefineACLWithACLs(msgs []DefineACL) cTypes.Msg {
	return NewMsgDefineACLs(msgs)
}

// BuildMsgDefineACL : build define acls message
func BuildMsgDefineACL(from cTypes.AccAddress, to cTypes.AccAddress, aclAccount acl.ACLAccount) cTypes.Msg {
	defineACL := NewDefineACL(from, to, aclAccount)
	return NewMsgDefineACLs([]DefineACL{defineACL})
}

// #####ACL

// #####Comdex
