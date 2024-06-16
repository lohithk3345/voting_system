package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	middleware "github.com/lohithk3345/voting_system/apis/middlewares/user"
	"github.com/lohithk3345/voting_system/cache"
	"github.com/lohithk3345/voting_system/internal/voting"
)

type VoterAPIHandlers struct {
	cache *cache.CacheService
}

func NewVoterAPIHandlers(cache *cache.CacheService) *VoterAPIHandlers {
	return &VoterAPIHandlers{
		cache: cache,
	}
}

func (u *VoterAPIHandlers) SetupUserRouter() *gin.Engine {
	router := gin.Default()

	middlewares := middleware.NewUserMiddlewares(u.cache)

	router.Use(middlewares.ApiKeyCheck())

	manager, err := voting.Initialize(cache.NewCacheService())
	if err != nil {
		log.Panicln("Error In Initializing WebsocketManager")
	}

	router.GET("/ws", middlewares.CheckAccessTokenAuth(), manager.HandleWS)

	return router
}
