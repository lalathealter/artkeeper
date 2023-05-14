package controllers

import (
	"database/sql"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
)

var PutInCollectionHandler = factorAPIHandler(
	readPutInCollectionRequest,
	respondPutInCollection,
)

func readPutInCollectionRequest(r *http.Request) (models.Message, error) {
	pcr := models.PutInCollectionRequest{}
	pcr, err := parseJSONMessage(r, pcr)
	if err != nil {
		return nil, err
	}
	return parseURLValues(r, pcr)
}


func respondPutInCollection(w http.ResponseWriter, dbr models.DBResult) {
	if dbr == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Submitted link doesn't exist"))
		return 
	}

	execRes := dbr.(sql.Result)
	if affRows, _ := execRes.RowsAffected(); affRows < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Destined collection doesn't exist"))
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}
