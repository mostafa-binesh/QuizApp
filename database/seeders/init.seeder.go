package seeders

import "fmt"

func InitSeeder() {
	fmt.Println("user seeder")
	UserSeeder()
	fmt.Println("law seeder")
	// LawSeeder()
	fmt.Println("law comment seeder")
	// LawCommentsSeeder()
	fmt.Println("admin seeder")
	AdminSeeder()
	fmt.Println("course seeder")
	CourseSeeder()
	fmt.Println("question seeder")
	QuestionSeeder()
}
