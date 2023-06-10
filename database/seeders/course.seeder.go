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
    }
	for _, course := range medicalCourses {
		D.DB().Create(&M.Course{
			Title: course,
		})
	}
}