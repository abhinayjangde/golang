package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	var router *gin.Engine = gin.Default()

	router.SetTrustedProxies(nil)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "api is working",
			"status":    "success",
			"timestamp": time.Now().UTC(),
		})
	})

	router.Run(":8080")
}
