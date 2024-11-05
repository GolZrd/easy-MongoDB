package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongo() (*mongo.Client, error) {
	// Строка подключения MongoDB
	connectStr := "mongodb://localhost:27017"

	// Создаем клиента MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectStr))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println("Connected to MongoDB")

	return client, nil
}
