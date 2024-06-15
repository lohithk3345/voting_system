package rpcHandlers

import (
	"log"
	"strings"

	buffers "github.com/lohithk3345/voting_system/buffers/protobuffs"
	"github.com/lohithk3345/voting_system/cache"
	"github.com/lohithk3345/voting_system/internal/auth"
	userServices "github.com/lohithk3345/voting_system/services/user"
	"github.com/lohithk3345/voting_system/types"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	db *mongo.Database
	buffers.UnimplementedUserServiceServer
}

func NewUserServer(db *mongo.Database) *UserServer {
	return &UserServer{db: db}
}

func (s *UserServer) CreateUser(ctx context.Context, req *buffers.CreateUserRequest) (*buffers.CreateUserResponse, error) {
	newUser := types.ConvertUserRPCRequest(req)
	hash, errPass := auth.HashPassword([]byte(req.Password))
	if errPass != nil {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}
	newUser.AddHash(string(hash))
	_, err := userServices.NewUserService(s.db).CreateUser(newUser)
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "The Email Already Exists")
	}
	logger := log.Default()
	logger.Println("Success")

	return &buffers.CreateUserResponse{
		Status: 200,
		Body:   "User Created",
	}, nil
}

func (s *UserServer) GetUserById(ctx context.Context, req *buffers.GetUserByIdRequest) (*buffers.GetUserByIdResponse, error) {
	log.Println("User ID:", req.UserId)
	user, err := userServices.NewUserService(s.db).FindUserById(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User "+codes.NotFound.String())
	}
	if user.Role != types.VOTER {
		return nil, status.Error(codes.Canceled, "User "+codes.Canceled.String())
	}

	log.Println("Found Voter")

	return &buffers.GetUserByIdResponse{
		UserId: user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Age:    int32(user.Age),
	}, nil
}

func (s *UserServer) Login(ctx context.Context, req *buffers.UserLoginRequest) (*buffers.UserLoginResponse, error) {
	user, err := userServices.NewUserService(s.db).FindUserByEmail(req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User "+codes.NotFound.String())
	}

	errAuth := auth.VerifyHash([]byte(user.Hash), req.Password)
	if errAuth != nil {
		return nil, status.Error(codes.Unauthenticated, "The Email Or Password Is Wrong. Please Check Again")
	}

	accessToken, errToken := auth.GenerateAccessToken(user.Id, user.Role)
	if errToken != nil {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	refreshToken, errRToken := auth.GenerateRefreshToken(user.Id, user.Role)
	if errRToken != nil {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	cache.NewCacheService().SetTokens(user.Id, types.Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	})

	log.Println("Found Voter")

	return &buffers.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserServer) Logout(ctx context.Context, req *buffers.UserLogoutRequest) (*buffers.UserLogoutResponse, error) {

	data, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return &buffers.UserLogoutResponse{IsDone: buffers.Status_FAILED}, status.Error(codes.Unauthenticated, "Token Not Found")
	}

	log.Println(data)

	token := strings.Split(data["authorization"][0], " ")[1]
	log.Println(token)
	claims, err := auth.ValidateToken(token)
	if err != nil {

		return &buffers.UserLogoutResponse{IsDone: buffers.Status_FAILED}, status.Error(codes.Unauthenticated, "Token Invalid")
	}
	cache.NewCacheService().DeleteTokenKey(claims.Id)

	return &buffers.UserLogoutResponse{IsDone: buffers.Status_SUCCESS}, nil
}
