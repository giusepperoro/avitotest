package handlers

import (
	"encoding/json"
	"github.com/giusepperoro/avitotest/internal/database"
	"io"
	"log"
	"net/http"
)

type CreateUserRequest struct {
	InitialBalance int64 `json:"initial_balance"`
}

type CreateUserResponse struct {
	ClientId int64  `json:"client_id,omitempty"`
	Err      string `json:"error,omitempty"`
}

func sendRequestCreate(w http.ResponseWriter, status int, e error, errMSg string, clientId int64) {
	response := CreateUserResponse{
		ClientId: clientId,
	}

	log.Printf("error handling create user: %v", e)

	if errMSg != "" {
		response.Err = errMSg
		w.WriteHeader(status)
	} else {
		response.ClientId = clientId
	}
	rawData, err := json.Marshal(response)
	w.WriteHeader(status)

	_, err = w.Write(rawData)
	if err != nil {
		log.Printf("unable to write data: %v", err)
	}
}

func HandleCreateUser(manager database.DbManager) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		var req CreateUserRequest

		if request.Method != "POST" {
			sendRequestCreate(writer, http.StatusMethodNotAllowed, nil, "method not allowed", 0)
			return
		}

		if req.InitialBalance < 0 {
			sendRequestCreate(writer, http.StatusBadRequest, nil, "negative balance", 0)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			sendRequestCreate(writer, http.StatusBadRequest, err, "invalid body", 0)
			return
		}

		err = json.Unmarshal(body, &req)
		if err != nil {
			sendRequestCreate(writer, http.StatusBadRequest, err, "invalid body fields", 0)
			return
		}

		userId, err := manager.CreateUser(ctx, req.InitialBalance)
		if err != nil {
			sendRequestCreate(writer, http.StatusBadRequest, err, "unable to create user", 0)
			return
		}
		sendRequestCreate(writer, http.StatusOK, nil, "", userId)
	}
}
