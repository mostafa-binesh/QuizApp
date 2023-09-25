package controllers

import (
	D "docker/database"
	M "docker/models"
	S "docker/services"
	U "docker/utils"
	"fmt"

	// ut "github.com/go-playground/universal-translator"
	// "github.com/go-playground/validator/v10"
	"os"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func TranslationTest(c *fiber.Ctx) error {
	type User struct {
		Username string `validate:"required" json:"username"`
		Password string `validate:"required" json:"wirdpass"`
	}

	user := User{Username: "kurosh79"}
	if errs := U.Validate(user); errs != nil {
		return c.JSON(fiber.Map{"errors": U.Validate(user)})
	}
	return c.JSON(fiber.Map{"msg": "everything is fine"})
}

func DevAllUsers(c *fiber.Ctx) error {
	users := []M.User{}
	pagination := new(U.Pagination)
	if err := c.QueryParser(pagination); err != nil {
		U.ResErr(c, err.Error())
	}
	D.DB().Scopes(U.Paginate(users, pagination)).Find(&users)
	return c.JSON(fiber.Map{
		"meta": pagination,
		"data": users,
	})
}
func UploadFile(c *fiber.Ctx) error {
	type Upload struct {
		FirstName string `json:"firstName" validate:"required"`
		LastName  string `json:"lastName" validate:"required"`
		File      string `json:"file" validate:"required"`
	}
	payload := new(Upload)
	if err := c.BodyParser(payload); err != nil {
		return c.JSON(fiber.Map{
			"error": err,
		})
	}
	file, err := c.FormFile("file")
	// if err != nil {
	// ! if file not exists, we get error: there is no uploaded file associated with the given key
	// 	return c.JSON(fiber.Map{"error": err.Error()})
	// }
	if file != nil {
		payload.File = file.Filename
	}
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// check if file with this name already exists
	if U.FileExistenceCheck(file.Filename, U.UploadLocation) {
		return U.ResErr(c, "file already exists")
	}
	// ! file extension check
	if !(U.HasImageSuffixCheck(file.Filename) || U.HasSuffixCheck(file.Filename, []string{"pdf"})) {
		return c.SendString("file should be image or pdf! please fix it")
	}
	// Save file to disk
	err = c.SaveFile(file, fmt.Sprintf(U.UploadLocation+"/%s", file.Filename))
	if err != nil {
		return U.ResErr(c, "cannot save | "+err.Error())
	}
	return c.JSON(fiber.Map{"msg": "فایل آپلود شد"})
}

func ExistenceCheck(c *fiber.Ctx) error {
	filename := c.FormValue("fileName")
	directory := c.FormValue("dir")
	if _, err := os.Stat(directory + "/" + filename); os.IsNotExist(err) {
		return c.SendString("File does not exist")
	} else {
		return c.SendString("File exists")
	}
}
func GormG(c *fiber.Ctx) error {
	type pashm struct {
		Name         string `json:"name" validate:"required,dunique=users.name"` // users table, name column
		PersonalCode string `json:"personalCode" validate:"required,dexists=users"`
	}
	payload := new(pashm)
	// parse payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// ! if you're in edit and wanna ignore the user's information rows
	// ! - you need to pass the id to validation function as well
	// ! -- eg. the user's phoneNumber is 1234 and you've used dunique in phoneNumber field
	// ! --- but if you check the user's row, you'll get the user's phoneNumber and unique validation will fail
	// ! ---- but you don't want this. so you need to ignore that specific id
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	return c.SendString("no error")
}

// messaing
type Guest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func MessageRegister(c *fiber.Ctx) error {
	payload := new(Guest)
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	sess := U.Session(c)
	sess.Set("type", "guest")
	sess.Set("guest", payload)
	if err := sess.Save(); err != nil {
		return U.ResErr(c, err.Error())
	}
	return U.ResMsg(c, "ورود انجام شد")
}
func SeeMessages(c *fiber.Ctx) error {
	sess := U.Session(c)
	// sess, err := U.Store.Get(c)
	// if err != nil {
	// 	panic("cannot get the session")
	// }

	// return c.SendString(fmt.Sprintf("%s", sess.Get("type")))
	// userID := sess.Get(U.USER_TYPE)
	// if userID == nil {
	// 	return ReturnError(c, "not authenticated", fiber.StatusUnauthorized) // ! notAuthorized is notAuthenticated
	// } else {
	// 	return c.SendString(fmt.Sprintf("user id is: %s", sess.Get(U.USER_ID)))
	// }
	// session := U.Session(c)
	if sess != nil {
		return U.ResMsg(c, fmt.Sprintf("type: %v, guest: %v", sess.Get("type"), sess.Get("guest")))
	}
	// return U.ResMsg(c, "session KHALIE")
	// if sess.Get("type") == "guest" {
	// 	guest := sess.Get("guest").(Guest)
	// 	return c.JSON(fiber.Map{
	// 		"data":      guest,
	// 		"sessionID": sess.ID(),
	// 	})
	// }
	return U.ResErr(c, "شما باید وارد شوید")
}

//	func FiberContextMemoryAddress(c *fiber.Ctx) error {
//		fmt.Printf("utility memory: %p\n, function context memory ad.: %p\n", U.FiberCtx(), c)
//		return c.SendString("ss")
//	}
func StructInfo(c *fiber.Ctx) error {
	type Post struct {
		PostName string
	}
	type User struct {
		Name     string `gorm:"varchar(255)"`
		LastName string
		Body     string `gorm:"type:text"`
		Posts    Post
	}

	u := User{}
	s := reflect.TypeOf(u)

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		fmt.Printf("%s %s %s\n", field.Name, field.Type, field.Tag.Get("gorm"))
	}
	return c.SendString("sdassds")
}
func ResetMemory(c *fiber.Ctx) error {
	U.Memory.Reset()
	return c.SendString("done")
}
func WCProducts(c *fiber.Ctx) error {
	wcProducts, err := S.GetAllProducts()
	if err != nil {
		return U.ResErr(c, err.Error())
	}
	return c.JSON(fiber.Map{"data": wcProducts})
}
func AllQuestions(c *fiber.Ctx) error {
	var questions []M.Question
	if err := D.DB().Find(&questions).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": questions})
}

// get all answers, checks if the answer is correct or not
func AnswerCorrection(c *fiber.Ctx) error {
	var answers []M.UserAnswer
	if err := D.DB().
		Preload("Question.Options").
		Find(&answers).Error; err != nil {
		return err
	}
	for i := range answers {
		answers[i].IsCorrect = answers[i].IsChosenOptionsCorrect()
		// set the question to null because we're saving the options later in the handler
		// and don't want to save the questions and options again
		answers[i].Question = nil
	}

	if err := D.DB().Save(&answers).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": answers})
}

// fix all question course_id
func QuestionCorrect(c *fiber.Ctx) error {
	var questions []M.Question
	if err := D.DB().
		Preload("System.Subject").
		Find(&questions).Error; err != nil {
		return err
	}
	for i := range questions {
		questions[i].CourseID = &questions[i].System.Subject.CourseID
	}

	if err := D.DB().Save(&questions).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": questions})
}
func MigrateNewSubjects(c *fiber.Ctx) error {
	realSubjects := []string{
		"Critical Care",
		"Fundamentals",
		"Leadership and Management",
		"Pharmacology",
		"Prometric",
		"Saunders",
		"Self Assessments",
	}
	realSystems := [][]string{
		{
			"Critical Care Concepts",
		},
		{
			"Basic care and Comfort",
			"Medication Administration",
			"SafetyInfection Control",
			"Skills Procedures",
		},
		{
			"Management Concepts",
			"Assignment Delegation",
			"Ethical Legal",
			"Prioritization",
		},
		{
			"Analgesics",
			"Cardiovascular",
			"Endorcine",
			"Gastrointestinal Nutrition",
			"Hematological Oncological",
			"Immune",
			"Infection Decease",
			"Integumentary",
			"Musculoskeletal",
			"Neurologic",
			"Psychiatric Medications",
			"Reproductive Maternity Newborn",
			"Respiratory",
			"Urinary Renal",
			"Visual Auditory",
		},
		{"Prometric"},
		{
			"Chapter 8",
			"Chapter 9",
			"Chapter 10",
			"Chapter 11",
			"Chapter 12",
			"Chapter 13",
			"Chapter 14",
			"Chapter 15",
			"Chapter 16",
			"Chapter 17",
			"Chapter 18",
			"Chapter 19",
			"Chapter 20",
			"Chapter 21",
			"Chapter 22",
			"Chapter 23",
			"Chapter 24",
			"Chapter 25",
			"Chapter 26",
			"Chapter 27",
			"Chapter 28",
			"Chapter 29",
			"Chapter 30",
			"Chapter 31",
			"Chapter 32",
			"Chapter 33",
			"Chapter 34",
			"Chapter 35",
			"Chapter 36",
			"Chapter 37",
			"Chapter 38",
			"Chapter 39",
			"Chapter 40",
			"Chapter 41",
			"Chapter 42",
			"Chapter 43",
			"Chapter 44",
			"Chapter 45",
			"Chapter 46",
			"Chapter 47",
			"Chapter 48",
			"Chapter 48",
			"Chapter 49",
			"Chapter 50",
			"Chapter 50",
			"Chapter 51",
			"Chapter 52",
			"Chapter 53",
			"Chapter 54",
			"Chapter 55",
			"Chapter 56",
			"Chapter 57",
			"Chapter 58",
			"Chapter 59",
			"Chapter 60",
			"Chapter 61",
			"Chapter 62",
			"Chapter 63",
			"Chapter 64",
			"Chapter 65",
			"Chapter 66",
			"Chapter 67",
			"Chapter 68",
			"Chapter 69",
			"Chapter 70",
		},
		{
			"Self Assessments 1",
			"Self Assessments 2",
		},
	}
	// add real subject and systems
	NCLEXCourse := &M.Course{}
	if result := D.DB().Where("Title = ?", "NCLEX-RN مادر").Find(&NCLEXCourse); result.RowsAffected == 0 {
		return U.ResErr(c, "No course found with title of NC-LEX RN")
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
	return c.JSON(fiber.Map{"msg": "done"})
}
