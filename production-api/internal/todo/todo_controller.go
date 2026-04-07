package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTodoController(c *gin.Context) {
	var todo Todo

	// Bind JSON request
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := CreateTodo(todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo created"})
}

func GetTodosController(c *gin.Context) {
	todos, err := GetTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}
