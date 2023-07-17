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
			SystemID: systems[systemID].ID,
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
				QuestionID: &question.ID,
				IsCorrect:  U.ConvertBoolToUint(j == correctOption), // Set oneof the options as correct randomly // TODO changed to 1 for test of isCorrect uint test
			}
			D.DB().Create(&option)
		}
	}
	// next gen questions seeder
	nextGenQuestionTypes := []M.QuestionType{M.NextGenerationSingleSelect, M.NextGenerationMultipleSelect, M.NextGenerationTableSingleSelect,
		M.NextGenerationTableMultipleSelect, M.NextGenerationTableDropDown}
	nextGenTexts := []string{"Next Gen Single", "Next Gen Multiple", "Next Gen Table Single", "Next Gen Table Multiple", "Next Gen Dropdown"}
	stringCounter := []string{"First of ", "Second of ", "Third of ", "Fourth of ", "Fifth of "}
	for i, gen := range nextGenQuestionTypes {
		dropDowns := []M.Dropdown{}
		systemID := uint(rand.Intn(int(systemsCount)))
		tabs := []M.Tab{M.Tab{
			Tables: []M.Table{M.Table{
				Title: "A New Table",
				Rows:  [][]string{{"row1 col1", "row1 col2", "row1 col3"}, {"row2 col1", "row2 col2", "row2 col1"}},
			}},
		}}
		if gen == M.NextGenerationTableDropDown {
			dropDowns = append(dropDowns, M.Dropdown{})
			dropDowns = append(dropDowns, M.Dropdown{})
		}
		question := M.Question{
			Title:       stringCounter[i] + nextGenTexts[i],
			Description: "A Random Description",
			// Images:      []M.Image{},
			SystemID:  systems[systemID].ID,
			CourseID:  &systems[systemID].Subject.CourseID,
			Type:      gen,
			Tabs:      tabs,
			Dropdowns: dropDowns,
		}
		D.DB().Create(&question)
		optionsLen := len(options[i])
		correctOption := rand.Intn(optionsLen)
		for j, optionText := range options[i] {
			var dropDownID *uint
			var questionID *uint
			questionID = &question.ID
			if len(dropDowns) > 0 {
				dropDownsCount := len(dropDowns)
				selectedDropdownID := rand.Intn(dropDownsCount)
				dropDownID = &dropDowns[selectedDropdownID].ID
				questionID = nil
			}
			option := M.Option{
				Title:      optionText,
				Index:      string('A' + j),
				QuestionID: questionID,
				IsCorrect:  U.ConvertBoolToUint(j == correctOption), // Set oneof the options as correct randomly // TODO changed to 1 for test of isCorrect uint test
				DropdownID: dropDownID,
			}
			D.DB().Create(&option)
		}
	}
	fmt.Println("Questions and options inserted successfully.")
}
