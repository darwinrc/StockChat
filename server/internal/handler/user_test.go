package handler

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/internal/model"
	mock_service "server/internal/service/mocks"
	"testing"
)

func TestHandleSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockUserService(ctrl)
	handler := NewUserHandler(mockService)

	user := &model.User{
		Username: "Bob",
		Password: "12345",
	}

	mockService.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(user, nil)

	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/signup", bytes.NewReader(userJSON))
	rec := httptest.NewRecorder()

	handler.HandleSignup(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rec.Code)
	}

	assert.Contains(t, rec.Body.String(), "\"username\":\"Bob\"")
}

func TestHandleLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockUserService(ctrl)
	handler := NewUserHandler(mockService)

	user := &model.User{
		Username: "Bob",
		Password: "12345",
	}

	mockService.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(user)

	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(userJSON))
	rec := httptest.NewRecorder()

	handler.HandleLogin(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rec.Code)
	}

	assert.Contains(t, rec.Body.String(), "\"username\":\"Bob\"")
}
