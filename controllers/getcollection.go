package controllers

import (
	"net/http"

	"github.com/lalathealter/artkeeper/models"
)

var GetCollectionHandler = factorAPIHandler(
	readGetCollectionRequest,
	respondGetCollection,
)

func readGetCollectionRequest(r *http.Request) (models.Message, error) {
	return parseURLValues(r, models.GetCollectionRequest{})
}

func respondGetCollection(w http.ResponseWriter, dbr models.DBResult) {
	encodeJSONResponses(w, dbr, models.GetCollectionResponse{})
}

var GetURLsFromCollectionHandler = factorAPIHandler(
	readGetURLsFromCollectionRequest,
	respondGetURLsFromCollection,
)

func readGetURLsFromCollectionRequest(r *http.Request) (models.Message, error) {
	return parseURLValues(r, models.GetURLsFromCollectionRequest{})
}

func respondGetURLsFromCollection(w http.ResponseWriter, dbr models.DBResult) {
	encodeJSONResponses(w, dbr, models.GetURLResponse{})
}

