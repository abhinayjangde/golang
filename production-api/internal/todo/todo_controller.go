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

func GetTodoByIDController(c *gin.Context) {
	id := c.Param("id")

	todo, err := GetTodoById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func UpdateTodoController(c *gin.Context) {
	id := c.Param("id")

	var todo Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := UpdateTodoById(id, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo updated"})
}

func DeleteTodoController(c *gin.Context) {
	id := c.Param("id")

	err := DeleteTodoById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
