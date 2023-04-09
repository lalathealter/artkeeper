package controllers

import (
	"database/sql"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
	"github.com/lalathealter/artkeeper/psql"
)

var DeleteURLHandler = factorAPIHandler(
	readDeleteURLRequest,
	deleteURL,
	respondDeleteURL,
)

func readDeleteURLRequest(r *http.Request) (models.Message, error) {
	return parseJSONMessage(r, models.DeleteURLRequest{})
}

func deleteURL(db *sql.DB) dbcaller {
	return func(msg models.Message) (dbresult, error) {
		dr := msg.(models.DeleteURLRequest)
		sqlstatement := psql.DeleteOneURL
		return db.Exec(sqlstatement, extractFieldValues(&dr)...)
	}
}

func respondDeleteURL(w http.ResponseWriter, _ dbresult) {
	w.WriteHeader(http.StatusNoContent)
}
