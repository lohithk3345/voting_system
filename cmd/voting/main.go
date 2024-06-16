package main

import (
	"net/http"

	handlers "github.com/lohithk3345/voting_system/apis/handlers/voting"
	"github.com/lohithk3345/voting_system/cache"
	"github.com/lohithk3345/voting_system/config"
)

func main() {
	c := cache.NewCacheService()

	v := handlers.NewVoterAPIHandlers(c)
	http.ListenAndServe(":"+config.EnvMap["VOTING_WEBSOCKET_PORT"], v.SetupUserRouter())
}
