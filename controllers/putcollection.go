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
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Submitted link doesn't exist"))
		return 
	}

	execRes := dbr.(sql.Result)
	if affRows, _ := execRes.RowsAffected(); affRows < 1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Destined collection doesn't exist"))
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}


var AttachTagToCollectionHandler = factorAPIHandler(
	readAttachTagToCollectionRequest,
	respondAttachTagToCollection,
)

func readAttachTagToCollectionRequest(r *http.Request) (models.Message, error) {
	return parseURLValues(r, models.AttachTagToCollectionRequest{})
}


func respondAttachTagToCollection(w http.ResponseWriter, dbr models.DBResult) {
	execRes := dbr.(sql.Result)
	if affRows, _ := execRes.RowsAffected(); affRows < 1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Destined collection doesn't exist"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
