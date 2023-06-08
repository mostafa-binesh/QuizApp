package seeders

import (
	D "docker/database"
	M "docker/models"
	"math/rand"
	"strconv"
	"time"
)

func AdminSeeder() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator with current time
	D.DB().Create(&M.User{
		Name:         "محمدمهدی کاظمی",
		PhoneNumber:  "09121318520",
		Password:     "$2a$10$f2vsTfJbqox.my6oPJ0eZeuEuhSVBqO3BUj4EExEE.UIqhhfEOwoG", // = password
		Role:         2,
		NationalCode: strconv.Itoa(rand.Intn(9000000000) + 1000000000), // Generate 10-digit number
		PersonalCode: "1234567890",                                     // Generate 10-digit number
		Verified:     true,
	})
}
