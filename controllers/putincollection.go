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
	return parseJSONMessage(r, models.PutInCollectionRequest{})
}


func respondPutInCollection(w http.ResponseWriter, _ models.DBResult) {
	w.WriteHeader(http.StatusNoContent)
}
