package seeders

import (
	S "docker/services"
)

func CourseSeeder() error {
	_, err := S.ImportCoursesFromWoocommerce()
	if err != nil {
		return err
	}
	return nil
}
