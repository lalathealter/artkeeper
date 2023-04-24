package controllers

import (
	"database/sql"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
	"github.com/lalathealter/artkeeper/psql"
)
var PutInCollectionHandler = factorAPIHandler(
	readPutInCollectionRequest,
	savePutInCollection,
	respondPutInCollection,
)

func readPutInCollectionRequest(r *http.Request) (models.Message, error) {
	return parseJSONMessage(r, models.PutInCollectionRequest{})
}

func savePutInCollection(db *sql.DB) dbcaller {
	return func(msg models.Message) (dbresult, error) {
		putcr := msg.(models.PutInCollectionRequest)

		sqlstatement := psql.UpdateLinksInCollection
		sqlargs := []any{
			putcr.LinkID,
			putcr.CollID,
		}
		return db.Query(sqlstatement, sqlargs...)
	}
}

func respondPutInCollection(w http.ResponseWriter, _ dbresult) {
	w.WriteHeader(http.StatusNoContent)
}
