package controllers

import (
	"net/http"

	"github.com/lalathealter/artkeeper/models"
)

var PostURLhandler = factorAPIHandler(
	readPostURLRequest,
	respondPostURL,
)

func readPostURLRequest(r *http.Request) (models.Message, error) {
	return parseJSONMessage(r, models.PostURLRequest{})
}


func respondPostURL(w http.ResponseWriter, _ models.DBResult) {
	w.WriteHeader(http.StatusNoContent)
}
