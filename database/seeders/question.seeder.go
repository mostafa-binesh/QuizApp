package seeders

import (
	D "docker/database"
	M "docker/models"
	U "docker/utils"
	"fmt"
	"math/rand"
	"time"
)

// course seeder needed first
func QuestionAndOptionsSeeder() {
	// Define an array of questions
	questions := [20]string{
		"What is the capital of Spain?",
		"What is the largest planet in our solar system?",
		"What is the smallest country in the world?",
		"What is the largest ocean in the world?",
		"What is the highest mountain in the world?",
		"What is the smallest planet in our solar system?",
		"What is the largest country in the world by area?",
		"What is the fastest land animal?",
		"What is the tallest animal in the world?",
		"What is the longest river in the world?",
		"What is the largest desert in the world?",
		"What is the most populous country in the world?",
		"What is the hottest continent on Earth?",
		"What is the deepest ocean in the world?",
		"What is the largest bird in the world?",
		"What is the largest continent in the world?",
		"What is the largest mammal in the world?",
		"What is the largest fish in the world?",
		"What is the fastest sea animal?",
		"What is the largest reptile in the world?",
	}

	// Define an array of options for each question
	options := [20][4]string{
		{"Madrid", "Barcelona", "Seville", "Valencia"},
		{"Mars", "Venus", "Jupiter", "Saturn"},
		{"Vatican City", "Monaco", "Nauru", "Maldives"},
		{"Atlantic Ocean", "Indian Ocean", "Arctic Ocean", "Pacific Ocean"},
		{"Mount Everest", "K2", "Kangchenjunga", "Lhotse"},
		{"Mercury", "Venus", "Earth", "Mars"},
		{"Russia", "Canada", "China", "United States"},
		{"Cheetah", "Lion", "Gazelle", "Leopard"},
		{"Giraffe", "Elephant", "Hippopotamus", "Rhinoceros"},
		{"Amazon River", "Nile River", "Yangtze River", "Mississippi River"},
		{"Sahara Desert", "Arabian Desert", "Gobi Desert", "Antarctic Desert"},
		{"China", "India", "United States", "Indonesia"},
		{"Africa", "Asia", "Australia", "Antarctica"},
		{"Pacific Ocean", "Atlantic Ocean", "Indian Ocean", "Southern Ocean"},
		{"Ostrich", "Eagle", "Albatross", "Condor"},
		{"Asia", "Africa", "North America", "South America"},
		{"Blue Whale", "African Elephant", "Giraffe", "Hippopotamus"},
		{"Whale Shark", "Great White Shark", "Tiger Shark", "Hammerhead Shark"},
		{"Sailfish", "Tuna", "Marlin", "Shark"},
		{"Saltwater Crocodile", "Nile Crocodile", "American Alligator", "Gharial"},
	}

	// Insert the questions and options into the database
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	// get systems count for question refrence
	systems := []M.System{}
	D.DB().Preload("Subject.Course").Find(&systems)
	systemsCount := len(systems)
	fmt.Printf("systems count: %d\n", systemsCount)
	for i, questionText := range questions {
		systemID := uint(rand.Intn(int(systemsCount)))
		question := M.Question{
			// ! rand int may return 0 ! so i added 1 to solve this
			Title:    questionText,
			SystemID: systemID,
			Status:   "unvisited", Description: "Random description for " + questionText,
			CourseID: &systems[systemID].Subject.CourseID,
		}
		D.DB().Create(&question)
		optionsLen := len(options[i])
		correctOption := rand.Intn(optionsLen)
		for j, optionText := range options[i] {
			option := M.Option{
				Title:      optionText,
				Index:      string('A' + j),
				QuestionID: question.ID,
				IsCorrect:  U.ConvertBoolToUint(j == correctOption), // Set oneof the options as correct randomly // TODO changed to 1 for test of isCorrect uint test
			}
			D.DB().Create(&option)
		}
	}

	fmt.Println("Questions and options inserted successfully.")
}
