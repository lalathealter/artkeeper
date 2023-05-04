package controllers

import (
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


func respondPutInCollection(w http.ResponseWriter, _ models.DBResult) {
	w.WriteHeader(http.StatusNoContent)
}
