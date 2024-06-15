package types

// import (
// 	buffers "ecommerce/buffers/userpb/protobuffs"
// 	"fmt"
// 	"math/rand"
// 	"time"

// 	"github.com/google/uuid"
// )

// type UserRequest struct {
// 	Name     string `json:"name"`
// 	Age      int    `json:"age"`
// 	Email    Email  `json:"email"`
// 	Address  string `json:"address"`
// 	IsActive bool   `json:"isActive"`
// 	Password string `json:"password"`
// 	Role     Role   `json:"role"`
// }

// func (u UserRequest) Convert() *User {
// 	return &User{
// 		Name:     u.Name,
// 		Age:      u.Age,
// 		Address:  u.Address,
// 		IsActive: u.IsActive,
// 		Email:    u.Email,
// 		Role:     u.Role,
// 	}
// }

// func ConvertUserRPCRequest(req *buffers.CreateUserRequest) *User {
// 	return &User{
// 		Name:    req.Name,
// 		Age:     int(req.Age),
// 		Address: req.Address,
// 		Email:   req.Email,
// 		Role:    req.Role,
// 		Phone:   req.Phone,
// 	}
// }

// type User struct {
// 	Id       UserID `bson:"_id,omitempty"`
// 	Name     string `bson:"name"`
// 	Age      int    `bson:"age"`
// 	Email    Email  `bson:"email"`
// 	Address  string `bson:"address"`
// 	Phone    string `bson:"phone"`
// 	Hash     string `bson:"hashPass"`
// 	IsActive bool   `bson:"isActive"`
// 	Role     Role   `bson:"role"`
// }

// func (u *User) SetID() {
// 	u.Id = uuid.New().String()
// }

// func (u *User) GetID() ID {
// 	return u.Id
// }

// func (u *User) GetName() string {
// 	return u.Name
// }

// func (u *User) AddHash(hash string) {
// 	u.Hash = hash
// }

// type UserID = ID

// // func (u *UserID) SetID() {
// // 	u = uuid.New().String()
// // }

// // func (u *UserID) GetID() ID {
// // 	return u.UUID
// // }

// // func (u *User) DecodeRaw(v *mongo.SingleResult) {
// // 	v.Decode(&u)
// // 	log.Println("DECODE:", u.Id)
// // 	uid := u.GetID()
// // 	u.Id = uid
// // }

// func GenerateRandomEmail() string {
// 	rand.Seed(time.Now().UnixNano())

// 	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 	randomString := make([]byte, 10)
// 	for i := range randomString {
// 		randomString[i] = allowedChars[rand.Intn(len(allowedChars))]
// 	}

// 	email := fmt.Sprintf("%s@%s", randomString, "example.domain")

// 	return email
// }
