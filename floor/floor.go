package floor

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("floor", floor)
}

type request struct {
	Input string `json:"input"`
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
	inputVal, err := strconv.ParseFloat(req.Input, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Fatalf("strconv.ParseFloat failed; %s; %v", req.Input, err.Error())
	}
	result := math.Floor(inputVal)
	resp := fmt.Sprintf("%f", result)
	log.Printf("response: %s\n", resp)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resp))
}
