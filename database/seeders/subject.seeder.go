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
	realSubjects := []string{
		"PDA in common Health Scenarios",
		"Adult Health",
		"Child Health",
		"Maternal and Newborn Health care",
		"Mental health",
	}
	realSystems := [][]string{
		{
			"Chapter 1. Pain",
			"Chapter 10. Visual and Auditory Problems",
			"Chapter 11. Musculoskeletal Problems",
			"Chapter 12. Gastrointestinal and Nutritional Problems",
			"Chapter 13. Diabetes Mellitus",
			"chapter 14. Other Endocrine Problems",
			"Chapter 15. Integumentary Problems",
			"Chapter 16. Renal and Urinary Problems",
			"Chapter 17. Reproductive Problems",
			"Chapter 2. Cancer",
			"Chapter 3. Immunologic Problems",
			"Chapter_4_Fluid,_Electrolyte,_and_Acid_Base_Balance_Problems",
			"Chapter 5. Safety and Infection Control",
			"Chapter 6. Respiratory Problems",
			"Chapter 7. Cardiovascular Problems",
			"Chapter 8. Hematologic Problems",
			"Chapter 9. Neurologic Problems",
			"Chapter 18. Problems in Pregnancy and Childbearing",
			"Chapter 19. Pediatric Problems",
			"Chapter 20. Pharmacology",
			"Chapter 21. Emergencies and Disasters",
			"Chapter 22. PsychiatricMental Health Problems",
		},
		{
			"Cardiovascular",
			"Endocrine",
			"GastrointestinalNutrition",
			"HematologicalOncological",
			"Immune",
			"Infectious Disease",
			"Integumentary",
			"Musculoskeletal",
			"Neurologic",
			"Reproductive",
			"respiratory",
			"UrinaryRenal",
			"VisualAuditory",
		},
		{
			"Cardiovascular",
			"Endocrine",
			"GastrointestinalNutrition",
			"Growth & Development",
			"Immune",
			"Infectious Disease",
			"Integumentary",
			"Musculoskeletal",
			"Neurologic",
			"Respiratory",
			"UrinaryRenal",
			"VisualAuditory",
		},
		{
			"Antepartum",
			"LaborDelivery",
			"Newborn",
			"PostPartum",
		},
		{"Mental Health Concepts"},
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
	// add real subject and systems
	NCLEXCourse := &M.Course{}
	D.DB().Where("Title = ?", "NC-LEX RN").Find(&NCLEXCourse)
	for i, subject := range realSubjects {
		newSubject := &M.Subject{
			Title:    subject,
			CourseID: NCLEXCourse.ID,
		}
		D.DB().Create(&newSubject)
		for _, system := range realSystems[i] {
			newSystem := &M.System{
				Title:     system,
				SubjectID: newSubject.ID,
			}
			D.DB().Create(&newSystem)
		}
	}

	fmt.Println("Subjects and systems inserted successfully.")
}
