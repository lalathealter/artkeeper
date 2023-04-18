package controllers

import (
	"database/sql"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
	"github.com/lalathealter/artkeeper/psql"
	"github.com/lib/pq"
)

var PostCollectionHandler = factorAPIHandler(
	readPostCollectionRequest,
	savePostCollection,
	respondPostCollection,
)

func readPostCollectionRequest(r *http.Request) (models.Message, error) {
	return parseJSONMessage(r, models.PostCollectionRequest{})
}

func savePostCollection(db *sql.DB) dbcaller {
	return func(msg models.Message) (dbresult, error) {
		pcr := msg.(models.PostCollectionRequest)

		sqlstatement := psql.InsertOneCollection
		sqlargs := []any{
			pq.Array(pcr.LinkIDs),
			pcr.UserID,
			pcr.Description,
		}
		return db.Query(sqlstatement, sqlargs...)
	}
}

func respondPostCollection(w http.ResponseWriter, _ dbresult) {
	w.WriteHeader(http.StatusNoContent)
}
