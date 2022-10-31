package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/giusepperoro/avitotest/internal/database/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetBalanceSuccess(t *testing.T) {
	reqData := GetBalanceRequest{ClientId: 1234}
	data, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/balance", bytes.NewReader(data))
	assert.NoError(t, err)

	dbMock := mocks.NewDbManager(t)
	dbMock.EXPECT().GetBalance(mock.Anything, int64(1234)).Return(10000, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleGetBalance(dbMock))
	handler.ServeHTTP(rr, req)

	respData := &GetBalanceResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), respData)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, int64(10000), respData.Balance)
}

func TestHandleGetFailed(t *testing.T) {
	reqData := GetBalanceRequest{ClientId: 1234}
	data, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/balance", bytes.NewReader(data))
	assert.NoError(t, err)

	dbMock := mocks.NewDbManager(t)
	dbMock.EXPECT().GetBalance(mock.Anything, int64(1234)).Return(10000, errors.New("unable to create user"))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleGetBalance(dbMock))
	handler.ServeHTTP(rr, req)

	respData := &GetBalanceResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), respData)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, int64(0), respData.Balance)
}
