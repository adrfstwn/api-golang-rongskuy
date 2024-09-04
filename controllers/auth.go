package controllers

import (
	"backend-rongskuy/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"unicode"

	"gopkg.in/gomail.v2"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("bajSBHJ344A35sfga0vaa!#44")

// AuthController operations
type AuthController struct {
	beego.Controller
}

type Claims struct {
	UserID int64  `json:"id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// SendEmail function for sending email
func SendEmail(to, subject, body string) error {
	// Konfigurasi SMTP server Anda
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUser := "rusterevolt@gmail.com"
	smtpPass := "wpxuwgbxrzasmgwd"

	// Buat pesan email
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Konfigurasi dialer untuk SMTP server
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Kirim email
	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent successfully to", to)
	return nil
}

// ValidatePassword checks if the password meets security requirements
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}

	return nil
}

func (c *AuthController) Login() {
	// Read the raw request body
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]string{
			"error": "Failed to read request body",
		}
		c.ServeJSON()
		return
	}

	// Convert the body to a string for logging
	bodyStr := string(body)
	fmt.Println("Request Body:", bodyStr)

	// Parse the JSON request body
	var loginData struct {
		UserID   int64  `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &loginData); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest) // 400 Bad Request
		c.Data["json"] = map[string]string{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		}
		c.ServeJSON()
		return
	}

	fmt.Printf("Parsed login data: %+v\n", loginData)

	// Fetch user from database
	user, err := models.GetUserModelByEmail(loginData.Email)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized) // 401 Unauthorized
		c.Data["json"] = map[string]string{
			"error": "Invalid email or password",
		}
		c.ServeJSON()
		return
	}

	fmt.Printf("Fetched user: %+v\n", user) // Log user object

	// Compare the provided password with the hashed password in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized) // 401 Unauthorized
		c.Data["json"] = map[string]string{
			"error": "Invalid email or password",
		}
		c.ServeJSON()
		return
	}

	// Password matched, generate JWT
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		UserID: user.Id,
		Email:  user.Email,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]string{
			"error": "Failed to generate token",
		}
		c.ServeJSON()
		return
	}

	// Respond with the token
	c.Ctx.Output.SetStatus(http.StatusOK) // 200 OK
	c.Data["json"] = map[string]interface{}{
		"token": tokenString,
		"role":  user.Role,
		"email": user.Email,
		"id":    user.Id,
	}
	c.ServeJSON()
}

func (c *AuthController) Register() {
	// Read the raw request body
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]string{
			"error": "Failed to read request body",
		}
		c.ServeJSON()
		return
	}

	// Convert the body to a string for logging
	bodyStr := string(body)
	fmt.Println("Request Body:", bodyStr)

	// Parse the JSON request body
	var v models.UserModel
	if err := json.Unmarshal(body, &v); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest) // 400 Bad Request
		c.Data["json"] = map[string]string{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Log the parsed request body for debugging
	parsedJSON, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println("Parsed JSON:", string(parsedJSON))

	// Check if the username already exists
	if existingUser, err := models.GetUserModelByName(v.Name); err == nil && existingUser != nil {
		c.Ctx.Output.SetStatus(http.StatusConflict) // 409 Conflict
		c.Data["json"] = map[string]string{
			"error": "Username already exists",
		}
		c.ServeJSON()
		return
	}

	// Check if the email already exists
	if existingUser, err := models.GetUserModelByEmail(v.Email); err == nil && existingUser != nil {
		c.Ctx.Output.SetStatus(http.StatusConflict) // 409 Conflict
		c.Data["json"] = map[string]string{
			"error": "Email already exists",
		}
		c.ServeJSON()
		return
	}

	// Validate the password
	if err := validatePassword(v.Password); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest) // 400 Bad Request
		c.Data["json"] = map[string]string{
			"error":   "Invalid password",
			"details": err.Error(),
		}
		c.ServeJSON()
		return
	}

	//hash pasword
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(v.Password), bcrypt.DefaultCost)
	v.Password = string(hashedPassword)

	// Set default role
	if v.Role == "" {
		v.Role = "user" // Default role
	}

	// Insert the new UserModel into the database
	if id, err := models.AddUserModel(&v); err == nil {
		c.Ctx.Output.SetStatus(http.StatusCreated) // 201 Created
		c.Data["json"] = map[string]interface{}{
			"id":      id,
			"success": "User successfully created",
		}
	} else {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]interface{}{
			"error":   "Failed to create user",
			"details": err.Error(),
		}
	}

	c.ServeJSON()
}

func (c *AuthController) ForgotPassword() {
	// Read the raw request body
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]string{
			"error": "Failed to read request body",
		}
		c.ServeJSON()
		return
	}

	// Parse the JSON request body
	var requestData struct {
		Email string `json:"email"`
	}

	if err := json.Unmarshal(body, &requestData); err != nil {
		log.Println("Invalid JSON format:", err)
		c.Ctx.Output.SetStatus(http.StatusBadRequest) // 400 Bad Request
		c.Data["json"] = map[string]string{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Fetch user from database
	user, err := models.GetUserModelByEmail(requestData.Email)
	if err != nil {
		log.Println("User not found:", err)
		c.Ctx.Output.SetStatus(http.StatusNotFound) // 404 Not Found
		c.Data["json"] = map[string]string{
			"error": "User not found",
		}
		c.ServeJSON()
		return
	}

	// Generate reset token and save it in the database
	resetToken, err := user.GenerateResetToken()
	if err != nil {
		log.Println("Failed to generate reset token:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]string{
			"error": "Failed to generate reset token",
		}
		c.ServeJSON()
		return
	}

	// Create reset email content
	resetURL := fmt.Sprintf("http://127.0.0.1:3000/reset-password?token=%s", resetToken)
	emailBody := fmt.Sprintf("Please reset your password using the following link: %s", resetURL)

	// Send reset token to user via email
	if err := SendEmail(requestData.Email, "Password Reset", emailBody); err != nil {
		log.Println("Failed to send reset token:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]string{
			"error": "Failed to send reset token",
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK) // 200 OK
	c.Data["json"] = map[string]string{
		"success": "Reset token sent to email",
	}
	c.ServeJSON()
}

func (c *AuthController) ResetPassword() {
	// Read the raw request body
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]string{
			"error": "Failed to read request body",
		}
		c.ServeJSON()
		return
	}

	// Convert body to string for logging
	bodyStr := string(body)
	fmt.Println("Request Body:", bodyStr)

	// Parse the JSON request body
	var requestData struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &requestData); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest) // 400 Bad Request
		c.Data["json"] = map[string]string{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Fetch user from database using the reset token
	user, err := models.GetUserModelByResetToken(requestData.Token)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized) // 401 Unauthorized
		c.Data["json"] = map[string]string{
			"error": "Invalid or expired token",
		}
		c.ServeJSON()
		return
	}

	// Validate the new password
	if err := validatePassword(requestData.Password); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest) // 400 Bad Request
		c.Data["json"] = map[string]string{
			"error":   "Invalid password",
			"details": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]string{
			"error": "Failed to hash password",
		}
		c.ServeJSON()
		return
	}
	user.Password = string(hashedPassword)

	// Update the password in the database
	if err := models.UpdateUserModelById(user); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]string{
			"error":   "Failed to update password",
			"details": err.Error(),
		}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK) // 200 OK
	c.Data["json"] = map[string]string{
		"success": "Password reset successfully",
	}
	c.ServeJSON()
}

func (c *AuthController) Logout() {
	// Ini adalah langkah sederhana untuk logout
	// Token JWT biasanya dihapus dari sisi klien, bukan dari server

	// Kirimkan respons sukses
	c.Ctx.Output.SetStatus(http.StatusOK) // 200 OK
	c.Data["json"] = map[string]string{
		"message": "Successfully logged out",
	}
	c.ServeJSON()
}
