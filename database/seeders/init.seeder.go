package seeders

import "fmt"

func InitSeeder() {
	fmt.Println("user seeder")
	UserSeeder()
	fmt.Println("admin seeder")
	AdminSeeder()
	fmt.Println("course seeder")
	if CourseSeeder() != nil {
		panic("error on seeding courses from woocommerce code")
	}
	fmt.Println("subject and system seeder")
	SubjectAndSystemSeeder()
	fmt.Println("question seeder")
}
