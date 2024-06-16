package middleware

import (
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/lohithk3345/voting_system/cache"
	"github.com/lohithk3345/voting_system/config"
	"github.com/lohithk3345/voting_system/helpers"
	"github.com/lohithk3345/voting_system/internal/auth"
	"github.com/lohithk3345/voting_system/types"
)

type UserMiddlewares struct {
	cache *cache.CacheService
}

func NewUserMiddlewares(cache *cache.CacheService) *UserMiddlewares {
	return &UserMiddlewares{
		cache: cache,
	}
}

func (m *UserMiddlewares) ApiKeyCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.GetHeader(types.API_KEY)
		if key == types.EmptyString && key != config.EnvMap["API_KEY"] {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *UserMiddlewares) CheckAccessTokenAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := helpers.Tokens.GetToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Token"})
			ctx.Abort()
			return
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		isLoggedIn, err := m.cache.GetTokens(claims.Id)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		if helpers.IsEmpty(isLoggedIn.Access) {
			log.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		log.Println("TOKEN", claims.Id)

		ctx.Set("userId", claims.Id)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}

func (m *UserMiddlewares) CheckRefreshTokenAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := helpers.Tokens.GetToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Token"})
			ctx.Abort()
			return
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		tokens, err := m.cache.GetTokens(claims.Id)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		if tokens.Refresh != token {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.Id)

		m.cache.DeleteTokenKey(claims.Id)
		ctx.Next()
	}
}

func (m *UserMiddlewares) CheckIfAdminRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		iRole := ctx.MustGet("role")
		role := iRole.(types.Role)

		if role != types.ADMIN || !slices.Contains(types.Roles, role) {
			log.Println("NOT A DEALER")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *UserMiddlewares) CheckIfVoterRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		iRole := ctx.MustGet("role")
		role := iRole.(types.Role)

		if role != types.VOTER || !slices.Contains(types.Roles, role) {
			log.Println("NOT A CUSTOMER")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
