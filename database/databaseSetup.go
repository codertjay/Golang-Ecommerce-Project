package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func DBSet() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("Failed to connect to mongo db")
		return nil
	}
	log.Println("Successfully connected to mongodb")
	return client
}

var Client *mongo.Client = DBSet()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Golang-Ecommerce").Collection(collectionName)
	return collection
}
