package rest

import (
	"fmt"
	"net/http"

	"github.com/commitHub/commitBlockchain/client/context"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/slashing"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec) {
	r.HandleFunc(
		"/slashing/signing_info/{validator}",
		signingInfoHandlerFn(cliCtx, "slashing", cdc),
	).Methods("GET")
}

// http request handler to query signing info
func signingInfoHandlerFn(cliCtx context.CLIContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		pk, err := sdk.GetConsPubKeyBech32(vars["validator"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		key := slashing.GetValidatorSigningInfoKey(sdk.ConsAddress(pk.Address()))

		res, err := cliCtx.QueryStore(key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query signing info. Error: %s", err.Error())))
			return
		}

		var signingInfo slashing.ValidatorSigningInfo

		err = cdc.UnmarshalBinary(res, &signingInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't decode signing info. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(signingInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
