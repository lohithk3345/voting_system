package userServices

import (
	"slices"

	"github.com/lohithk3345/voting_system/helpers"
	userRepository "github.com/lohithk3345/voting_system/internal/repositories/user"
	"github.com/lohithk3345/voting_system/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServices struct {
	userRepository *userRepository.UserRepository
}

func NewUserService(database *mongo.Database) *UserServices {
	return &UserServices{
		userRepository: userRepository.NewUserRepo(database),
	}
}

func (u UserServices) CreateUser(user *types.User) (types.ID, error) {
	if user.Role == types.EmptyString || !slices.Contains(types.Roles, user.Role) {
		user.Role = types.VOTER
	}
	return u.userRepository.InsertUser(user)
}

// func (u UserServices) FindById(id types.ID) (*types.User, error) {
// 	return u.userRepository.FindUserByID(id)
// }

func (u UserServices) GetRoleById(id types.ID) (types.Role, error) {
	user, err := u.userRepository.FindUserByID(id)
	if err != nil {
		return types.EmptyString, err
	}
	return user.Role, nil
}

func (u UserServices) FindUserByEmail(email types.Email) (*types.User, error) {
	filter := &helpers.Filter{}
	return u.userRepository.FindOneByFilter(filter.ByEmail(email))
}

func (u UserServices) FindUserById(id types.ID) (*types.User, error) {
	filter := &helpers.Filter{}
	return u.userRepository.FindOneByFilter(filter.ByID(id))
}
