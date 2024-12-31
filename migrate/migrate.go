package main

import (
	"fmt"
	"log"

	"github.com/peymanh/sms_gateway/initializers"
	"github.com/peymanh/sms_gateway/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	initializers.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.SMSLog{})
	fmt.Println("👍 Migration complete")
}
