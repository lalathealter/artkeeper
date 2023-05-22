package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
)

const SNONCE_SIZE = 16
const HeaderAuthReqId = "Authentication-Request-ID"
const HeaderAuthServerNonce = "Authentication-Server-Nonce"
var AuthRequests = make(map[string][SNONCE_SIZE]byte)

func ServerNonceHandler(w http.ResponseWriter, r *http.Request) {
	snonce := generateRandomByteArr(SNONCE_SIZE)
	requestid := hex.EncodeToString(generateRandomByteArr(8))
	AuthRequests[requestid] = [SNONCE_SIZE]byte(snonce)
	hexSnonce := hex.EncodeToString(snonce)

	w.Header().Set(HeaderAuthReqId, requestid)
	w.Header().Set(HeaderAuthServerNonce, hexSnonce)
	w.WriteHeader(http.StatusOK)
}

func generateRandomByteArr(size int) []byte {
	ranBytes := make([]byte, size)
	if _, err := rand.Read(ranBytes); err != nil {
		log.Panicln(err)
	}
	return ranBytes
}


func GetServerNonce(r *http.Request) ([]byte, error) {
	reqid := r.Header.Get(HeaderAuthReqId)
	snonce, ok := AuthRequests[reqid]
	if !ok {
		return nil, fmt.Errorf("the server nonce for the provided id wasn't issued")
	}
	delete(AuthRequests, reqid)
	return snonce[:], nil
}
