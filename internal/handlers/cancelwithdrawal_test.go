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

func TestHandleCancelWithdrawalSuccess(t *testing.T) {
	reqData := CancelWithdrawalRequest{
		ClientId:  1234,
		ServiceId: 111,
		OrderId:   17,
		Amount:    3000,
	}
	data, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/cancelWithdrawal", bytes.NewReader(data))
	assert.NoError(t, err)

	dbMock := mocks.NewDbManager(t)
	dbMock.EXPECT().CancelWithdrawal(mock.Anything, int64(1234), int64(111), int64(17), int64(3000)).Return(true, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleCancelWithdrawal(dbMock))
	handler.ServeHTTP(rr, req)

	respData := &CancelWithdrawalResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), respData)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, true, respData.Success)
}

func TestHandleCancelWithdrawalFailed(t *testing.T) {
	reqData := CancelWithdrawalRequest{
		ClientId:  1234,
		ServiceId: 111,
		OrderId:   17,
		Amount:    3000,
	}
	data, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/cancelWithdrawal", bytes.NewReader(data))
	assert.NoError(t, err)

	dbMock := mocks.NewDbManager(t)
	dbMock.EXPECT().CancelWithdrawal(mock.Anything, int64(1234), int64(111), int64(17), int64(3000)).Return(false, errors.New("unable to process withdrawal"))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleCancelWithdrawal(dbMock))
	handler.ServeHTTP(rr, req)

	respData := &CancelWithdrawalResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), respData)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, false, respData.Success)
}
