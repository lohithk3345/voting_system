package database

import (
	"context"
	"log"

	"github.com/lohithk3345/voting_system/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func configOptions() *options.ClientOptions {
	connString := config.EnvMap["MONGODB_URL"]
	ioptions := options.Client().ApplyURI(connString)
	return ioptions
}

func connect() *mongo.Client {
	ioptions := configOptions()
	client, err := mongo.Connect(context.Background(), ioptions)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
