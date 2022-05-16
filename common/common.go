package common

import (
	"encoding/json"
	"log"
	"os"
	"smbteam4/models"

	"github.com/joho/godotenv"
)

func GetenvData(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	value := os.Getenv(key)
	return value
}
func JSON_Marshell(res models.Response) string {
	resJson, _ := json.Marshal(res)
	return string(resJson)
}
