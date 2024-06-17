package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type ClientManager struct {
	client *mongo.Client
}

func NewClient(ctx context.Context) *ClientManager {
	return &ClientManager{
		client: connect(ctx),
	}
}

func (c *ClientManager) Disconnect(ctx context.Context) {
	disconnect(c.client, ctx)
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
