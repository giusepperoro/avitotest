package handlers

import (
	"encoding/json"
	"github.com/giusepperoro/avitotest/internal/database"
	"io"
	"log"
	"net/http"
)

type WithdrawalRequest struct {
	ClientId  int64 `json:"client_id"`
	ServiceId int64 `json:"service_id"`
	OrderId   int64 `json:"order_id"`
	Amount    int64 `json:"amount"`
}

type WithdrawalResponse struct {
	Success bool   `json:"success"`
	Err     string `json:"error,omitempty"`
}

func sendRequestWithdrawal(w http.ResponseWriter, status int, e error, errMSg string, approved bool) {
	response := WithdrawalResponse{
		Success: approved,
	}

	log.Printf("error handling create user: %v", e)

	if errMSg != "" {
		response.Err = errMSg
		w.WriteHeader(status)
	} else {
		response.Success = approved
	}
	rawData, err := json.Marshal(response)
	w.WriteHeader(status)

	_, err = w.Write(rawData)
	if err != nil {
		log.Printf("unable to write data: %v", err)
	}
}

func HandleWithdrawal(manager database.DbManager) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req WithdrawalRequest
		ctx := request.Context()

		if request.Method != "POST" {
			sendRequestWithdrawal(writer, http.StatusMethodNotAllowed, nil, "method not allowed", false)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			sendRequestWithdrawal(writer, http.StatusInternalServerError, err, "invalid body", false)
			return
		}

		err = json.Unmarshal(body, &req)
		if err != nil {
			sendRequestWithdrawal(writer, http.StatusBadRequest, err, "invalid body fields", false)
			return
		}

		success, err := manager.Withdrawal(ctx, req.ClientId, req.ServiceId, req.OrderId, req.Amount)
		if err != nil {
			sendRequestWithdrawal(writer, http.StatusBadRequest, err, "unable to create withdrawal", false)
			return
		}
		sendRequestWithdrawal(writer, http.StatusOK, nil, "", success)
	}
}
