package types

import (
	"github.com/google/uuid"
	buffers "github.com/lohithk3345/voting_system/buffers/protobuffs"
)

// User Types

type User struct {
	Id    ID     `bson:"_id,omitempty"`
	Name  string `bson:"name"`
	Age   uint8  `bson:"age"`
	Email Email  `bson:"email"`
	Hash  string `bson:"hashPass"`
	Role  Role   `bson:"role"`
}

func (d *User) SetID() {
	d.Id = uuid.New().String()
}

func (u *User) AddHash(hash string) {
	u.Hash = hash
}

func ConvertUserRPCRequest(req *buffers.CreateUserRequest) *User {
	return &User{
		Name:  req.Name,
		Age:   uint8(req.Age),
		Email: req.Email,
		Role:  VOTER,
		Id:    uuid.NewString(),
		// Hash: ,
	}
}

// Voter Types

type VoterID = ID

type Voter struct {
	User
}

func (v *Voter) SetID() {
	v.Id = uuid.New().String()
}

func (u *Voter) AddHash(hash string) {
	u.Hash = hash
}

type VoterRequest struct {
	Name     string   `json:"name"`
	Age      uint8    `json:"age"`
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

func (v *VoterRequest) Convert() *User {
	return &User{
		Id:    uuid.NewString(),
		Name:  v.Name,
		Email: v.Email,
		Age:   v.Age,
		Role:  VOTER,
	}
}

// Admin Types

type AdminID = ID

type Admin struct {
	User
}

type AdminRequest struct {
	Name     string   `json:"name"`
	Age      int8     `json:"age"`
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

func (d *Admin) SetID() {
	d.Id = uuid.New().String()
}

// Token Types

type Tokens struct {
	Access  Token `json:"access"`
	Refresh Token `json:"refresh"`
}
