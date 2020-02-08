package pkg

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// MongoClient interface for work with Mongo
type MongoClient interface {
	GetCollection(dbName, collectionName string) *mongo.Collection
	UpsertSingle(collection *mongo.Collection, item interface{}, filter bson.M) error
	CreateIndex(collection *mongo.Collection, keys bsonx.Doc, indexName string, opts *options.IndexOptions)
	BulkWrite(writeModels []mongo.WriteModel, collection *mongo.Collection)
}
