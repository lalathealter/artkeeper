package controllers

import (
	"net/http"

	"github.com/lalathealter/artkeeper/controllers/auth"
	"github.com/lalathealter/artkeeper/models"
	"golang.org/x/crypto/bcrypt"
)

var PostSessionHandler = factorAPIHandler(
	readPostSessionRequest,
	respondPostSessionRequest,
)

func readPostSessionRequest(r *http.Request) (models.Message, error) {
	msg, err := parseJSONMessage(r, models.PostSessionRequest{})
	if err != nil {
		return nil, err
	}
	passhash, err := auth.DecryptPassword(r, msg.Password, msg.ClientNonce)
	if err != nil {
		return nil, err
	}
	msg.Password.ReplaceWith(passhash)
	return msg, err
}

func respondPostSessionRequest(w http.ResponseWriter, dbr models.DBResult) {
	psr := dbr.(models.PostSessionDBResult)

	var foundHash models.Password
	err := psr.Row.Scan(&foundHash)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Authentication failed: user doesn't exist"))  
		return 
	}

	validHash := foundHash.String()

	errComparingHashes := bcrypt.CompareHashAndPassword(
		[]byte(validHash), 
		[]byte(psr.Password.String()),
	)

	if errComparingHashes != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authentication failed: wrong password"))
		return
	}

	jwtoken := auth.IssueJWT(psr.Username.String()) 
	jwtCookie := auth.BakeCookieWithJWT(jwtoken)

	http.SetCookie(w, jwtCookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authentication has been passed successfully"))
}


