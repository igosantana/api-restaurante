package controllers

import (
	"api-restaurante/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	userMock := CreateUserMock()
	fmt.Println("userMock", userMock)
	defer DeleteUserMock(userMock.ID.String())
	r := SetUpRouter()
	r.POST("/login", AuthController.Login)
	userLogin := models.UserLogin{Email: userMock.Email, Password: "123456"}
	jsonValue, err := json.Marshal(userLogin)
	if err != nil {
		fmt.Println("marshal error")
	}
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("new request error")
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestLogout(t *testing.T) {
	r := SetUpRouter()
	r.GET("/logout", AuthController.LogoutUser)
	req, err := http.NewRequest("GET", "/logout", nil)
	if err != nil {
		fmt.Println("new request error")
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}
