package kafka

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	dbm "github.com/tendermint/tendermint/libs/db"

	"github.com/commitHub/commitBlockchain/codec"
)

// SetTicketIDtoDB : initiates ticketid in Database
func SetTicketIDtoDB(ticketID Ticket, kafkaDB *dbm.GoLevelDB, cdc *codec.Codec, msg []byte) {

	ticketid, err := cdc.MarshalJSON(ticketID)
	if err != nil {
		panic(err)
	}

	kafkaDB.Set(ticketid, msg)
	return
}

// AddResponseToDB : Updates response to DB
func AddResponseToDB(ticketID Ticket, response []byte, kafkaDB *dbm.GoLevelDB, cdc *codec.Codec) {
	ticketid, err := cdc.MarshalJSON(ticketID)
	if err != nil {
		panic(err)
	}
	kafkaDB.SetSync(ticketid, response)
	return
}

// GetResponseFromDB : gives the response from DB
func GetResponseFromDB(ticketID Ticket, kafkaDB *dbm.GoLevelDB, cdc *codec.Codec) []byte {
	ticketid, err := cdc.MarshalJSON(ticketID)
	if err != nil {
		panic(err)
	}

	return kafkaDB.Get(ticketid)
}

// QueryDB : REST outputs info from DB
func QueryDB(cdc *codec.Codec, r *mux.Router, kafkaDB *dbm.GoLevelDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		iDByte, err := cdc.MarshalJSON(vars["ticketid"])
		if err != nil {
			panic(err)
		}
		var response []byte
		if kafkaDB.Has(iDByte) == true {
			response = GetResponseFromDB(Ticket(vars["ticketid"]), kafkaDB, cdc)
		} else {
			output, err := cdc.MarshalJSON("The ticket ID does not exist, it must have been deleted, Query the chain to know")
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(fmt.Sprintf("ticket ID does not exist. Error: %s", err.Error())))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write(response)
		return
	}
}
