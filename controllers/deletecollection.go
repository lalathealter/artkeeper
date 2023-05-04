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
	return parseURLValues(r, models.DeleteCollectionRequest{})
}

func respondDeleteCollection(w http.ResponseWriter, _ models.DBResult) {
	w.WriteHeader(http.StatusNoContent)
}

var DeleteURLFromCollection = factorAPIHandler(
	readDeleteURLFromCollectionRequest,
	respondDeleteURLFromCollection,
)

func readDeleteURLFromCollectionRequest(r *http.Request) (models.Message, error) {
	return parseURLValues(r, models.DeleteURLFromCollectionRequest{})
}

func respondDeleteURLFromCollection(w http.ResponseWriter, _ models.DBResult) {
	w.WriteHeader(http.StatusNoContent)
}
