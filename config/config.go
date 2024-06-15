package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lohithk3345/voting_system/helpers"
)

var EnvMap = getENV()

func getENV() map[string]string {
	envMap := make(map[string]string)
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}
	envMap["REDIS_URL"] = os.Getenv("REDIS_URL")
	envMap["API_KEY"] = os.Getenv("API_KEY")
	envMap["TOKEN_SECRET"] = os.Getenv("TOKEN_SECRET")
	envMap["MONGODB_URL"] = os.Getenv("MONGODB_URL")
	envMap["REST_API_PORT"] = os.Getenv("REST_API_PORT")
	envMap["GRPC_PORT"] = os.Getenv("GRPC_PORT")
	for key, value := range envMap {
		if helpers.IsEmpty(value) {
			log.Fatalf("ENV not found for %s", key)
		}
	}
	return envMap
}
