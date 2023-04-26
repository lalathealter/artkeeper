package controllers

import (
	"net/http"

	"github.com/lalathealter/artkeeper/models"
)

var PostCollectionHandler = factorAPIHandler(
	readPostCollectionRequest,
	respondPostCollection,
)

func readPostCollectionRequest(r *http.Request) (models.Message, error) {
	return parseJSONMessage(r, models.PostCollectionRequest{})
}

func respondPostCollection(w http.ResponseWriter, _ models.DBResult) {
	w.WriteHeader(http.StatusNoContent)
}
