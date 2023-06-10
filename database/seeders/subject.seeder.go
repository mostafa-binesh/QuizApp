package seeders

import (
	D "docker/database"
	M "docker/models"
	"fmt"
	"math/rand"
	"time"
)

func SubjectAndSystemSeeder() {
	// Define an array of questions
	subjects := []string{
		"PDA in common Health Scenarios",
		"Adult Health",
	}

	// Define an array of options for each question
	systems := [][]string{
		{"Pain", "Visual and Auditory Problems", "Cancer", "Reproductive Problems"},
		{"Immune", "Infectious Disease", "Reproductive"},
	}

	// Insert the questions and options into the database
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	for i, subject := range subjects {
		newSubject := M.Subject{Title: subject}
		D.DB().Create(&newSubject)
		for _, systemText := range systems[i] {
			option := M.System{
				Title:     systemText,
				SubjectID: newSubject.ID,
			}
			D.DB().Create(&option)
		}
	}

	fmt.Println("Subjects and systems inserted successfully.")
}
