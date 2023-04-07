package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
)

var Posturlhandler = factorapihandler(
	readposturl,
	saveposturl,
	posturlrespond,
)

func readposturl(r *http.Request) (models.Message, error) {
	return parsejsonmessage(r, models.Posturl{})
}

func saveposturl(db *sql.DB) dbcaller {
	return func(m models.Message) (dbresult, error) {
		p := m.(models.Posturl)
		sqlstatement := dbinserturl
		_, err := db.Exec(sqlstatement, p.Link, p.Description, 1)
		return nil, err
	}
}

func posturlrespond(w http.ResponseWriter, _ dbresult) {
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "Your resource has been accepted")
}
