package main

import (
	"log"
	"net/http"
	"time"

	"github.com/abhinayjangde/todo/internal/config"
	"github.com/abhinayjangde/todo/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	cfg, err := config.MustLoad()

	if err != nil {
		log.Fatal("Failed to load configuration: " + err.Error())
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect to database: " + err.Error())
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()

	router.SetTrustedProxies(nil)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "api is working",
			"status":    "success",
			"timestamp": time.Now().UTC(),
		})
	})

	router.Run(":" + cfg.Port)
}
