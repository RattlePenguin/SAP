package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/RattlePenguin/SAP/models"
	"gorm.io/gorm"

	// "github.com/pquerna/otp/totp"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

// Registers a user into the system
func (ac *AuthController) RegisterUser(ctx *gin.Context) {
	var payload *models.RegisterUserInput

	// Error with payload binding JSON
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	newUser := models.User{
		Name:		payload.Name,
		Email:		strings.ToLower(payload.Email),
		Password:	payload.Password,
	}

	result := ac.DB.Create(&newUser)

	// Error Handling
	// Duplicate unique key value error
	if result.Error != nil && strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Email is already in use, please use another email address."})
		return
	// Other errors
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Registered successfully!"})
}