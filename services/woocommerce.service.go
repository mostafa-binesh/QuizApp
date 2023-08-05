package services

import (
	D "docker/database"
	M "docker/models"
	U "docker/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// import (
// 	U "docker/utils"
// 	"github.com/chenyangguang/woocommerce"
// )

type WCProduct struct {
	ID         int                   `json:"id"`
	Name       string                `json:"name"`
	Attributes []WCProductAttributes `json:"attributes"`
}
type WCProductAttributes struct {
	Name    string    `json:"name"`
	Options []*string `json:"options"`
}

// in wc. products model, we have serveral attributes:
// 1. parent courses (products) have only duration attr.
// 2. non-parent courses have ValidityPeriod and ParentCourse (parentCourseID)

func GetAllProducts() ([]WCProduct, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/products?per_page=100", U.Env("WC_BASE_URL"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(U.Env("WC_CONSUMER_KEY"), U.Env("WC_CONSUMER_SECRET"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve products. status code: %d", resp.StatusCode)
	}
	var products []WCProduct
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func GetProductByID(productID int) (WCProduct, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/products/%d", U.Env("WC_BASE_URL"), productID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return WCProduct{}, err
	}

	req.SetBasicAuth(U.Env("WC_CONSUMER_KEY"), U.Env("WC_CONSUMER_SECRET"))

	resp, err := client.Do(req)
	if err != nil {
		return WCProduct{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WCProduct{}, fmt.Errorf("failed to retrieve product. status code: %d", resp.StatusCode)
	}

	var product WCProduct
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return WCProduct{}, err
	}

	return product, nil
}
func ConvertWCCourseToCourseModel(wcCourses *[]WCProduct) (*[]M.Course, error) {
	courses := make([]M.Course, len(*wcCourses))
	for _, wcCourse := range *wcCourses {
		newCourse := M.Course{}
		for _, att := range wcCourse.Attributes {
			if len(att.Options) != 1 || att.Options[0] == nil {
				continue
			}
			// todo : change second and third conditions to else if
			// > because an attr. can only have one value of the checked ones
			// duration is only for parent course
			if att.Name == "Duration" {
				newCourse.Duration, _ = strconv.ParseUint(*att.Options[0], 10, 32)
			}
			// parentCourse and ValidityPeriod is only for children courses
			if att.Name == "ParentID" {
				parsedUint64, _ := strconv.ParseUint(*att.Options[0], 10, 64)
				parsedUint := uint(parsedUint64)
				newCourse.ParentID = &parsedUint
			}
			if att.Name == "ValidityPeriod" {
				parsedUint64, _ := strconv.ParseUint(*att.Options[0], 10, 64)
				newCourse.ValidityDaysPeriod = uint(parsedUint64)
			}
		}
		newCourse.WoocommerceID = uint(wcCourse.ID)
		newCourse.Title = wcCourse.Name
		courses = append(courses, newCourse)
	}
	return &courses, nil
}

// converts woocommerce time string to time.Time
func ConvertWCTimeToStandard(wcTime string) time.Time {
	// Layout defines the format of the input string
	layout := "2006-01-02T15:04:05"
	// Parse the input string and convert it to a time.Time object
	time, _ := time.Parse(layout, wcTime)
	return time

}

// gets all products from wc. api and save them into the database
func ImportCoursesFromWoocommerce() (*[]M.Course, error) {
	// get all woocommerce products from its service
	wcCourses, err := GetAllProducts()
	if err != nil {
		return nil, err
	}
	// convert the wooc. retreived products into course model
	convertedCourses, err := ConvertWCCourseToCourseModel(&wcCourses)
	if err != nil {
		return nil, err
	}
	// because some courses need to find reference to their parent in the database
	// and maybe it's the first time that the parent course is being inserted into the database
	// we need to insert all of courses into the database first and ignore all errors
	for _, course := range *convertedCourses {
		// for some reason, some of retreived courses were empty
		if course.Title == "" {
			continue
		}
		// saving the courses to the database
		// not using create func. to create courses to the db
		// becausew some of courses variables such as duration may been changed
		D.DB().Save(&course)
	}
	// create a map to store parent courses id and its courses
	DBCourseWoocommerceIDMap := make(map[uint]*M.Course)
	for _, course := range *convertedCourses {
		// for some reason, some of retreived courses were empty
		if course.Title == "" {
			continue
		}
		var existingCourse M.Course
		var parentCourse M.Course
		// Try to find the course with the given woocommerce_id
		wcCoruseResult := D.DB().Where("woocommerce_id = ?", course.WoocommerceID).First(&existingCourse)
		// also need to find the parent course with the given course.parentid
		// NOTE that course.ParentID comes from woocommerce !
		if course.ParentID != nil {
			// if the parent id wasn't submitted in the map before
			// try to find the parent dbCourse and submit the parent id record
			if DBCourseWoocommerceIDMap[*course.ParentID] == nil {
				// parent course map doesn't exist in the map
				// try to find it in the database
				// parentResult := D.DB().First(&parentCourse, *course.ParentID)
				parentResult := D.DB().Where("woocommerce_id = ?", course.ParentID).First(&parentCourse)
				// if the parent exist
				if parentResult.RowsAffected == 1 {
					DBCourseWoocommerceIDMap[*course.ParentID] = &parentCourse
				}
			} else {
				parentCourse = *DBCourseWoocommerceIDMap[*course.ParentID]
			}
		}

		// if course with desired woocommerce_id not found, we need to create one
		if wcCoruseResult.RowsAffected == 0 {
			// if parent exist, set parentID to dbCourse id
			if parentCourse.Title != "" {
				course.ParentID = &parentCourse.ID
			}
			course.Duration = parentCourse.Duration
			// insert it to the database
			D.DB().Create(&course)
		} else {
			// if found, just update it
			// course is convereted wc course
			existingCourse.Title = course.Title
			if parentCourse.Title != "" {
				existingCourse.ParentID = &parentCourse.ID
				existingCourse.Duration = parentCourse.Duration
			}
			existingCourse.ValidityDaysPeriod = course.ValidityDaysPeriod
			// update it to db
			D.DB().Save(&existingCourse)
		}
	}
	return convertedCourses, nil
}
