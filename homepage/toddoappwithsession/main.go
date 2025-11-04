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
	store := cookie.NewStore([]byte("key"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/addtoddo", addtoddo)

	r.GET("/gettoddo", gettoddo)

	r.DELETE("/deletetoddo", deletetoddo)

	fmt.Println("server running at port 8080")
	r.Run(":8080")
}
func addtoddo(c *gin.Context) {
	session := sessions.Default(c)

	var req struct {
		Item string `json:"item"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	// getting the existing todo
	toddos := session.Get("toddos")
	var toddolist []string

	if toddos != nil {
		toddolist = toddos.([]string)
	}
	// add new toddo
	toddolist = append(toddolist, req.Item)
	//saving the section
	session.Set("toddos", toddolist)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"messate": "toddo added successfully"})
}
func gettoddo(c *gin.Context) {
	session := sessions.Default(c)

	toddos := session.Get("toddos")

	if toddos == nil {
		c.JSON(http.StatusOK, gin.H{"toddos": []string{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"toddos": toddos})
}
func deletetoddo(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("toddos")
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "all toddos clear for the section"})
}
