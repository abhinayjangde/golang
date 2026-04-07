package main

import (
	"github.com/abhinayjangde/prod/internal/config"
	"github.com/abhinayjangde/prod/internal/todo"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	todo.TodoRoutes(r)

	r.Run(":8089")
}
