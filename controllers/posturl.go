package controllers

import (
	"database/sql"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
	"github.com/lalathealter/artkeeper/psql"
)

var PostURLhandler = factorAPIHandler(
	readPostURLRequest,
	savePostURL,
	respondPostURL,
)

func readPostURLRequest(r *http.Request) (models.Message, error) {
	return parseJSONMessage(r, models.PostURLRequest{})
}

func savePostURL(db *sql.DB) dbcaller {
	return func(msg models.Message) (dbresult, error) {
		pr := msg.(models.PostURLRequest)
		sqlstatement := psql.InsertOneURL
		return db.Exec(sqlstatement, extractFieldValues(&pr)...)
	}
}

func respondPostURL(w http.ResponseWriter, _ dbresult) {
	w.WriteHeader(http.StatusNoContent)
}
