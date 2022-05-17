package helpers

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/renu-ramesh/robot-apocalypse-docker/models"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func GetenvData(key string) string {
	value := os.Getenv(key)
	return value
}
func JSON_Marshell(res models.Response) string {
	resJson, _ := json.Marshal(res)
	return string(resJson)
}
