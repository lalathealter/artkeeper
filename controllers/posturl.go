package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
	"github.com/lalathealter/artkeeper/psql"
)

var PostURLhandler = factorAPIHandler(
	readPostURL,
	savePostURL,
	respondPostURL,
)

func readPostURL(r *http.Request) (models.Message, error) {
	return parseJSONMessage(r, models.PostURLRequest{})
}

func savePostURL(db *sql.DB) dbcaller {
	return func(m models.Message) (dbresult, error) {
		p := m.(models.PostURLRequest)
		sqlstatement := psql.InsertURL
		_, err := db.Exec(sqlstatement, p.Link, p.Description, 1)
		return nil, err
	}
}

func respondPostURL(w http.ResponseWriter, _ dbresult) {
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "Your resource has been accepted")
}
