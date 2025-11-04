package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var userStore = map[string]string{} // for username->hashded password
var sessions = map[string]string{}  // session token ->username

// func for check the input is valid or not
func validateInput(user User) error {
	if user.Username == "" || user.Password == "" {
		return fmt.Errorf("username and password cannot be empty")
	}
	if len(user.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	return nil
}

// Hash password before storing

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// compare password with hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// signup endpoint
func signupHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format"})
		return
	}
	if err := validateInput(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, exists := userStore[user.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"erro": "user name already exists"})
		return
	}

	hashed, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		return
	}
	userStore[user.Username] = hashed
	c.JSON(http.StatusOK, gin.H{"message": "signup successfull"})
}

// Login endpoint
func loginHandler(c *gin.Context) {
	var loginData User
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
		return
	}

	storedHash, exists := userStore[loginData.Username]
	if !exists || !checkPasswordHash(loginData.Password, storedHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Create a simple session (for demo purposes)
	sessionToken := fmt.Sprintf("%s_token", loginData.Username)
	sessions[sessionToken] = loginData.Username

	// Set cookie
	c.SetCookie("session_token", sessionToken, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

// Middleware: Authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if _, exists := sessions[sessionToken]; !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Protected route
func dashboardHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to your secure dashboard!"})
}

// Logout endpoint
func logoutHandler(c *gin.Context) {
	sessionToken, err := c.Cookie("session_token")
	if err == nil {
		delete(sessions, sessionToken)
	}
	c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// Main function
func main() {
	r := gin.Default()

	r.POST("/signup", signupHandler)
	r.POST("/login", loginHandler)

	protected := r.Group("/dashboard")
	protected.Use(AuthMiddleware())
	protected.GET("/", dashboardHandler)
	protected.POST("/logout", logoutHandler)

	fmt.Println("🚀 Server running on http://localhost:8080")
	r.Run(":8080")
}
