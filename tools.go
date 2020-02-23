package undao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
)

func MongoIndexesCreateMany(coll *mongo.Collection, keys []string) {
	models := make([]mongo.IndexModel, 0)
	for _, item := range keys {
		models = append(models, mongo.IndexModel{
			Keys:    bsonx.Doc{{Key: item, Value: bsonx.Int64(1)}},
			Options: options.Index().SetBackground(true),
		})
	}
	if _, err := coll.Indexes().CreateMany(context.Background(), models); err != nil {
		log.Printf("[  error  ] %s store MongoIndexesCreateMany(): %s\n", coll.Name(), err)
	}
	return
}

func MongoIndexesCreateUnique(coll *mongo.Collection, key string) {
	models := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: key, Value: bsonx.Int64(1)}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := coll.Indexes().CreateOne(context.Background(), models); err != nil {
		log.Printf("[  error  ] %s store MongoIndexesCreateMany(): %s\n", coll.Name(), err)
	}
	return
}
