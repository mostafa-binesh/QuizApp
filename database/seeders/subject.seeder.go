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
		"Gastrointestinal Health",
		"Musculoskeletal Health",
		"Neurological Health",
		"Endocrine Health",
		"Renal Health",
		"Dermatological Health",
		"Psychiatric Health",
		"Pediatric Health",
	}

	// Define an array of options for each question
	systems := [][]string{
		{"Pain", "Visual and Auditory Problems", "Cancer", "Reproductive Problems"},
		{"Immune", "Infectious Disease", "Reproductive"},
		// Gastrointestinal Health
		{
			"Digestive System Anatomy",
			"Gastrointestinal Disorders",
			"Gastrointestinal Diseases",
			"Gastrointestinal Infections",
		},

		// Musculoskeletal Health
		{
			"Musculoskeletal System Anatomy",
			"Musculoskeletal Disorders",
			"Musculoskeletal Diseases",
			"Musculoskeletal Injuries",
		},

		// Neurological Health
		{
			"Nervous System Anatomy",
			"Neurological Disorders",
			"Neurological Diseases",
			"Neurological Injuries",
		},

		// Endocrine Health
		{
			"Endocrine System Anatomy",
			"Endocrine Disorders",
			"Endocrine Diseases",
			"Endocrine Injuries",
		},

		// Renal Health
		{
			"Urinary System Anatomy",
			"Renal Disorders",
			"Renal Diseases",
			"Renal Injuries",
		},

		// Dermatological Health
		{
			"Skin Anatomy",
			"Dermatological Disorders",
			"Dermatological Diseases",
			"Dermatological Injuries",
		},

		// Psychiatric Health
		{
			"Mental Health Disorders",
			"Mental Health Diseases",
			"Mental Health Treatments",
			"Psychiatric Medications",
		},

		// Pediatric Health
		{
			"Child Development",
			"Pediatric Disorders",
			"Pediatric Diseases",
			"Pediatric Treatments",
		},
	}

	// Insert the questions and options into the database
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	var coursesCount int64
	D.DB().Model(&M.Course{}).Count(&coursesCount)
	for i, subject := range subjects {
		newSubject := M.Subject{Title: subject, CourseID: uint(rand.Intn(int(coursesCount)))}
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
