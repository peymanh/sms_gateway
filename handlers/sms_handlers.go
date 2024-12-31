package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/peymanh/sms_gateway/services"
)

type SMSHandler struct {
	SMSService  *services.SMSService
	UserService *services.UserService
}

func NewSMSHandler(smsService *services.SMSService, userService *services.UserService) *SMSHandler {
	return &SMSHandler{SMSService: smsService, UserService: userService}
}

func (h *SMSHandler) SendSMS(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request context (you'll need to implement this)
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Invalid or missing user ID", http.StatusUnauthorized)
		return
	}

	var request struct {
		To      string `json:"to"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.SMSService.SendSMS(string(userID), request.To, request.Message, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "SMS sent successfully"})
}

func (h *SMSHandler) GetSMSLog(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request context (you'll need to implement this)
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid or missing user ID", http.StatusUnauthorized)
		return
	}

	smsLog, err := h.SMSService.GetSMSLog(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(smsLog)
}
