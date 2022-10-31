package handlers

import (
	"encoding/json"
	"github.com/giusepperoro/avitotest/internal/database"
	"io"
	"log"
	"net/http"
)

type GetBalanceRequest struct {
	ClientId int64 `json:"client_id"`
}

type GetBalanceResponse struct {
	Balance int64  `json:"balance"`
	Err     string `json:"error,omitempty"`
}

func sendRequestBalance(w http.ResponseWriter, status int, e error, errMSg string, balance int64) {
	response := GetBalanceResponse{
		Balance: balance,
	}

	log.Printf("error handling create user: %v", e)

	if errMSg != "" {
		response.Err = errMSg
		w.WriteHeader(status)
	} else {
		response.Balance = balance
	}
	rawData, err := json.Marshal(response)
	w.WriteHeader(status)

	_, err = w.Write(rawData)
	if err != nil {
		log.Printf("unable to write data: %v", err)
	}
}

func HandleGetBalance(manager database.DbManager) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		req := new(GetBalanceRequest)
		ctx := request.Context()

		if request.Method != "GET" {
			sendRequestBalance(writer, http.StatusMethodNotAllowed, nil, "method not allowed", req.ClientId)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			sendRequestBalance(writer, http.StatusBadRequest, err, "invalid body", 0)
			return
		}

		err = json.Unmarshal(body, req)
		if err != nil {
			sendRequestBalance(writer, http.StatusBadRequest, err, "invalid body fields", 0)
			return
		}

		balance, err := manager.GetBalance(ctx, req.ClientId)
		if err != nil {
			sendRequestBalance(writer, http.StatusBadRequest, err, "unable to create user", 0)
			return
		}
		sendRequestBalance(writer, http.StatusOK, nil, "", balance)
	}
}
