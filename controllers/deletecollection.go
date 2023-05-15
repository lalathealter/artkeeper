package controllers

import (
	"database/sql"
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


var DetachTagFromCollectionHandler = factorAPIHandler(
	readDetachTagFromCollection,
	respondDetachTagFromCollection,
)

func readDetachTagFromCollection(r *http.Request) (models.Message, error ) {
	return parseURLValues(r, models.DetachTagFromCollectionRequest{})
}

func respondDetachTagFromCollection(w http.ResponseWriter, dbr models.DBResult) {
	execRes := dbr.(sql.Result)
	if affRows, _ := execRes.RowsAffected(); affRows < 1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Destined collection doesn't exist"))
		return
	}


	w.WriteHeader(http.StatusNoContent)
}
