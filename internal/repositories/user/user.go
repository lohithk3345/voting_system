package userRepository

import (
	"github.com/lohithk3345/voting_system/helpers"
	reporesult "github.com/lohithk3345/voting_system/internal/repositories/repo_result"

	"log"

	"github.com/lohithk3345/voting_system/types"

	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection string = "users"

type UserRepository struct {
	store UserStore[*types.User]
}

func NewUserRepo(database *mongo.Database) *UserRepository {
	return &UserRepository{
		store: NewUserDatabase[*types.User](database, userCollection),
	}
}

func (u UserRepository) InsertUser(user *types.User) (types.ID, error) {
	user.SetID()
	insertedID, err := u.store.insertOne(user)
	if err != nil {
		log.Println(err.(reporesult.StoreError).Message)
		return types.EmptyString, err
	}
	log.Println(insertedID)
	return insertedID, nil
}

func (u UserRepository) FindUserByID(id types.ID) (*types.User, error) {
	result, err := u.store.findOne(H.ByID(id))
	if err != nil {
		log.Println("UserRepository", err.(reporesult.StoreError).Message)
		return nil, err
	}

	var user *types.User
	result.Decode(&user)
	return user, err
}

// func (u UserRepository) FindUserByEmail(email types.Email) (*types.User, error) {
// 	log.Println("UserRepositoryEmailLog:", H.ByEmail(email))
// 	result, err := u.store.findOne(H.ByEmail(email))
// 	if err != nil {
// 		log.Println(err.(reporesult.StoreError).Message)
// 		return nil, err
// 	}
// 	var user *types.User
// 	result.Decode(&user)
// 	log.Println(user)
// 	return user, nil
// }

func (u UserRepository) FindOneByFilter(filter helpers.Filter) (*types.User, error) {
	result, err := u.store.findOne(filter.Doc())
	if err != nil {
		log.Println(err.(reporesult.StoreError).Message)
		return nil, err
	}
	var user *types.User
	result.Decode(&user)
	return user, nil
}

func (u UserRepository) UpdateOneByFilter(filter helpers.Filter, update helpers.UserUpdate) {
	result, err := u.store.updateOne(filter.Doc(), update.Doc())
	if err != nil {
		log.Println(err.(reporesult.StoreError).Message)
		return
	}
	log.Println(result)
}

func (u UserRepository) DeleteOneByFilter(filter helpers.Filter) {
	result, err := u.store.deleteOne(filter)
	if err != nil {
		log.Println(err.(reporesult.StoreError).Message)
		return
	}
	log.Println(result)
}
