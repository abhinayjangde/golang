package todo

import "github.com/gin-gonic/gin"

func TodoRoutes(router *gin.Engine) {
	router.POST("/todos", CreateTodoController)
	router.GET("/todos", GetTodosController)

	router.GET("/todos/:id", GetTodoByIDController)
	router.PUT("/todos/:id", UpdateTodoController)
	router.DELETE("/todos/:id", DeleteTodoController)
}
