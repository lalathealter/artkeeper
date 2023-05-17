package controllers

import (
	"log"
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
	return msg, err
}

func respondPostSessionRequest(w http.ResponseWriter, dbr models.DBResult) {
	psr := dbr.(models.PostSessionDBResult)
	var foundHash models.Password
	err := psr.Row.Scan(&foundHash)
	if err != nil {
		log.Panicln(err)
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
	jwtCookie := bakeCookieWithJWT(jwtoken)

	http.SetCookie(w, jwtCookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authentication has been passed successfully"))
}


const COOKIE_TOKEN_NAME = "token"
func bakeCookieWithJWT(jwt string) *http.Cookie {
	
	return &http.Cookie{
		Name: COOKIE_TOKEN_NAME,
		Value: jwt, 
		Secure: true,
		MaxAge: auth.MAX_TOKEN_AGE,
		HttpOnly: true,
	}
}
