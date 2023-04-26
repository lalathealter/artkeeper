package controllers

import (
	"net/http"

	"github.com/lalathealter/artkeeper/models"
)

var DeleteURLHandler = factorAPIHandler(
	readDeleteURLRequest,
	respondDeleteURL,
)

func readDeleteURLRequest(r *http.Request) (models.Message, error) {
	return parseJSONMessage(r, models.DeleteURLRequest{})
}

func respondDeleteURL(w http.ResponseWriter, _ models.DBResult) {
	w.WriteHeader(http.StatusNoContent)
}
