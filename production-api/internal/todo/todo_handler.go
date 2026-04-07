package todo

func CreateTodo(todo Todo) error {
	return Create(todo)
}

func GetTodos() ([]Todo, error) {
	return GetAll()
}

func GetTodoById(id string) (Todo, error) {
	return GetById(id)
}

func UpdateTodoById(id string, update Todo) error {
	return UpdateTodo(id, update)
}

func DeleteTodoById(id string) error {
	return DeleteTodo(id)
}
