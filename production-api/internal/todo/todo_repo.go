package todo

import (
	"context"

	"github.com/abhinayjangde/prod/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCollection() *mongo.Collection {
	return config.DB.Database("todo_db").Collection("todos")
}
func Create(todo Todo) error {
	_, err := getCollection().InsertOne(context.Background(), todo)
	return err
}

func GetAll() ([]Todo, error) {
	var todos []Todo

	cursor, err := getCollection().Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var todo Todo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}

	return todos, nil
}
