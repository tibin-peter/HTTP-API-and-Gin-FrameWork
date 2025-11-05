package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Book struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}
type User struct {
	Username string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var (
	userstore = map[string]User{} //store username and hashpassword
	bookstore = map[int]Book{}    // store book id
	nextid    = 1
)

func main() {

	r := gin.Default()
	store := cookie.NewStore([]byte("key"))
	r.Use(sessions.Sessions("mysession", store))

	// public routes
	r.POST("/signup", signup)
	r.POST("/login", login)
	r.POST("/logout", logout)

	// protected routes — only logged-in users can access
	protected := r.Group("/api")
	protected.Use(authmiddleware())
	{
		protected.GET("/books", listallbook)
		protected.GET("/books/:id", getbook)
	}

	// admin-only routes — only admin role can access
	admin := r.Group("/admin")
	admin.Use(authmiddleware(), adminmiddleware())
	{
		admin.POST("/books", addbook)
	}

	fmt.Println("server running at port 8080")
	r.Run(":8080")
}

// function for signup
func signup(c *gin.Context) {
	var newuser User

	if err := c.ShouldBindJSON(&newuser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if newuser.Username == "" || newuser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}
	if newuser.Role == "" {
		newuser.Role = "user"
	}
	if newuser.Role != "user" && newuser.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role must be user or admin"})
		return
	}
	if _, exist := userstore[newuser.Username]; exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already existing"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newuser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	userstore[newuser.Username] = User{
		Username: newuser.Username,
		Password: string(hash),
		Role:     newuser.Role,
	}
	c.JSON(http.StatusOK, gin.H{"message": "singup successfull"})

}

// func for login
func login(c *gin.Context) {
	var req struct {
		Username string `json:"name"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	user, ok := userstore[req.Username]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}
	session := sessions.Default(c)
	session.Set("user", user.Username)
	session.Set("role", user.Role)
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "login successfull",
		"user":    user.Username,
		"role":    user.Role,
	})
}

// func for logout
func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// func for add a new book
func addbook(c *gin.Context) {
	var newbook Book
	if err := c.ShouldBindJSON(&newbook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}
	if newbook.Author == "" || newbook.Name == "" || newbook.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and author is required and the price should be greaterthan zero"})
		return
	}
	id := nextid
	nextid++
	book := Book{Id: id, Name: newbook.Name, Author: newbook.Author, Price: newbook.Price}
	bookstore[id] = book
	c.JSON(http.StatusOK, gin.H{"message": "book added successfully"})

}

// func for listall book
func listallbook(c *gin.Context) {
	book := []Book{}
	for _, b := range bookstore {

		book = append(book, b)
	}
	c.JSON(http.StatusOK, gin.H{"Books": book})
}

// func for get a specific from its id
func getbook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}
	book, ok := bookstore[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Book": book})
}

////////////  Middleware  /////////////

// Authentication middleware
func authmiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized please login"})
			c.Abort()
			return
		}
		c.Next()

	}
}

// Admin middleware
func adminmiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role := session.Get("role")
		if role == nil || role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin access only"})
			c.Abort()
			return
		}
		c.Next()
	}

}
