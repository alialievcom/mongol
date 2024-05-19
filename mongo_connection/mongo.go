package mongo_connection

import (
	"alialiev/sites-core/constants"
	sites_core "alialiev/sites-core/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

func ConnectMongo(cfg *sites_core.Config) []*mongo.Collection {
	var collections = make([]*mongo.Collection, 0, len(cfg.Mongo.Collections))
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	collections = append(collections, client.Database(cfg.Mongo.DB).Collection(constants.UsersCollection))
	for _, v := range cfg.Mongo.Collections {
		collection := client.Database(cfg.Mongo.DB).Collection(v)
		if collection == nil {
			panic(fmt.Sprintf("collection: %s is nil", v))
		}
		if _, err := collection.EstimatedDocumentCount(ctx); err != nil {
			log.Fatalf("Failed to access collection %s: %v", v, err)
		}
		collections = append(collections, client.Database(cfg.Mongo.DB).Collection(v))
	}

	return collections
}
