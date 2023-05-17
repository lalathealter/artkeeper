package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/lalathealter/artkeeper/models"
	"golang.org/x/crypto/bcrypt"
)

var UserRegistrationHandler = factorAPIHandler(
	readUserRegistrationRequest,
	respondUserRegistrationRequest,
)

func readUserRegistrationRequest(r *http.Request) (models.Message, error) {
	msg, err := parseJSONMessage(r, models.RegisterUserRequest{})
	msg.Password.Update(bcryptString)
	return msg, err 
}

func bcryptString(input string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
	}
	return string(hash)
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
