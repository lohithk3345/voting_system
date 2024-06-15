package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	userHandlers "github.com/lohithk3345/voting_system/apis/handlers"
	buffers "github.com/lohithk3345/voting_system/buffers/protobuffs"
	"github.com/lohithk3345/voting_system/cache"
	"github.com/lohithk3345/voting_system/config"
	"github.com/lohithk3345/voting_system/internal/database"
	"github.com/lohithk3345/voting_system/rpc/rpcHandlers"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func setupGRPC(db *mongo.Database, wg *sync.WaitGroup) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic("err")
	}

	s := grpc.NewServer()
	server := rpcHandlers.NewUserServer(db)
	buffers.RegisterUserServiceServer(s, server)
	log.Printf("Starting gRPC at: %s\n", config.EnvMap["GRPC_PORT"])
	s.Serve(lis)

	wg.Done()
}

func setupREST(db *mongo.Database, wg *sync.WaitGroup) {
	u := userHandlers.NewUserApiHandler(db, cache.NewCacheService())
	router := u.SetupUserRouter()

	log.Printf("Starting HTTP server at: %s\n", config.EnvMap["REST_API_PORT"])
	http.ListenAndServe(":3000", router)
	wg.Done()

}

func main() {
	db := database.NewClient().TestDatabase()

	var wg sync.WaitGroup

	wg.Add(1)
	go setupGRPC(db, &wg)
	wg.Add(1)
	go setupREST(db, &wg)

	wg.Wait()
}
