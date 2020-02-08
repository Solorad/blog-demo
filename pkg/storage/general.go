package storage

import (
	"github.com/solorad/blog-demo/server/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"golang.org/x/net/context"
	"os"
	"runtime/debug"
	"time"
)

const (
	bulkSize = 500
)

// MongoClient is an implementation of pkg.MongoClient
type MongoClient struct {
	Client *mongo.Client
}

// NewMongoClient init connection with MongoDb
func NewMongoClient() *MongoClient {
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	writeOptions := writeconcern.New(
		writeconcern.W(1),
		writeconcern.J(true),
	)
	// configuration was hardcoded here to make this demo app easy to run
	options := options.Client().SetWriteConcern(writeOptions).ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, options)
	if err != nil {
		log.Infof("Error occurred during mongoDb init %v", err)
		os.Exit(1)
	}
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Infof("Error occurred during mongoDb ping %v", err)
	}
	return &MongoClient{
		Client: client,
	}
}

// GetCollection create collection if not present and returns it
func (m *MongoClient) GetCollection(dbName, collectionName string) *mongo.Collection {
	return m.Client.Database(dbName).Collection(collectionName)
}

// UpsertSingle func upsert 1 item
func (m *MongoClient) UpsertSingle(collection *mongo.Collection, item interface{}, filter bson.M) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	upsert := true
	_, err := collection.UpdateOne(
		ctx,
		filter,
		bson.M{"$set": item},
		&options.UpdateOptions{Upsert: &upsert, Collation: &options.Collation{
			Locale: "de",
		}},
	)
	return err
}

// CreateIndex add i not present index in MongoDb
func (m *MongoClient) CreateIndex(collection *mongo.Collection, keys bsonx.Doc, indexName string, opts *options.IndexOptions) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, e := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    keys,
		Options: opts,
	})
	if e != nil {
		log.Errorf("Error occurred on %v index creation: %v", indexName, e)
	} else {
		log.Infof("Ensure in index for %v", indexName)
	}
}

// BulkWrite bulk write objects into db
func (m *MongoClient) BulkWrite(writeModels []mongo.WriteModel, collection *mongo.Collection) {
	start := time.Now()
	defer func() {
		log.TimeTrack(start, "BulkWrite to "+collection.Name())
	}()
	jobs := make(chan []mongo.WriteModel, 10000)
	results := make(chan struct{}, 10000)
	for i := 0; i < 4; i++ {
		go startBulkWorker(jobs, results, collection)
	}
	for i := 0; i < (len(writeModels)-1)/bulkSize+1; i++ {
		from := bulkSize * i
		to := bulkSize * (i + 1)
		if to > len(writeModels) {
			to = len(writeModels)
		}
		jobs <- writeModels[from:to]
	}
	close(jobs)
	for i := 0; i < (len(writeModels)-1)/bulkSize+1; i++ {
		<-results
	}
	close(results)
	debug.FreeOSMemory()
}

func startBulkWorker(jobs chan []mongo.WriteModel, results chan struct{}, collection *mongo.Collection) {
	for job := range jobs {
		writeBulk(collection, job)
		results <- struct{}{}
	}
}

func writeBulk(collection *mongo.Collection, writeModels []mongo.WriteModel) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	_, err := collection.BulkWrite(ctx, writeModels)
	if err != nil {
		log.Errorf("error occurred on bulk write in %s: %v", collection.Name(), err)
	}
}
