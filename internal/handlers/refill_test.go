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

func TestHandleRefillSuccess(t *testing.T) {
	reqData := RefillRequest{ClientId: 1234, Amount: 3000}
	data, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/refill", bytes.NewReader(data))
	assert.NoError(t, err)

	dbMock := mocks.NewDbManager(t)
	dbMock.EXPECT().Refill(mock.Anything, int64(1234), int64(3000)).Return(true, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleRefill(dbMock))
	handler.ServeHTTP(rr, req)

	respData := &RefillResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), respData)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, true, respData.Approved)
}

func TestHandleRefillFailedAmount(t *testing.T) {
	reqData := RefillRequest{ClientId: 1234, Amount: -3000}
	data, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/refill", bytes.NewReader(data))
	assert.NoError(t, err)

	dbMock := mocks.NewDbManager(t)
	dbMock.EXPECT().Refill(mock.Anything, int64(1234), int64(-3000)).Return(false, errors.New("negative amount"))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleRefill(dbMock))
	handler.ServeHTTP(rr, req)

	respData := &RefillResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), respData)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, false, respData.Approved)
}

func TestHandleRefillFailedUser(t *testing.T) {
	reqData := RefillRequest{ClientId: 1234, Amount: 3000}
	data, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/refill", bytes.NewReader(data))
	assert.NoError(t, err)

	dbMock := mocks.NewDbManager(t)
	dbMock.EXPECT().Refill(mock.Anything, int64(1234), int64(3000)).Return(false, errors.New("unable to refill balance"))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleRefill(dbMock))
	handler.ServeHTTP(rr, req)

	respData := &RefillResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), respData)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, false, respData.Approved)
}
