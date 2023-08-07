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
	value, ok := os.LookupEnv(key)
	if !ok {
		err := env.Load(".env")
		if err != nil {
			panic(fmt.Sprintf("can't get env variable: %s", key))
		}
		return os.Getenv(key)
	}
	return value
}

// returns true if key variable in .env file is "true"
func EnvBool(key string) bool {
	return Env(key) == "true"
}
