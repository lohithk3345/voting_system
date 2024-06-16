package database

import "go.mongodb.org/mongo-driver/mongo"

type ClientManager struct {
	client *mongo.Client
}

func NewClient() *ClientManager {
	return &ClientManager{
		client: connect(),
	}
}

func (store *ClientManager) Database(name string) *mongo.Database {
	return store.client.Database(name)
}

func (store *ClientManager) AppDatabase() *mongo.Database {
	return store.client.Database("app")
}

func (store *ClientManager) TestDatabase() *mongo.Database {
	return store.client.Database("test")
}
