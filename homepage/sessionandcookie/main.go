package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	// session store intialize
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/login", func(c *gin.Context) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if creds.Username == "tibin" && creds.Password == "1234" {
			session := sessions.Default(c)
			session.Set("user", creds.Username)
			session.Save()
			c.JSON(http.StatusOK, gin.H{"message": "login successfull"})
		}
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "logout successfull"})
	})

	r.Run(":8080")
}
