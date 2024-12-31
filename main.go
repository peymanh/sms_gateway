package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/peymanh/sms_gateway/handlers"
	"github.com/peymanh/sms_gateway/services"
)

func main() {
	db, err := services.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	smsService := services.NewSMSService(db)
	userService := services.NewUserService(db)

	smsHandler := handlers.NewSMSHandler(smsService, userService)

	router := mux.NewRouter()
	router.HandleFunc("/sms/send", smsHandler.SendSMS).Methods("POST")
	router.HandleFunc("/sms/log", smsHandler.GetSMSLog).Methods("GET")
	// ... (Add other routes)
	log.Print("helllooooo")
	log.Fatal(http.ListenAndServe(":8080", router))
}
