package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
}

func ConnectMongoDb(uri string) (*MongoClient, error) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return nil, err
	}
	fmt.Println("Connected to MongoDB")
	return &MongoClient{Client: client}, nil
}
