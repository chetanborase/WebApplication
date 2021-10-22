package env

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func Get(key string, defaultValue string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	if strings.Contains(defaultValue, "nil") {
		panic(errors.New("Environment variable : " + key + " not found"))
	}
	return defaultValue
}
func LoadProfile() {
	if strings.EqualFold(Get("ACTIVE_PROFILE", "Dev"), "Dev") {
		if err := godotenv.Load(os.ExpandEnv(".env")); err != nil {
			log.Fatal("error reading .env ", err)
		}
	}
	log.Println("Profile Loaded")
}