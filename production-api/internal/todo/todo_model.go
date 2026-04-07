package todo

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title string             `json:"title" bson:"title"`
	Done  bool               `json:"done" bson:"done"`
}
