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
		// todo: uncomment order.Status line in production
		// if order.Status != "completed" {
		// 	continue
		// }
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
