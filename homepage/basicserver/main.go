package main

import "github.com/gin-gonic/gin"

func main() {
	// Create a default Gin router
	r := gin.Default()

	// Basic homepage route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Gin Framework!",
		})
	})

	// Run the server on port 8080
	r.Run(":8080")
}
