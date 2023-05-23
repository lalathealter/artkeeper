package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"
)

const SNONCE_SIZE = 16
const REQID_SIZE = 8
const WAIT_TIME = 20 * time.Second
const HeaderAuthReqId = "Authentication-Request-ID"
const HeaderAuthServerNonce = "Authentication-Server-Nonce"
var AuthRequests = make(map[string][SNONCE_SIZE]byte)

func ServerNonceHandler(w http.ResponseWriter, r *http.Request) {
	snonce := generateRandomByteArr(SNONCE_SIZE)
	requestid := generateUniqueRequestId(REQID_SIZE)
	AuthRequests[requestid] = [SNONCE_SIZE]byte(snonce)
	hexSnonce := hex.EncodeToString(snonce)

	w.Header().Set(HeaderAuthReqId, requestid)
	w.Header().Set(HeaderAuthServerNonce, hexSnonce)
	w.WriteHeader(http.StatusOK)

	time.AfterFunc(WAIT_TIME, removeRequestId(requestid))
}

func generateUniqueRequestId(size int) string {
	idHexStr := hex.EncodeToString(generateRandomByteArr(size))
	_, alreadyExists := AuthRequests[idHexStr]
	if !alreadyExists {
		return idHexStr
	}
	return generateUniqueRequestId(size)
}

func generateRandomByteArr(size int) []byte {
	ranBytes := make([]byte, size)
	if _, err := rand.Read(ranBytes); err != nil {
		log.Panicln(err)
	}
	return ranBytes
}

func removeRequestId(idStr string) func() {
	return func() {
		_, isThere := AuthRequests[idStr]
		if isThere {
			delete(AuthRequests, idStr)
		}
	}
}

func GetServerNonce(r *http.Request) ([]byte, error) {
	reqid := r.Header.Get(HeaderAuthReqId)
	snonce, ok := AuthRequests[reqid]
	if !ok {
		return nil, fmt.Errorf("the server nonce for the provided id wasn't issued")
	}
	removeRequestId(reqid)()
	return snonce[:], nil
}
