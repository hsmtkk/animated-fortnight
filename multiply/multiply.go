package multiply

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("multiply", multiply)
}

type request struct {
	Input int `json:"input"`
}

type response struct {
	Multiplied int `json:"multiplied"`
}

func multiply(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Fatalf("io.ReadAll failed; %v", err.Error())
	}
	var req request
	if err := json.Unmarshal(reqBytes, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Fatalf("json.Unmarshal failed; %s; %v", string(reqBytes), err.Error())
	}
	respBytes, err := json.Marshal(response{Multiplied: 2 * req.Input})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Fatalf("json.Marshal failed; %v", err.Error())
	}
	log.Printf("response: %v\n", string(respBytes))
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}
