package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//middleware session store intialization
	store := cookie.NewStore([]byte("key"))
	r.Use(sessions.Sessions("mysession", store))

	// Middleware login and logout
	r.Use(LogginMiddleware())

	// public routes
	r.POST("/login", LoginHandler)
	r.POST("/logout", LogoutHandler)

	// Protected routes require auth
	protected := r.Group("/api")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/dashboard", DashboardHandler)
		protected.GET("/profile", ProfileHandler)
	}

	fmt.Println("server running ")
	r.Run(":8080")
}

// loginhandler func for login and set session cookie
func LoginHandler(c *gin.Context) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	if creds.Username == "tibin" && creds.Password == "1234" {
		session := sessions.Default(c)
		session.Set("user", creds.Username)
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "login successfull"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login"})
	}
}

// Logouthandler func for logout and clear the session and cookie
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not loggedin"})
		return
	}
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "loggedout successfully"})
}

// Dashboardhandler only accessible after login
func DashboardHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("welcome to your dashboard,%v!", user)})
}

// profile handler func protected
func ProfileHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"profile": fmt.Sprintf("User: %v|Status:Active", user),
	})
}

///////////MiddleWare/////////////

//Logginmiddleware for login and logout

func LogginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/login" || path == "/logout" {
			fmt.Printf("[LOG]%s request made to %s\n", c.Request.Method, path)
		}
		c.Next()
	}
}

// Authroute func for ensure the user is logged in before accessing the protected route

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized access"})
			c.Abort()
		}
		c.Next()
	}
}
