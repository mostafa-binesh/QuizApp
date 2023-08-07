package utils

import (
	"fmt"
	"os"

	env "github.com/joho/godotenv"
)

// returns envoirment variable if exist in the .env file
// otherwise, try to get the env. variable from host
// if there was no env. variable, panics
func Env(key string) string {
	// load .env file
	err := env.Load(".env")
	if err != nil {
		value, ok := os.LookupEnv(key)
		if !ok {
			panic(fmt.Sprintf("can't get env variable: %s", key))
		}
		return value
	}
	return os.Getenv(key)
}

// returns true if key variable in .env file is "true"
func EnvBool(key string) bool {
	return Env(key) == "true"
}
