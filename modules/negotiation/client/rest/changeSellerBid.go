package rest

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	rest2 "github.com/persistenceOne/comdexCrust/client/rest"
	"github.com/persistenceOne/comdexCrust/kafka"
	"github.com/persistenceOne/comdexCrust/modules/acl"
	negotiationTypes "github.com/persistenceOne/comdexCrust/modules/negotiation/internal/types"
	"github.com/persistenceOne/comdexCrust/types"
)

type changeSellerBidBody struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	To       string       `json:"to" valid:"required~Enter the ToAddress,matches(^persist[a-z0-9]{39}$)~ToAddress is Invalid"`
	Bid      int64        `json:"bid" valid:"required~Enter the Bid,matches(^[1-9]{1}[0-9]*$)~Enter valid Bid"`
	Time     int64        `json:"time" valid:"required~Enter the Time,matches(^[1-9]{1}[0-9]*$)~Enter valid Time"`
	PegHash  string       `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	Password string       `json:"password" valid:"required~Enter the Password"`
	Mode     string       `json:"mode"`
}

func ChangeSellerBidRequestHandlerFn(cliCtx context.CLIContext, kafkaBool bool, kafkaState kafka.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req changeSellerBidBody

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			rest2.WriteErrorResponse(w, cTypes.NewError(negotiationTypes.DefaultCodeSpace, http.StatusBadRequest, err.Error()))
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, name, err := context.GetFromFields(req.BaseReq.From, false)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrFromName(negotiationTypes.DefaultCodeSpace))
			return
		}

		cliCtx = cliCtx.WithFromAddress(fromAddr)
		cliCtx = cliCtx.WithFromName(name)

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryACLAccount", fromAddr), nil)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(negotiationTypes.DefaultCodeSpace, "ACL Account"))
			return
		}

		if len(res) == 0 {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(negotiationTypes.DefaultCodeSpace))
			return
		}

		var account acl.ACLAccount
		cliCtx.Codec.MustUnmarshalJSON(res, &account)
		if !account.GetACL().ChangeSellerBid {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(negotiationTypes.DefaultCodeSpace))
			return
		}

		to, err := cTypes.AccAddressFromBech32(req.To)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(negotiationTypes.DefaultCodeSpace, req.To))
			return
		}

		pegHashHex, err := types.GetAssetPegHashHex(req.PegHash)
		negotiationID := types.NegotiationID(append(append(to.Bytes(), fromAddr.Bytes()...), pegHashHex...))

		proposedNegotiation := &types.BaseNegotiation{
			NegotiationID: negotiationID,
			BuyerAddress:  to,
			SellerAddress: fromAddr,
			PegHash:       pegHashHex,
			Bid:           req.Bid,
			Time:          req.Time,
		}

		msg := negotiationTypes.BuildMsgChangeSellerBid(proposedNegotiation)

		if kafkaBool == true {
			ticketID := kafka.TicketIDGenerator("CHSB")
			jsonResponse := kafka.SendToKafka(kafka.NewKafkaMsgFromRest(msg, ticketID, req.BaseReq, cliCtx, req.Mode, req.Password), kafkaState, cliCtx.Codec)
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write(jsonResponse)
		} else {
			output, err := rest2.SignAndBroadcast(req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
			if err != nil {
				rest2.WriteErrorResponse(w, err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(output)
		}
	}
}
