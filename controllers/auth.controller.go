package controllers

import (
	"net/http"
	"strings"

	"github.com/RattlePenguin/SAP/models"
	"github.com/gin-gonic/gin"
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
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: payload.Password,
	}

	// TODO: SQL Injection Site
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

// Logs in a user
func (ac *AuthController) LoginUser(ctx *gin.Context) {
	var payload *models.LoginUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// TODO: SQL Injection Site
	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid Email or Password"})
		return
	}

	userResponse := gin.H{
		"id":          user.ID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"otp_enabled": user.Otp_enabled,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "user": userResponse})
}

// Generate TOTP
func (ac *AuthController) GenerateOTP(ctx *gin.Context) {
	var payload *models.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message", err.Error()})
		return
	}

	// Returns key in base32 and URL encoded
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:			"rattlepenguin.com"
		AccountName:	"admin@admin.com"
		SecretSize:		15,
	})

	if err != nil {
		panic(err)
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid Email or Password"})
		return
	}

	dataToUpdate := models.User{
		Otp_secret:		key.Secret(),
		Otp_auth_url:	key.URL(),
	}

	
}