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
	lookupURL,
	respondGetURL,
)

func readGetURLRequest(r *http.Request) (models.Message, error) {
	return parseURLParams(r, models.GetURLRequest{})
}

func lookupURL(db *sql.DB) dbcaller {
	return func(msg models.Message) (dbresult, error) {
		greq := msg.(models.GetURLRequest) 
		sqlstatement := psql.SelectOneURL
		sqlargs := []any{greq.ID}
		return db.Query(sqlstatement, sqlargs...)
	}
}

var GetLatestURLsHandler = factorAPIHandler(
	readGetLatestURLsRequest,
	getLatestURLs,
	respondGetURL,
)

func readGetLatestURLsRequest(r *http.Request) (models.Message, error) {
	return parseURLParams(r, models.GetLatestURLsRequest{})
}

func getLatestURLs(db *sql.DB) dbcaller {
	return func(msg models.Message) (dbresult, error) {
		greqLatest := msg.(models.GetLatestURLsRequest)
		sqlstatement := psql.SelectLatestURLsWithPagination
		var sqlargs []any 
		if (*greqLatest.Limit == "0") {
			sqlargs = []any{psql.DefaultPaginationLimit, greqLatest.Offset}
		} else {
			sqlargs = []any{greqLatest.Limit, greqLatest.Offset}
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
