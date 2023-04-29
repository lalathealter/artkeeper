package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"

	"github.com/lalathealter/artkeeper/models"
	"github.com/lalathealter/artkeeper/psql"
	_ "github.com/lib/pq"
)

type readmethod func(*http.Request) (models.Message, error)
// type dbcaller func(models.Message) (dbresult, error)
// type calldbmethod func(*sql.DB) dbcaller
type respondmethod func(http.ResponseWriter, models.DBResult)

func factorAPIHandler(
	read readmethod,
	respond respondmethod,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				fmt.Println(rec)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(rec.(string)))
			}
		}()

		msg, err := read(r)
		if err != nil {
			log.Panicln(err)
		}
		err = msg.VerifyValues()
		if err != nil {
			log.Panicln(err)
		}

		fmt.Println("Gone through", msg)

		db := psql.GetDB()
		dbres, err := msg.Call(db)
		if err != nil {
			log.Panicln(err)
		}

		respond(w, dbres)
	}
}

func parseJSONMessage[T models.Message](r *http.Request, target T) (T, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&target)
	if err != nil {
		return *(new(T)), (err)
	}

	return target, nil
}

func parseURLParams[T models.Message](r *http.Request, target T) (T, error) {
	urlvals, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return *(new(T)), err
	}
	iterm := reflect.ValueOf(&target).Elem()
	for i := 0; i < iterm.NumField(); i++ {
		key := iterm.Type().Field(i).Tag.Get("urlparam")
		paramval := urlvals.Get(key) // may be empty string

		typedfield := iterm.Field(i).Interface()
		reffedval, err := models.ReflectCastedStringlike(paramval, typedfield)
		if err != nil {
			return *(new(T)), err
		}
		iterm.Field(i).Set(reffedval)
	}
	return target, nil
}

// func getURLParam(uvals *url.Values, key string) (string, error) {
// 	paramfound := uvals.Has(key)
// 	if !paramfound {
// 		return "", fmt.Errorf("URL Parameter %v wasn't provided;", key)
// 	}
// 	param := uvals.Get(key)
// 	return param, nil
// }

func parseSQLRows[T any](responseFormat T, rows *sql.Rows) ([]*T, error) {
	defer rows.Close()

	results := make([]*T, 0)
	i := 0
	for rows.Next() {
		results = append(results, new(T))
		fieldMap, err := ExtractFieldPointersIntoNamedMap(results[i])
		if err != nil {
			return nil, err
		}
		sqlColumns, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		orderedPointersArr := make([]any, len(fieldMap))
		for i, column := range sqlColumns {
			orderedPointersArr[i] = fieldMap[column]
		}
		err = rows.Scan(orderedPointersArr...)
		if err != nil {
			return nil, err
		}
		i++
	}

	return results, rows.Err()
}

func ExtractFieldPointersIntoNamedMap[T any](in *T) (map[string]any, error) {
	fieldMap := make(map[string]any)
	iter := reflect.ValueOf(in).Elem()
	for i := 0; i < iter.NumField(); i++ {
		currPtr	 := iter.Field(i).Addr().Interface()
		columnName := iter.Type().Field(i).Tag.Get("field") // sql field tag
		if columnName == "" {
			return nil, fmt.Errorf("Struct type %T doesn't provide the necessary field tags for successful sql parsing", *in)
		}
		fieldMap[columnName] = currPtr
		
	}
	return fieldMap, nil
}

// func ExtractFieldValues[T any](in *T) []any {
// 	iter := reflect.ValueOf(in).Elem()
// 	fieldvals := make([]any, iter.NumField())
// 	for i := 0; i < iter.NumField(); i++ {
// 		fieldvals[i] = iter.Field(i).Interface()
// 	}
// 	return fieldvals
// }
