package controllers

import (
	"api-restaurante/controllers"
	"api-restaurante/initializers"
	"api-restaurante/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	m.Run()
}

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	return router
}

var Email string

func DeleteUserMock() {
	var user models.User
	initializers.DB.Delete(&user, "email = ?", Email)
}

func TestCreateUser(t *testing.T) {
	r := SetUpRouter()
	r.POST("/user", controllers.CreateUser)
	createUser := models.ToCreateUser{Name: "Test", Email: "test@email.com", Password: "123456"}
	jsonValue, err := json.Marshal(createUser)
	if err != nil {
		fmt.Println("marshal error")
	}
	Email = createUser.Email
	defer DeleteUserMock()
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("new request error")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}
