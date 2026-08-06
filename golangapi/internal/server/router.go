package server

import (
	"net/http"
	"time"

	"github.com/abhinayjangde/todo/internal/notes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(database *mongo.Database) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"ok":        true,
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
		})
	})

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome to the Todo API! Use /health to check the health of the server.")
	})

	notes.RegisterRoutes(r, database)

	return r
}
