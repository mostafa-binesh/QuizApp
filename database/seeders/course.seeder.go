package seeders

import (
	S "docker/services"
)

func CourseSeeder() error {
	_, err := S.ImportCoursesFromWoocommerce()
	return err
}
