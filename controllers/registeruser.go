package controllers

import (
	"database/sql"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
)

var UserRegistrationHandler = factorAPIHandler(
	readUserRegistrationRequest,
	respondUserRegistrationRequest,
)

func readUserRegistrationRequest(r *http.Request) (models.Message, error) {
	msg, err := parseJSONMessage(r, models.RegistrateUserRequest{})

	return msg, err 
}

func respondUserRegistrationRequest(w http.ResponseWriter, dbr models.DBResult) {
	execRes := dbr.(sql.Result)
	if rowsAff, _ := execRes.RowsAffected(); rowsAff < 1 {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Submitted username is already taken"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
