package handlers

import (
	"encoding/json"
	"github.com/giusepperoro/avitotest/internal/database"
	"io"
	"log"
	"net/http"
)

type RefillRequest struct {
	ClientId int64 `json:"client_id"`
	Amount   int64 `json:"amount"`
}

type RefillResponse struct {
	Approved bool   `json:"approved"`
	Err      string `json:"error,omitempty"`
}

func sendRequestRefill(w http.ResponseWriter, status int, e error, errMSg string, approved bool) {
	response := RefillResponse{
		Approved: approved,
	}

	log.Printf("error handling refill: %v", e)

	if errMSg != "" {
		response.Err = errMSg
		w.WriteHeader(status)
	} else {
		response.Approved = approved
	}
	rawData, err := json.Marshal(response)
	w.WriteHeader(status)

	_, err = w.Write(rawData)
	if err != nil {
		log.Printf("unable to write data: %v", err)
	}
}

func HandleRefill(manager database.DbManager) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		req := new(RefillRequest)
		ctx := request.Context()

		if request.Method != "POST" {
			sendRequestRefill(writer, http.StatusMethodNotAllowed, nil, "method not allowed", false)
			return
		}

		if req.Amount < 0 {
			sendRequestRefill(writer, http.StatusBadRequest, nil, "negative amount", false)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			sendRequestRefill(writer, http.StatusInternalServerError, err, "invalid body", false)
			return
		}

		err = json.Unmarshal(body, req)
		if err != nil {
			sendRequestRefill(writer, http.StatusBadRequest, err, "invalid body fields", false)
			return
		}

		refill, err := manager.Refill(ctx, req.ClientId, req.Amount)
		if err != nil {
			sendRequestRefill(writer, http.StatusBadRequest, err, "unable to refill balance", false)
			return
		}
		sendRequestRefill(writer, http.StatusOK, nil, "", refill)
	}
}
