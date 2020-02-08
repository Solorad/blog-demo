package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// BlogItem is an MongoDB object for blog persistence
type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}
