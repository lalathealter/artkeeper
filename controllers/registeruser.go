package controllers

import (
	"database/sql"
	"net/http"

	"github.com/lalathealter/artkeeper/controllers/auth"
	"github.com/lalathealter/artkeeper/models"
)

var UserRegistrationHandler = factorAPIHandler(
	readUserRegistrationRequest,
	respondUserRegistration,
)

func readUserRegistrationRequest(r *http.Request) (models.Message, error) {
	msg, err := parseJSONMessage(r, models.RegisterUserRequest{})
	if err != nil {
		return nil, err
	}
	passhash, err := auth.DecryptPassword(r, msg.Password, msg.ClientNonce)
	if err != nil {
		return nil, err
	}

	msg.Password.ReplaceWith(passhash)
	msg.Password.Update(auth.BcryptString)
	return msg, err 
}



func respondUserRegistration(w http.ResponseWriter, dbr models.DBResult) {
	execRes := dbr.(sql.Result)
	if rowsAff, _ := execRes.RowsAffected(); rowsAff < 1 {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Submitted username is already taken"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

var UpdateUserHandler = factorAPIHandler(
	readUpdateUserRequest,
	respondUpdateUser,
)

func readUpdateUserRequest(r *http.Request) (models.Message, error) {
	msg, err := parseJSONMessage(r, models.UpdateUserRequest{})
	if err != nil {
		return nil, err
	}
	
	jwtStr, err := auth.ParseJWTCookie(r)
	if err != nil {
		return nil, err
	}
	jwt, err := auth.ParseLoadFromJWT(jwtStr)
	if err != nil {
		return nil, err
	}
	passedUsername := models.Username(jwt.Name)
	msg.OldUsername = &passedUsername 
	if msg.Password != nil {
		msg.Password.Update(auth.BcryptString)
	}
	return msg, err
}



func respondUpdateUser(w http.ResponseWriter, dbr models.DBResult) {
	execRes, wasNotEmpty := dbr.(sql.Result)
	if !wasNotEmpty {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Profile update failed: recieved nothing to change;"))
		return
	}

	if rowsAff, _ := execRes.RowsAffected(); rowsAff < 1 {
		w.WriteHeader(http.StatusNotFound)
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}
