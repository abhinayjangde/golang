package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/api/products/:id/:price", func(c *gin.Context) {
		id := c.Param("id")
		price := c.Param("price")

		c.JSON(http.StatusCreated, gin.H{
			"message": "user created",
			"id":      id,
			"price":   price,
		})
	})
	router.Run(":8080")
}
