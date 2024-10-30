package routes

import (
	"github.com/RattlePenguin/SAP/controllers"
	"github.com/gin-gonic/gin"
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
	router.POST("/login", rc.AuthController.LoginUser)
}
