// Product API

package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Products struct {
	Id    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"gt=0"`
}

var products []Products
var Id = 1

func main() {
	r := gin.Default()

	api := r.Group("/api")
	api.GET("/getall", getall)
	api.GET("/getbyid/:id", getproductbyid)

	api.POST("/add", addproduct)

	api.PUT("/update/:id", update)

	api.PATCH("/partial/:id",partialupdate)

	api.DELETE("/delete/:id", delete)

	fmt.Println("server running at port 8080")
	r.Run(":8080")
}

// func for get all product
func getall(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"prodcts": products})
}

// func for get a single product by id
func getproductbyid(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for _, p := range products {
		if p.Id == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

// func for add new product
func addproduct(c *gin.Context) {
	var newproduct Products

	if err := c.ShouldBindJSON(&newproduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	newproduct.Id = Id
	Id++
	products = append(products, newproduct)
	c.JSON(http.StatusCreated, gin.H{"message": "product added successfully"})
}

// func for update Product by id
func update(c *gin.Context) {
	var update Products
	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	for i, p := range products {
		if p.Id == id {
			update.Id = id
			products[i] = update
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "product updated successfully"})
}

// func for delete a product using id
func delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for i, p := range products {
		if p.Id == id {
			products = append(products[:i], products[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
			return
		}
	}
	c.JSON(http.StatusNoContent, gin.H{"error": "product not found"})
}
//func for partial update using the product id
func partialupdate(c*gin.Context){
	id,_:= strconv.Atoi(c.Param("id"))
	var partial map[string]interface{}

	if err:= c.ShouldBindJSON(&partial);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"invalid input or mehod"})
		return
	}
    
	for i,p := range products{
		if p.Id==id{
            if name,ok:= partial["name"].(string);ok{
				products[i].Name=name
			}
			if price,ok:= partial["name"].(int):ok{
				products[i].Price=price
			}
			c.JSON(http.StatusOK, gin.H{"message": "Product partially updated"})
			return
		}
	}
   c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}
