package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/peymanh/sms_gateway/controllers"
	"github.com/peymanh/sms_gateway/middleware"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("users")
	router.GET("/me", middleware.DeserializeUser(), uc.userController.GetMe)
	router.POST("/balance/update", middleware.DeserializeUser(), uc.userController.UpdateBalance)
}
