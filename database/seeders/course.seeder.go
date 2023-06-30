package seeders

import (
	D "docker/database"
	M "docker/models"
)

func CourseSeeder() {
	medicalCourses := []string{
        "Anatomy and Physiology",
        "Pharmacology",
        "Medical Terminology",
        "Pathophysiology",
        "Medical Ethics",
		"NC-LEX RN", // main course , id  = 5
    }
	for i, course := range medicalCourses {
		D.DB().Create(&M.Course{
			Title: course,
			WoocommerceID: uint(i+1),
		})
	}
}