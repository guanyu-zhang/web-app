package initializers

import "web-app-backend/models"

func SyncDatabase() {
	ConnectToDb()
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Message{})
	DB.AutoMigrate(&models.Group{})
}
