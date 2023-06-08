package database

import (
	"fmt"
	"log"
	"os"
	"time"

	// "docker/config"
	// U "docker/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormDatabase *gorm.DB

func DB() *gorm.DB {
	if gormDatabase != nil {
		return gormDatabase
	}
	panic("CANNOT CONNECT TO THE DATABASE")
}

// func ConnectToDB() {
func ConnectToDB(DB_HOST string, DB_USERNAME string, DB_PASSWORD string, DB_NAME string, DB_PORT string) {
	var err error
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DB_USERNAME, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	gormDatabase, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database.\nConnection string: %s", dsn))
	}
	fmt.Println("database connection stablished")
}
func RowsCount(query string, searchValue string, ignoreID ...uint64) int {
	rows, err := gormDatabase.Raw(query, searchValue).Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	count := 0
	var id uint64
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			if len(ignoreID) > 0 {
				if id != ignoreID[0] {
					count++
					break
				}
			}
		}
	}
	return count
}
