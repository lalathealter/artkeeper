package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
	"github.com/lalathealter/artkeeper/psql"
)

var GetCollectionHandler = factorAPIHandler(
	readGetCollectionRequest,
	lookupCollection,
	respondGetCollection,
)

func readGetCollectionRequest(r *http.Request) (models.Message, error) {
	return parseURLParams(r, models.GetCollectionRequest{})
}

func lookupCollection(db *sql.DB) dbcaller {
	return func(msg models.Message) (dbresult, error) {

		gcr := msg.(models.GetCollectionRequest)

		sqlstatement := psql.SelectOneCollection
		sqlargs := []any{ gcr.ID }

		return db.Query(sqlstatement, sqlargs...)
	}
}

func respondGetCollection(w http.ResponseWriter, dbr dbresult) {
	rows := dbr.(*sql.Rows)
	responses, err := parseSQLRows(models.GetCollectionResponse{}, rows)

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


