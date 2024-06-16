package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	middleware "github.com/lohithk3345/voting_system/apis/middlewares/user"
	"github.com/lohithk3345/voting_system/cache"
	"github.com/lohithk3345/voting_system/helpers"
	"github.com/lohithk3345/voting_system/internal/auth"
	userServices "github.com/lohithk3345/voting_system/services/user"
	"github.com/lohithk3345/voting_system/types"
	"go.mongodb.org/mongo-driver/mongo"
)

func (u *UserAPIHandlers) SetupUserRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	middlewares := middleware.NewUserMiddlewares(u.cache)

	router.Use(middlewares.ApiKeyCheck())

	router.POST("/register/voter", u.registerVoter)
	router.POST("/login", u.login)
	router.GET("/protected", middlewares.CheckAccessTokenAuth(), u.protected)
	router.GET("/refresh", middlewares.CheckRefreshTokenAuth(), u.refresh)
	router.GET("/logout", middlewares.CheckAccessTokenAuth(), u.logout)
	router.GET("/status", u.status)

	return router
}

type UserAPIHandlers struct {
	service *userServices.UserServices
	cache   *cache.CacheService
}

func NewUserApiHandler(database *mongo.Database, cache *cache.CacheService) *UserAPIHandlers {
	return &UserAPIHandlers{
		service: userServices.NewUserService(database),
		cache:   cache,
	}
}

func (u UserAPIHandlers) registerVoter(ctx *gin.Context) {
	var user *types.VoterRequest
	ctx.BindJSON(&user)

	if !(helpers.Validators.Email(user.Email)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please enter proper email"})
		return
	}

	newUser := user.Convert()

	hash, errPass := auth.HashPassword([]byte(user.Password))
	if errPass != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	newUser.AddHash(string(hash))

	_, err := u.service.CreateUser(newUser)
	if err != nil {
		ctx.JSON(409, gin.H{"error": "The Email Already Exists"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User Created"})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserAPIHandlers) login(ctx *gin.Context) {
	var data *LoginRequest
	ctx.BindJSON(&data)
	user, err := u.service.FindUserByEmail(data.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "The User is not registered to the app"})
		ctx.Abort()
		return
	}
	errAuth := auth.VerifyHash([]byte(user.Hash), data.Password)
	if errAuth != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "The Email Or Password Is Wrong. Please Check Again"})
		ctx.Abort()
		return
	}
	accessToken, errToken := auth.GenerateAccessToken(user.Id, user.Role)
	if errToken != nil {
		ctx.Abort()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	refreshToken, errRToken := auth.GenerateRefreshToken(user.Id, user.Role)
	if errRToken != nil {
		ctx.Abort()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	u.cache.SetTokens(user.Id, types.Tokens{Access: accessToken, Refresh: refreshToken})
	ctx.Header("Authorization", "Bearer "+accessToken)
	ctx.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "refreshToken": refreshToken})
}

func (u *UserAPIHandlers) logout(ctx *gin.Context) {
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
	u.cache.DeleteTokenKey(claims.Id)
}

func (u *UserAPIHandlers) protected(ctx *gin.Context) {
	id, _ := ctx.Get("userId")
	ctx.JSON(200, gin.H{"id": id})
}

func (u *UserAPIHandlers) refresh(ctx *gin.Context) {
	id := ctx.MustGet("userId").(types.ID)
	role, _ := u.service.GetRoleById(id)
	log.Println(id)
	accessToken, errToken := auth.GenerateAccessToken(id, role)
	if errToken != nil {
		ctx.Abort()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	refreshToken, errRToken := auth.GenerateRefreshToken(id, role)
	if errRToken != nil {
		ctx.Abort()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	u.cache.SetTokens(id, types.Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	})
	ctx.Header("Authorization", "Bearer "+accessToken)
	ctx.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "refreshToken": refreshToken})
}

func (u *UserAPIHandlers) status(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "up"})
}
