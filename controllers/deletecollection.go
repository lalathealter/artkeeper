package controllers

import (
	"net/http"

	"github.com/lalathealter/artkeeper/models"
)

var DeleteCollectionHandler = factorAPIHandler(
	readDeleteCollectionRequest,
	respondDeleteCollection,
)

func readDeleteCollectionRequest(r *http.Request) (models.Message, error) {
	return parseJSONMessage(r, models.DeleteCollectionRequest{})
}

func respondDeleteCollection(w http.ResponseWriter, _ models.DBResult) {
	w.WriteHeader(http.StatusNoContent)
}
