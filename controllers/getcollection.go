package controllers

import (
	"database/sql"
	"log"
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
	rows := dbr.(*sql.Rows)
	responsesArr, err := parseSQLRows(models.CollectionResponse{}, rows)

	if err != nil {
		log.Panicln(err)
	}
	
	if len(responsesArr) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	singleColl := responsesArr[0]
	sendEncodedJSON(w, singleColl)
}

var GetURLsFromCollectionHandler = factorAPIHandler(
	readGetURLsFromCollectionRequest,
	respondGetURLsFromCollection,
)

func readGetURLsFromCollectionRequest(r *http.Request) (models.Message, error) {
	return parseURLValues(r, models.GetURLsFromCollectionRequest{})
}

func respondGetURLsFromCollection(w http.ResponseWriter, dbr models.DBResult) {
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

