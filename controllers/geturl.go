package controllers

import (
	"net/http"

	"github.com/lalathealter/artkeeper/models"
)

var GetURLHandler = factorAPIHandler(
	readGetURLRequest,
	respondGetURL,
)

func readGetURLRequest(r *http.Request) (models.Message, error) {
	return parseURLValues(r, models.GetURLRequest{})
}
var GetLatestURLsHandler = factorAPIHandler(
	readGetLatestURLsRequest,
	respondGetURL,
)

func readGetLatestURLsRequest(r *http.Request) (models.Message, error) {
	return parseURLValues(r, models.GetLatestURLsRequest{})
}
func respondGetURL(w http.ResponseWriter, dbr models.DBResult) {
	encodeJSONResponses(w, dbr, models.GetURLResponse{})
}
