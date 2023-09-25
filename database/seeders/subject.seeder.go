package seeders

import (
	D "docker/database"
	M "docker/models"
	"fmt"
)

func SubjectAndSystemSeeder() {
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
	// add real subject and systems
	NCLEXCourse := &M.Course{}
	if courseResult := D.DB().Where("Title = ?", "NCLEX-RN مادر").Find(&NCLEXCourse); courseResult.Error != nil {
		panic("cannot seed subjects because there was an error on getting course with title NCLEX-RN مادر")
	}
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
