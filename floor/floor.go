package floor

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("floor", floor)
}

type request struct {
	Input float64 `json:"input"`
}

func floor(w http.ResponseWriter, r *http.Request) {
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
	result := math.Floor(req.Input)
	resp := fmt.Sprintf("%f", result)
	log.Printf("response: %s\n", resp)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resp))
}
