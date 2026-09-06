package notes

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Repo struct {
	coll *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		coll: db.Collection("notes"),
	}
}

func (r *Repo) Create(ctx context.Context, note Note) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.coll.InsertOne(opCtx, note)
	if err != nil {
		return Note{}, fmt.Errorf("Insert note failed")
	}

	note.ID = result.InsertedID.(primitive.ObjectID)
	return note, nil
}

func (r *Repo) List(ctx context.Context) ([]Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.coll.Find(opCtx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("Find notes failed")
	}
	defer cursor.Close(opCtx)

	var notes []Note
	for cursor.Next(opCtx) {
		var note Note
		if err := cursor.Decode(&note); err != nil {
			return nil, fmt.Errorf("Decode note failed")
		}
		notes = append(notes, note)
	}
	return notes, nil
}
