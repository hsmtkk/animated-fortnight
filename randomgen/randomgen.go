package randomgen

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("randomgen", randomGen)
}

type response struct {
	Random int `json:"random"`
}

func randomGen(w http.ResponseWriter, r *http.Request) {
	value := rand.Intn(100) + 1
	respBytes, err := json.Marshal(response{Random: value})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Fatalf("json.Marshal failed; %v", err.Error())
	}
	log.Printf("response: %v\n", string(respBytes))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
