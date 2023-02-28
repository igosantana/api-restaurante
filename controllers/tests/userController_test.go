package controllers

import (
	"api-restaurante/app"
	"api-restaurante/controllers"
	"api-restaurante/initializers"
	"api-restaurante/models"
	"api-restaurante/routers"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routers.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routers.UserRouterController
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routers.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routers.NewUserRouteController(UserController)
	server = gin.Default()
}

func SetUpRouter() *gin.Engine {
	server = gin.Default()
	return server
}

func TestMain(m *testing.M) {
	router := server.Group("/api")
	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
}

func CreateUserMock() models.User {
	var userMock models.User
	r := SetUpRouter()
	r.POST("api/user", UserController.CreateUser)
	userCreate := models.ToCreateUser{Name: "test", Email: "test@email.com", Password: "123456"}
	jsonValue, err := json.Marshal(userCreate)
	if err != nil {
		fmt.Println("marshal error")
	}
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("new request error")
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	initializers.DB.First(&userMock, "email = ?", userCreate.Email)
	return userMock
}

func DeleteUserMock(id string) {
	var user models.User
	fmt.Println(id)
	initializers.DB.First(&user, "id = ?", id)
	initializers.DB.Delete(&user)
}

func TestUserUpdate(t *testing.T) {
	userMock := CreateUserMock()
	defer DeleteUserMock(userMock.ID.String())
	r := SetUpRouter()
	r.PATCH("/user/:id", UserController.UpdateUser)
	userUpdate := models.UserUpdate{Name: "teste2", Email: "test2@email.com", Password: "7654321"}
	jsonValue, err := json.Marshal(userUpdate)
	if err != nil {
		fmt.Println("marshal error")
	}
	path := "/user/" + userMock.ID.String()
	req, err := http.NewRequest("PATCH", path, bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("new request error")
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestDeleteUser(t *testing.T) {
	userMock := CreateUserMock()
	r := SetUpRouter()
	r.DELETE("/user/:id", UserController.DeleteUser)
	path := "/user/" + userMock.ID.String()
	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		fmt.Println("new request error")
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetUser(t *testing.T) {
	userMock := CreateUserMock()
	defer DeleteUserMock(userMock.ID.String())
	r := SetUpRouter()
	r.GET("/user/:id", UserController.GetUser)
	path := "/user/" + userMock.ID.String()
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		fmt.Println("new request error")
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	var userRes app.User
	json.Unmarshal(res.Body.Bytes(), &userRes)
	assert.Equal(t, http.StatusOK, res.Code)
}
