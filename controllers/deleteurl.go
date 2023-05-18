package controllers

import (
	"database/sql"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
)

var DeleteURLHandler = factorAPIHandler(
	readDeleteURLRequest,
	respondDeleteURL,
)

func readDeleteURLRequest(r *http.Request) (models.Message, error) {
	return parseURLValues(r, models.DeleteURLRequest{})
}

func respondDeleteURL(w http.ResponseWriter, dbr models.DBResult) {
	execRes := dbr.(sql.Result)
	if affRows, _ := execRes.RowsAffected(); affRows < 1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Submitted link doesn't exist"))
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}
