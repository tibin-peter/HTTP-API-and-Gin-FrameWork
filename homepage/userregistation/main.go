package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Usename  string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var userstore = map[string]User{}

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("key"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/register", resisterhandler)
	r.POST("/login", loginhandler)
	r.POST("/logout", logouthandler)

	protected := r.Group("/")
	protected.Use(authentication())
	{
		protected.GET("/profile", profile)
		protected.GET("/dashboard", adminmiddleware(), dashboard)
	}

	fmt.Println("server running at the port 8080")
	r.Run(":8080")
}
func resisterhandler(c *gin.Context) {
	var newuser User
	if err := c.ShouldBindJSON(&newuser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	if newuser.Role != "admin" && newuser.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be 'admin' or 'user'"})
		return
	}
	if _, exists := userstore[newuser.Usename]; exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exist"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newuser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	newuser.Password = string(hash)
	userstore[newuser.Usename] = newuser

	c.JSON(http.StatusOK, gin.H{"message": "user registed successfully"})

}
func loginhandler(c *gin.Context) {
	var registeduser User

	if err := c.ShouldBindJSON(&registeduser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalide input"})
		return
	}
	hash, exist := userstore[registeduser.Usename]
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash.Password), []byte(registeduser.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}
	session := sessions.Default(c)
	session.Set("user", registeduser.Usename)
	session.Set("role", registeduser.Role)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
func logouthandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
func profile(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	role := session.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"profile": gin.H{
			"username": user,
			"role":     role,
			"status":   "Active",
		},
	})
}
func authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized please login first"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func dashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "welcome to dashboard"})
}
func adminmiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role := session.Get("role")

		if role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "access denided only admin has access to the dashboard"})
			c.Abort()
			return

		}
		c.Next()
	}
}
