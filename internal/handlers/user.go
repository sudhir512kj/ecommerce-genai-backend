package handlers

import (
	"math/rand"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sudhir512kj/ecommerce_backend/config"
	"github.com/sudhir512kj/ecommerce_backend/internal/models"
	"github.com/sudhir512kj/ecommerce_backend/internal/repository"
	"github.com/sudhir512kj/ecommerce_backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo repository.UserRepository
	conf     *config.Config
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	conf := config.GetConfig()
	return &UserHandler{userRepo: userRepo, conf: conf}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    hashedPassword,
		Permissions: []models.Permission{"seller"},
	}

	if err := h.userRepo.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send email notification
	if err := h.sendEmailNotification(user.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.UserResponse{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Permissions: user.Permissions,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	otp, err := h.generateOTP(user.ID, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send OTP to user's email
	if err := h.sendOTPEmail(user.Email, otp.OTP); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to your email",
	})

	token, err := jwt.GenerateToken(h.conf, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func (h *UserHandler) VerifyOTP(c *gin.Context) {
	var req struct {
		OTP    string `json:"otp" binding:"required"`
		UserID int    `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otp, err := h.userRepo.GetOTPByUserID(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	if otp.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP has expired"})
		return
	}

	user, err := h.userRepo.GetUserByID(c.Request.Context(), otp.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log the user in and generate a JWT token
	token, err := jwt.GenerateToken(h.conf, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the used OTP
	if err := h.userRepo.DeleteOTP(otp.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully",
		"token":   token,
	})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
		UserID      int    `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetUserByID(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
		return
	}

	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = hashedPassword
	if err := h.userRepo.UpdateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	token, err := jwt.GenerateResetToken(h.conf, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send reset password email with the token
	if err := h.sendResetPasswordEmail(user.Email, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset instructions sent to your email",
	})

}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetUserByID(c.Request.Context(), c.GetInt("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.Permissions = req.Permissions

	if err := h.userRepo.UpdateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Permissions: user.Permissions,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (h *UserHandler) generateOTP(userID int, c *gin.Context) (*models.OTP, error) {
	otp := &models.OTP{
		UserID:    userID,
		OTP:       generateRandomOTP(),
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	}

	if err := h.userRepo.CreateOTP(c.Request.Context(), otp); err != nil {
		return nil, err
	}

	return otp, nil
}

func (h *UserHandler) sendOTPEmail(email, otp string) error {
	// Implement email sending logic
	from := h.conf.Email.From
	to := []string{email}
	subject := "Your OTP for login"
	body := "Dear user,\n\nYour one-time password (OTP) for login is: " + otp + "\n\nThis OTP will expire in 5 minutes.\n\nBest regards,\nThe Ecommerce Team"

	msg := "From: " + from + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	// auth := smtp.PlainAuth("", conf.Email.Username, conf.Email.Password, conf.Email.Host)
	err := smtp.SendMail(h.conf.Email.Host+":"+h.conf.Email.Port, nil, from, to, []byte(msg))
	return err
}

func (h *UserHandler) sendResetPasswordEmail(email, token string) error {
	// Implement email sending logic
	from := h.conf.Email.From
	to := []string{email}
	subject := "Reset your password"
	body := "Dear user,\n\nTo reset your password, please click on the following link:\n\n" +
		"https://your-app.com/reset-password?token=" + token + "\n\nThis link will expire in 1 hour.\n\nBest regards,\nThe Ecommerce Team"

	msg := "From: " + from + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	// auth := smtp.PlainAuth("", conf.Email.Username, conf.Email.Password, conf.Email.Host)
	err := smtp.SendMail(h.conf.Email.Host+":"+h.conf.Email.Port, nil, from, to, []byte(msg))
	return err
}

func (h *UserHandler) sendEmailNotification(email string) error {
	from := h.conf.Email.From
	to := []string{email}
	subject := "Welcome to our Ecommerce Platform"
	body := "Dear user,\n\nThank you for registering with our ecommerce platform. We're excited to have you on board!\n\nBest regards,\nThe Ecommerce Team"

	msg := "From: " + from + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	// auth := smtp.PlainAuth("", conf.Email.Username, conf.Email.Password, conf.Email.Host)
	err := smtp.SendMail(h.conf.Email.Host+":"+h.conf.Email.Port, nil, from, to, []byte(msg))
	return err
}

func generateRandomOTP() string {
	// Implement OTP generation logic
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	otp := make([]byte, 6)
	for i := range otp {
		otp[i] = digits[rand.Intn(len(digits))]
	}
	return string(otp)
}

func (h *UserHandler) AuthMiddleware(c *gin.Context) {
	// Extract the JWT token from the request
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
		c.Abort()
		return
	}

	// Verify the JWT token
	userId, err := jwt.VerifyToken(h.conf, tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// Set the user ID in the context
	c.Set("user_id", userId)
	c.Next()
}
