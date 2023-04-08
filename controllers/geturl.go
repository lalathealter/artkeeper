package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
)

var GetOneURLHandler = factorAPIHandler(
	readGetURL,
	lookupOneURL,
	giveOneURL,
)

func readGetURL(r *http.Request) (models.Message, error) {
	gr := models.GetURLRequest{}
	err := parseURLParams(r, &gr)
	if err != nil {
		return nil, err
	}

	return gr, nil
}

func lookupOneURL(db *sql.DB) dbcaller {
	return func(m models.Message) (dbresult, error) {
		greq := m.(models.GetURLRequest)
		sqlstatement := dbselecturl
		v := db.QueryRow(sqlstatement, greq.ID)
		return v, nil
	}
}

func giveOneURL(w http.ResponseWriter, dbr dbresult) {
	w.WriteHeader(http.StatusNoContent)
	row := dbr.(*sql.Row)

	gres := models.GetURLResponse{}
	err := row.Scan(extractFieldPointers(&gres)...)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(gres)
}
