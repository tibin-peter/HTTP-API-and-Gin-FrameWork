package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/visits", homeHandler)

	fmt.Println("server running at port 8080")
	r.Run(":8080")
}
func homeHandler(c *gin.Context) {
	cookie, err := c.Cookie("visits")

	var count int

	if err != nil {
		count = 1
	} else {
		count, _ = strconv.Atoi(cookie)
		count++
	}
	c.SetCookie(
		"visits",
		strconv.Itoa(count),
		3600, // expires in one hour
		"/",  // for all the route
		"localhost",
		false, // this means not using the http
		true,  // http only
	)
	message := fmt.Sprintf("you have visited this site %d times", count)
	c.String(http.StatusOK, message)
}
