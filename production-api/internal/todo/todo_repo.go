package todo

import (
	"context"

	"github.com/abhinayjangde/prod/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetById(id string) (Todo, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Todo{}, err
	}

	var todo Todo
	err = getCollection().FindOne(context.Background(), bson.M{"_id": objID}).Decode(&todo)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func UpdateTodo(id string, update Todo) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	updateDoc := bson.M{
		"$set": bson.M{
			"title": update.Title,
			"done":  update.Done,
		},
	}
	_, err = getCollection().UpdateOne(context.Background(), bson.M{"_id": objID}, updateDoc)
	return err

}

func DeleteTodo(id string) error {
	collection := getCollection()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	return err
}
