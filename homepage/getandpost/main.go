package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users []User

func main() {
	// creating a default server
	r := gin.Default()

	// creating a route group
	api := r.Group("/api")

	//get all users
	api.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, users)
	})

	//post a new users
	api.POST("/users", func(c *gin.Context) {
		var newUser User

		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		users = append(users, newUser)
		c.JSON(http.StatusCreated, gin.H{"message": "user added successfully", "user": newUser})
	})
	r.Run(":8080")

}
