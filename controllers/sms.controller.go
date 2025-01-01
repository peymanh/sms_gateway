package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peymanh/sms_gateway/models"
	"github.com/peymanh/sms_gateway/services"

	"gorm.io/gorm"
)

const SMS_COST = 10

type SMSController struct {
	DB         *gorm.DB
	SMSService *services.SMSService
}

func NewSMSController(DB *gorm.DB, smsService *services.SMSService) SMSController {
	return SMSController{DB, smsService}
}

func (pc *SMSController) SendSMS(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.SendInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// check balance
	if currentUser.Balance >= SMS_COST {

		currentUser.Balance = currentUser.Balance - SMS_COST

		now := time.Now()
		newSMSLog := models.SMSLog{
			Body:      payload.Body,
			Language:  payload.Language,
			Receiver:  payload.Receiver,
			Cost:      SMS_COST,
			User:      currentUser,
			UserID:    currentUser.ID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		result := pc.DB.Create(&newSMSLog)
		if result.Error != nil {

			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
			return
		}

		// calling go routine to send SMS to the receiver asynchronously
		go pc.SMSService.SendSMS(ctx, &currentUser, &newSMSLog, payload.Receiver, payload.Body, payload.Language)

		ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newSMSLog})

	} else {
		ctx.JSON(http.StatusPaymentRequired, gin.H{"status": "error", "message": "Insufficeint balance"})
		return
	}
}

func (pc *SMSController) GetSMSLogs(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var logs []models.SMSLog
	results := pc.DB.Limit(intLimit).Offset(offset).Find(&logs, "user_id = ?", currentUser.ID)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(logs), "data": logs})
}
