package seeders

import (
	D "docker/database"
	M "docker/models"
	"math/rand"
	// "strconv"
	"time"
)

func AdminSeeder() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator with current time
	D.DB().Create(&M.User{
		Email:    "admin@admin.com",
		Password: "$2a$10$f2vsTfJbqox.my6oPJ0eZeuEuhSVBqO3BUj4EExEE.UIqhhfEOwoG", // = password
		Role:     2,
		Verified: true,
	})
}
