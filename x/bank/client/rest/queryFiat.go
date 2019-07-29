package rest

import (
	"fmt"
	"net/http"

	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/fiatFactory"
	"github.com/gorilla/mux"
)

// QueryFiatHandlerFn :
func QueryFiatHandlerFn(cliContext context.CLIContext, r *mux.Router, cdc *wire.Codec, storeName string, decoder sdk.FiatPegDecoder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		vars := mux.Vars(r)
		peghashstr := vars["peghash"]
		nodeURI := "tcp://0.0.0.0:56657"

		cliCtx := cliContext
		cliCtx = cliCtx.WithNodeURI(nodeURI)
		cliCtx = cliCtx.WithTrustNode(true)

		peghashHex, err := sdk.GetFiatPegHashHex(peghashstr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.QueryStore(fiatFactory.FiatPegHashStoreKey(peghashHex), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't query account. Error: %s", err.Error()))
			return
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		account, err := decoder(res)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't parse query result. Result: %s, Error: %s", res, err.Error()))
			return
		}

		output, err := cdc.MarshalJSON(account)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't marshall query result. Error: %s", err.Error()))
			return
		}

		w.Write(output)
	}

}
