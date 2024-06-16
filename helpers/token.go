package helpers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lohithk3345/voting_system/types"
	"go.mongodb.org/mongo-driver/bson"
)

type Token struct{}

var Tokens = Token{}

func (Token) ExtractUUIDFromInsertedID(insertedID interface{}) (*types.ID, error) {
	log.Println("START")
	bsonBytes, err := bson.Marshal(insertedID)
	if err != nil {
		log.Println("Error")
		return nil, fmt.Errorf("failed to marshal inserted ID to BSON: %w", err)
	}

	log.Println(bsonBytes)

	var insertedData types.ID
	err = bson.Unmarshal(bsonBytes, &insertedData)
	log.Println(insertedData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal BSON to MyData: %w", err)
	}

	return &insertedData, nil
}

func (h Token) GetToken(ctx *gin.Context) (types.Token, error) {
	token := ctx.GetHeader("Authorization")
	// ctx.header
	bearerFix := "Bearer "
	if token == types.EmptyString || len(token) < len(bearerFix) || token[:len(bearerFix)] != bearerFix {
		return types.EmptyString, gin.Error{}
	}
	actualToken := token[len(bearerFix):]
	return actualToken, nil
}
