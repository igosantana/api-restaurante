package initializers

import "api-restaurante/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
