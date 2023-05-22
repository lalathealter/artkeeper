package controllers

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/lalathealter/artkeeper/controllers/auth"
	"github.com/lalathealter/artkeeper/models"
)

var UserRegistrationHandler = factorAPIHandler(
	readUserRegistrationRequest,
	respondUserRegistration,
)

func readUserRegistrationRequest(r *http.Request) (models.Message, error) {
	msg, err := parseJSONMessage(r, models.RegisterUserRequest{})
	// TODO: replace with real key from server
	secretkey := make([]byte, 16)
	ciphblock, e := aes.NewCipher(secretkey)
	if e != nil {
		return nil, e
	}

	aesgcm, e := cipher.NewGCM(ciphblock)
	if e != nil {
		return nil, err
	}

	ciphtext := decodeHexedField(msg.Password)
	// TODO: replace with real nonce from client
	cnonce := make([]byte, 12)
	fmt.Println("CIPHER:", ciphtext)
	plainhash, e := aesgcm.Open(nil, cnonce, ciphtext, nil)
	if e != nil {
		return nil, e
	}
	fmt.Println("decrypted", plainhash)
	
	msg.Password.Update(func (_ string) string {
		return hex.EncodeToString(plainhash)
	})
	msg.Password.Update(auth.BcryptString)
	return msg, err 
}

func decodeHexedField[T models.Stringlike](hexedInput *T) []byte {
	str, e := hex.DecodeString((*hexedInput).String())
	if e != nil {
		log.Panicln(e)
	}
	return str
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
	
	jwt, err := parseJWT(r) 
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

func parseJWT(r *http.Request) (auth.PayloadJWT, error) {
	jwtObj := auth.PayloadJWT{}
	
	jwtCookie, err := r.Cookie(COOKIE_TOKEN_NAME)
	if err != nil {
		return jwtObj, err
	}
	hexPayload := strings.Split(jwtCookie.Value, ".")[1]
	payload, err := base64.RawURLEncoding.DecodeString(hexPayload)
	err = json.Unmarshal(payload, &jwtObj)
	return jwtObj, err 
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
