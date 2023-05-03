package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
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
	rows := dbr.(*sql.Rows)

	responses, err := parseSQLRows(models.GetURLResponse{}, rows)
	if err != nil {
		log.Panicln(err)
	}

	if len(responses) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(responses); err != nil {
		log.Panicln(err)
	}
}
