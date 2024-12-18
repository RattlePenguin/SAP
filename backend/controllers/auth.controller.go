package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/RattlePenguin/SAP/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/pquerna/otp/totp"
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

	// SQL Injection Site
	var user models.User
	query := fmt.Sprintf("SELECT * FROM users WHERE users.email = \"%s\" AND users.password = \"%s\" ORDER BY users.id", strings.ToLower(payload.Email), payload.Password)

	// Raw SQL
	result := ac.DB.Raw(query).First(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Returns key in base32 and URL encoded
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "rattlepenguin.com",
		AccountName: "admin@admin.com",
		SecretSize:  15,
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
		Otp_secret:   key.Secret(),
		Otp_auth_url: key.URL(),
	}

	ac.DB.Model(&user).Updates(dataToUpdate)

	otpResponse := gin.H{
		"base32":      key.Secret(),
		"otpauth_url": key.URL(),
	}
	ctx.JSON(http.StatusOK, otpResponse)
}

// Verify TOTP
func (ac *AuthController) VerifyOTP(ctx *gin.Context) {
	var payload *models.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "Token is invalid or user does not exist."

	var user models.User
	result := ac.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": message})
		return
	}

	valid := totp.Validate(payload.Token, user.Otp_secret)
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": message})
		return
	}

	dataToUpdate := models.User{
		Otp_enabled:  true,
		Otp_verified: true,
	}

	ac.DB.Model(&user).Updates(dataToUpdate)

	userResponse := gin.H{
		"id":          user.ID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"otp_enabled": user.Otp_enabled,
	}
	ctx.JSON(http.StatusOK, gin.H{"otp_verified": true, "user": userResponse})
}

// Validate OTP without modifying database
func (ac *AuthController) ValidateOTP(ctx *gin.Context) {
	var payload *models.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "Token is invalid or user does not exist."

	var user models.User
	result := ac.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": message})
		return
	}

	valid := totp.Validate(payload.Token, user.Otp_secret)
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"otp_valid": true})
}

// Disable TOTP for a user
func (ac *AuthController) DisableOTP(ctx *gin.Context) {
	var payload *models.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User does not exist!"})
		return
	}

	user.Otp_enabled = false
	ac.DB.Save(&user)

	userResponse := gin.H{
		"id":          user.ID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"otp_enabled": user.Otp_enabled,
	}
	ctx.JSON(http.StatusOK, gin.H{"otp_disabled": true, "user": userResponse})
}
