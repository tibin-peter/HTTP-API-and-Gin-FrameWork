package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

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
			c.JSON(http.StatusOK, gin.H{"message": "loggin successfully"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid"})
		}
	})
	r.Run(":8080")
}
