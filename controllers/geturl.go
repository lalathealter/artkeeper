package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
	"github.com/lalathealter/artkeeper/psql"
)

var GetURLHandler = factorAPIHandler(
	readGetURLRequest,
	switchLookupURL,
	respondGetURL,
)

func readGetURLRequest(r *http.Request) (models.Message, error) {
	return parseURLParams(r, models.GetURLRequest{})
}

func switchLookupURL(db *sql.DB) dbcaller {
	return func(m models.Message) (dbresult, error) {
		var sqlstatement string
		var sqlargs []any

		switch greq := m.(models.GetURLRequest); {
		case *greq.ID != "":
			sqlstatement = psql.SelectOneURL
			sqlargs = []any{greq.ID}
		// case *greq.Collection != "":

		default:
			sqlstatement = psql.SelectAllURLs

		}

		return db.Query(sqlstatement, sqlargs...)
	}
}

func respondGetURL(w http.ResponseWriter, dbr dbresult) {
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
