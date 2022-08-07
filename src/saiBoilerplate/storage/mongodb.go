package storage

import (
	"context"
	"fmt"

	"github.com/webmakom-com/saiBoilerplate/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	Collection *mongo.Collection
}

func GetStorageInstance(ctx context.Context, cfg *config.Configuration) (*Storage, *mongo.Client, error) {

	// use mongodb as a storage
	mongoClientOptions := &options.ClientOptions{}
	if cfg.Mongo.User != "" && cfg.Mongo.Pass != "" {
		mongoClientOptions = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.Mongo.Host, cfg.Mongo.Port)).SetAuth(options.Credential{
			Username: cfg.Mongo.User,
			Password: cfg.Mongo.Pass,
		})
	} else {
		mongoClientOptions = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.Mongo.Host, cfg.Mongo.Port))
	}
	client, err := mongo.Connect(ctx, mongoClientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("error when connect to mongo : %w", err)

	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error when ping mongo instance : %w", err)

	}

	mongoCollection := client.Database(cfg.Mongo.Database).Collection(cfg.Mongo.Collection)

	fmt.Printf("found collection : %s", mongoCollection.Name())

	return &Storage{
		Collection: mongoCollection,
	}, client, nil
}
