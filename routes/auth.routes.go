package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/RattlePenguin/SAP/controllers"
)

type AuthRouteController struct {
	AuthController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("auth")

	router.POST("/register", rc.AuthController.RegisterUser)
}