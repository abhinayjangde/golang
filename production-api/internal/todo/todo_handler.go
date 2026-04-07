package todo

func CreateTodo(todo Todo) error {
	return Create(todo)
}

func GetTodos() ([]Todo, error) {
	return GetAll()
}
