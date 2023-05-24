package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lalathealter/artkeeper/config"
	"github.com/lalathealter/artkeeper/models"
	"golang.org/x/crypto/bcrypt"
)
func DecryptPassword(r *http.Request, pass *models.Password, clientNonce *models.Nonce) (string, error) {
	
	secretkey, e := GetServerNonce(r)
	if e != nil {
		return "", e
	}

	ciphblock, e := aes.NewCipher(secretkey)
	if e != nil {
		return "", e
	}

	aesgcm, e := cipher.NewGCM(ciphblock)
	if e != nil {
		return "", e
	}

	ciphtext := decodeHexedField(pass)
	cnonce := decodeHexedField(clientNonce)
	fmt.Println("CIPHER:", ciphtext)
	plainhash, e := aesgcm.Open(nil, cnonce, ciphtext, nil)
	if e != nil {
		return "", e
	}
	fmt.Println("decrypted", plainhash)

	return hex.EncodeToString(plainhash), nil
}

func decodeHexedField[T models.Stringlike](hexedInput *T) []byte {
	str, e := hex.DecodeString((*hexedInput).String())
	if e != nil {
		log.Panicln(e)
	}
	return str
}

func BcryptString(input string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
	}
	return string(hash)
}


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
func EncodeJWT(payload PayloadJWT) string {
	header := HeaderJWT{JWT_ALG, JWT_TYP}

	headerStr := encodeJWTData(header)
	payloadStr := encodeJWTData(payload)
	bodyStr := headerStr + "." + payloadStr

	signature := signJWT(bodyStr)
	signedJWT := bodyStr + "." + signature
	return signedJWT
}

func IssueJWT(name string) string {
	currTime := time.Now().Unix()
	payload := PayloadJWT{
		currTime, 
		currTime + MAX_TOKEN_AGE,
		name,
	}
	return EncodeJWT(payload)
}

const MAX_TOKEN_AGE = 60
const COOKIE_TOKEN_NAME = "token"
func BakeCookieWithJWT(jwt string) *http.Cookie {
	
	return &http.Cookie{
		Name: COOKIE_TOKEN_NAME,
		Value: jwt, 
		Path: "/",
		HttpOnly: true,
		Secure: true,
		MaxAge: MAX_TOKEN_AGE,
	}
}

func ValidateJWT(r *http.Request) error {
	jwtPassed, err := ParseJWTCookie(r)
	if err != nil {
		return err
	}

	jwtPayload, err := ParseLoadFromJWT(jwtPassed)
	if err != nil {
		return err
	}
	if jwtPayload.EXP < time.Now().Unix() {
		return errors.New("JWTOKEN OUTDATED")
	}

	jwtValid := EncodeJWT(jwtPayload)
	if jwtPassed != jwtValid {
		return errors.New("INVALID JWTOKEN")
	}

	return nil
}

func ParseJWTCookie(r *http.Request) (string, error) {
	jwtCookie, err := r.Cookie(COOKIE_TOKEN_NAME)
	if err != nil {
		return "", err
	}
	jwtStr := jwtCookie.Value
	return jwtStr, nil
}

func ParseLoadFromJWT(jwtString string) (PayloadJWT, error) {
	jwtObj := PayloadJWT{}

	hexPayload := strings.Split(jwtString, ".")[1]
	payload, err := base64.RawURLEncoding.DecodeString(hexPayload)
	err = json.Unmarshal(payload, &jwtObj)
	return jwtObj, err 
}

var jwtSignSecret = []byte(config.Getnonempty("JWTSECRET"))
func signJWT(body string) string {
	hasher := hmac.New(sha256.New, jwtSignSecret)
	hasher.Write([]byte(body))
	signature := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
	return signature
}
