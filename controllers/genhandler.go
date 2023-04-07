package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
	_ "github.com/lib/pq"
)

type dbresult interface{}
type readmethod func(*http.Request) (models.Message, error)
type dbcaller func(models.Message) (dbresult, error)
type calldbmethod func(*sql.DB) dbcaller
type respondmethod func(http.ResponseWriter, dbresult)
type httpreciever func(http.ResponseWriter, *http.Request)

func factorapihandler(
	read readmethod,
	call calldbmethod,
	respond respondmethod,
) httpreciever {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				fmt.Println(rec)
				w.WriteHeader(http.StatusBadRequest)
			}
		}()

		msg, err := read(r)
		if err != nil {
			log.Panicln(err)
		}

		fmt.Println("Gone through", msg)

		db := dbconnect()
		defer db.Close()

		dbres, err := call(db)(msg)
		if err != nil {
			log.Panicln(err)
		}

		respond(w, dbres)
	}
}

func parsejsonmessage[T models.Message](r *http.Request, target T) (T, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&target)
	if err != nil {
		return *(new(T)), (err)
	}
	err = target.Verifyvalues()
	if err != nil {
		return *(new(T)), err
	}

	return target, nil
}
