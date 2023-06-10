package seeders

import (
	D "docker/database"
	M "docker/models"
	"fmt"
	"math/rand"
	"time"
)
// course seeder needed first 
func QuestionSeeder() {
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
	for i, questionText := range questions {
		question := M.Question{CourseID: 1, Title: questionText}
		D.DB().Create(&question)
		for j, optionText := range options[i] {
			option := M.Option{
				Title:      optionText,
				Index:      string('A' + j),
				QuestionID: question.ID,
				IsCorrect:  j == rand.Intn(4), // Set oneof the options as correct randomly
			}
			D.DB().Create(&option)
		}
	}

	fmt.Println("Questions and options inserted successfully.")
}
