package main

import (
	"api-restaurante/initializers"
	"api-restaurante/routers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	routers.HandleRouters()
}
