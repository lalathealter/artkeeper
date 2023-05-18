package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/lalathealter/artkeeper/config"
	"golang.org/x/crypto/bcrypt"
)

func BcryptString(input string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
	}
	return string(hash)
}

const MAX_TOKEN_AGE = 60

type HeaderJWT struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type PayloadJWT struct {
	IAT int64 `json:"iat"`
	EXP int64 `json:"exp"`
	Name string `json:"name"`
}

func encodeJWTData(v any) string {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		log.Panicln(err)
	}
	encodedString := base64.RawURLEncoding.EncodeToString(jsonBytes)
	return encodedString
}

const JWT_ALG = "HS256"
const JWT_TYP = "JWT"
func IssueJWT(name string) string {
	
	header := HeaderJWT{JWT_ALG, JWT_TYP}

	currTime := time.Now().Unix()
	payload := PayloadJWT{
		currTime,
		currTime + MAX_TOKEN_AGE,
		name,
	}

	headerStr := encodeJWTData(header)
	payloadStr := encodeJWTData(payload)
	bodyStr := headerStr + "." + payloadStr

	signature := signJWT(bodyStr)
	signedJWT := bodyStr + "." + signature
	return signedJWT
}

var jwtSignSecret = []byte(config.Getnonempty("JWTSECRET"))
func signJWT(body string) string {
	hasher := hmac.New(sha256.New, jwtSignSecret)
	hasher.Write([]byte(body))
	signature := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
	return signature
}
