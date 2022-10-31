package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/giusepperoro/avitotest/internal/database/mocks"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCreateUserSuccess(t *testing.T) {

	reqData := CreateUserRequest{InitialBalance: 10000}
	data, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/create", bytes.NewReader(data))
	assert.NoError(t, err)

	dbMock := mocks.NewDbManager(t)
	dbMock.EXPECT().CreateUser(mock.Anything, int64(10000)).Return(1234, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleCreateUser(dbMock))
	handler.ServeHTTP(rr, req)

	respData := &CreateUserResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), respData)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, int64(1234), respData.ClientId)
}

func TestHandleCreateUserFailed(t *testing.T) {

	reqData := CreateUserRequest{InitialBalance: -10000}
	data, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/create", bytes.NewReader(data))
	assert.NoError(t, err)

	dbMock := mocks.NewDbManager(t)
	dbMock.EXPECT().CreateUser(mock.Anything, int64(-10000)).Return(1234, errors.New("negative balance"))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleCreateUser(dbMock))
	handler.ServeHTTP(rr, req)

	respData := &CreateUserResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), respData)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
