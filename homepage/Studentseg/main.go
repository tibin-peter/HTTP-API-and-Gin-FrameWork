package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Student struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Mark int    `json:"mark"`
}

var students []Student
var nextId = 1

func main() {
	r := gin.Default()

	r.GET("/", homepage)
	r.GET("/students", getstudents)
	r.GET("/students/:id", getstudentsbyid)
	r.GET("/top", topmark)

	r.POST("/add", addstudent)

	r.DELETE("/students/:id", delete)

	r.PUT("/update/:id", updatestudent)

	r.PATCH("/updatestudentpartial/:id", updatestudentpartial)

	fmt.Println("Server running at port:8080")
	r.Run(":8080")
}
func homepage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "welcome to home page"})
}
func getstudents(c *gin.Context) {
	c.JSON(http.StatusOK, students)
}
func getstudentsbyid(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for _, s := range students {
		if s.Id == id {
			c.JSON(http.StatusOK, s)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"Error": "student not found"})
}
func addstudent(c *gin.Context) {
	var newstudent Student
	if err := c.ShouldBindJSON(&newstudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalide request"})
		return
	}
	newstudent.Id = nextId
	nextId++
	students = append(students, newstudent)

	c.JSON(http.StatusOK, gin.H{"message": "student added successfully", "student": newstudent})
}
func delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for i, s := range students {
		if s.Id == id {
			students = append(students[:i], students[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("student delete in the id of %d", id)})
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"messgae": "student not found"})
}
func updatestudent(c *gin.Context) {
	var update Student
	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid format"})
		return
	}

	for i, s := range students {
		if s.Id == id {
			students[i] = update
			c.JSON(http.StatusOK, gin.H{"message": "student updated successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "students not found"})
}
func updatestudentpartial(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var update map[string]interface{}

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid input"})
		return
	}

	for i, s := range students {
		if s.Id == id {

			if name, ok := update["name"].(string); ok {
				students[i].Name = name
			}
			if mark, ok := update["mark"].(int); ok {
				students[i].Mark = mark
			}
			c.JSON(http.StatusOK, gin.H{"message": "update sucessfull"})
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
}
func topmark(c *gin.Context) {
	var topstudent []Student
	for _, s := range students {
		if s.Mark > 80 {
			topstudent = append(topstudent, s)
		}
	}
	c.JSON(http.StatusOK, gin.H{"topstudents": topstudent})
}
