package services

import (
	D "docker/database"
	M "docker/models"
	U "docker/utils"
	"fmt"
	"time"
)

// returns bought courses using orderID
// returned courses are childCourses !
// second returned value is a map of [wcID]payTime
func ImportUserCoursesUsingOrderID(orderID uint) ([]M.Course, map[int64]time.Time, error) {
	// get user orders
	order, err := U.WCClient().Order.Get(int64(orderID), nil)
	if err != nil {
		return nil, nil, err
	}
	// todo: move this add user's bought courses to its service
	// add bought courses to user's courses
	var purchasedCoursesIDs []int64
	purchasedCourseIDPayDateMap := make(map[int64]time.Time)
	for _, item := range order.LineItems {
		// if order is not completed, don't process the order
		if order.Status != "completed" {
			continue
		}
		purchasedCoursesIDs = append(purchasedCoursesIDs, item.ProductID)
		var payTime time.Time
		// todo: dev only, delete order.DatePaidGMT=="" on production
		// if DatePaidGMT is empty, the order hasn't been paid
		if order.DatePaidGmt == "" {
			payTime = time.Now()
		} else {
			payTime = ConvertWCTimeToStandard(order.DatePaidGmt)
		}
		purchasedCourseIDPayDateMap[item.ProductID] = payTime
	}
	// get the courses
	var childrenCourses []M.Course
	if result := D.DB().Find(&childrenCourses, "woocommerce_id IN ?", purchasedCoursesIDs); result.Error != nil {
		return nil, nil, err
	}
	// get the parent course IDs using child courses
	var parentCourseIDs []uint
	for _, childrenCourse := range childrenCourses {
		if childrenCourse.ParentID != nil {
			parentCourseIDs = append(parentCourseIDs, *childrenCourse.ParentID)
		}
	}
	// if no the order doesn't contain any valid child course, show error
	if len(childrenCourses) == 0 || len(parentCourseIDs) == 0 {
		return nil, nil, fmt.Errorf("Your order doesn't contain any valid course")
	}
	return childrenCourses, purchasedCourseIDPayDateMap, nil
}

// DEPREACTED - NOT USED ANYWHERE
// converts courses to course_user
func AddCourseUserUsingCourses(childrenCourses []M.Course, purchasedCourseIDPayDateMap map[int64]time.Time, userID uint) []M.CourseUser {
	// create course_user records
	var userNewBoughtCourses []M.CourseUser
	for _, childCourse := range childrenCourses {
		// if the childCoruse doesn't have parent, don't insert it
		if childCourse.ParentID == nil {
			continue
		}
		userNewBoughtCourses = append(userNewBoughtCourses, M.CourseUser{
			UserID:         int(userID),
			CourseID:       int(*&childCourse.ID),
			ExpirationDate: purchasedCourseIDPayDateMap[int64(childCourse.WoocommerceID)].Add(time.Duration(childCourse.ValidityDaysPeriod) * time.Hour * 24),
		})
	}
	return userNewBoughtCourses
}

// returns courses which should be inserted and what courses needs update in course_update table based on orderID
func ExtractCourseToInsertAndToUpdate(childCourses []M.Course, purchasedCourseIDPayDateMap map[int64]time.Time, userID uint) (courseUsersToUpdate []M.CourseUser, newCourseUsers []M.CourseUser, err error) {
	// Create a map to store the courseUser data indexed by course ID
	courseUserMap := make(map[uint]*M.CourseUser)

	// Iterate over childCourses and update the courseUserMap with expiration dates
	for _, course := range childCourses {
		if cu, ok := purchasedCourseIDPayDateMap[int64(course.WoocommerceID)]; ok {
			courseUserMap[course.ID] = &M.CourseUser{
				UserID:         int(userID),
				CourseID:       int(course.ID),
				ExpirationDate: cu.AddDate(0, 0, int(course.ValidityDaysPeriod)),
			}
		}
	}

	// Create a list of courseUser records to update
	courseUsersToUpdate = make([]M.CourseUser, 0, len(courseUserMap))

	// Retrieve existing courseUser records from the database
	existingCourseUsers := make([]M.CourseUser, 0)
	if err := D.DB().Where("user_id = ?", userID).Find(&existingCourseUsers).Error; err != nil {
		return nil, nil, err
	}

	// Iterate over existing courseUser records and update expiration dates if needed
	for _, existing := range existingCourseUsers {
		if cu, ok := courseUserMap[uint(existing.CourseID)]; ok && existing.ExpirationDate != cu.ExpirationDate {
			existing.ExpirationDate = cu.ExpirationDate
			courseUsersToUpdate = append(courseUsersToUpdate, existing)
		}
		delete(courseUserMap, uint(existing.CourseID))
	}

	// Create a list of new courseUser records to insert
	newCourseUsers = make([]M.CourseUser, 0, len(courseUserMap))
	for _, cu := range courseUserMap {
		newCourseUsers = append(newCourseUsers, *cu)
	}
	return courseUsersToUpdate, newCourseUsers, nil
}
