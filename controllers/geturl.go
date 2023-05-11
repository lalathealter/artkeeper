package controllers

import (
	"database/sql"
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
	respondGetLatestURLs,
)

func readGetLatestURLsRequest(r *http.Request) (models.Message, error) {
	return parseURLValues(r, models.GetLatestURLsRequest{})
}

func respondGetLatestURLs(w http.ResponseWriter, dbr models.DBResult) {
	rows := dbr.(*sql.Rows)
	responsesArr, err := parseSQLRows(models.URLResponse{}, rows)
	if err != nil {
		log.Panicln(err)
	}

	sendEncodedJSON(w, responsesArr)
}


func respondGetURL(w http.ResponseWriter, dbr models.DBResult) {
	rows := dbr.(*sql.Rows)
	responsesArr, err := parseSQLRows(models.URLResponse{}, rows)
	if err != nil {
		log.Panicln(err)
	}
	
	if len(responsesArr) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sendEncodedJSON(w, responsesArr)
}
