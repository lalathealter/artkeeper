package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
)

var GetCollectionHandler = factorAPIHandler(
	readGetCollectionRequest,
	respondGetCollection,
)

func readGetCollectionRequest(r *http.Request) (models.Message, error) {
	return parseURLParams(r, models.GetCollectionRequest{})
}

func respondGetCollection(w http.ResponseWriter, dbr models.DBResult) {
	rows := dbr.(*sql.Rows)
	responses, err := parseSQLRows(models.GetCollectionResponse{}, rows)

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


