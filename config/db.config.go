package config

import (
// U "docker/utils"
)

// var dbConfig *U.DBConfig

// func Init() error {
// 	config, err := U.GetDBConfig()

// 	if err != nil {
// 		return err
// 	}

// 	dbConfig = config

// 	// Initialize database connection here

// 	return nil
// }

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

// func GetDBConfig() *DBConfig {
// 	dbConfig := &DBConfig{
// 		Host:     U.Env("DB_HOST"),
// 		Port:     U.Env("DB_PORT"),
// 		Username: U.Env("DB_USERNAME"),
// 		Password: U.Env("DB_PASSWORD"),
// 		Name:     U.Env("DB_NAME"),
// 	}
// 	return dbConfig
// }
