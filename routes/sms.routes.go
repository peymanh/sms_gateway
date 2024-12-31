package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/peymanh/sms_gateway/controllers"
	"github.com/peymanh/sms_gateway/middleware"
)

type SMSRouteController struct {
	smsController controllers.SMSController
}

func NewRouteSMSController(smsController controllers.SMSController) SMSRouteController {
	return SMSRouteController{smsController}
}

func (pc *SMSRouteController) SMSRoute(rg *gin.RouterGroup) {
	router := rg.Group("sms")
	router.Use(middleware.DeserializeUser())
	router.POST("/send", pc.smsController.SendSMS)
	router.GET("/logs", pc.smsController.GetSMSLogs)
}
