package routes

import (
	C "docker/controllers"
	AC "docker/controllers/admin"

	// "github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func APIInit(router *fiber.App) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"msg": "freeman was here <3",
		})
	})
	// ! quiz routes
	quiz := router.Group("/quiz")
	quiz.Post("/", C.CreateQuiz)
	// ! course routes
	user := router.Group("/user", C.AuthMiddleware)
	userCourse := user.Group("/courses")
	userCourse.Get("/", C.AllCourses)
	userCourse.Post("/", C.CreateQuiz)
	userCourse.Put("/update", C.UpdateUserCourses)
	userCourse.Get("/:courseID<int>/subjects", C.CourseSubjects)

	userQuiz := user.Group("/quizzes")
	userQuiz.Get("/", C.AllQuizzes)
	userQuiz.Get("/:id<int>", C.QuizByID)
	userQuiz.Put("/:id<int>", C.UpdateQuiz)
	userQuiz.Get("/:id<int>/report", C.QuizReport)
	userQuiz.Post("/", C.CreateQuiz)
	userQuiz.Post("/createFakeQuiz", C.CreateFakeQuiz)
	// userQuiz.Post("/createFakeNextGenerationQuiz", C.CreateFakeNextGenerationQuiz)
	userQuiz.Get("/overall", C.OverallReport)
	userQuiz.Get("/report", C.AllQuizzesReport)

	userNotes := user.Group("/notes")
	userNotes.Get("/", C.AllNotes)
	userNotes.Put("/:id<int>", C.EditNote)

	userQuestions := user.Group("/questions")
	userQuestions.Get("/search", C.AllQuestionsWithSearch)

	userStudyPlanner := user.Group("/studyPlanner")
	userStudyPlanner.Get("/", C.AllStudyPlans)
	userStudyPlanner.Post("/", C.CreateStudyPlanner)
	userStudyPlanner.Delete("/", C.DeleteStudyPlan)
	userStudyPlanner.Put("/finish", C.FinishDate)

	// ! admin routes
	// admin := router.Group("/admin")
	admin := router.Group("/admin", C.AuthMiddleware, C.RoleCheck([]string{"admin"})) // todo: replace main admin router with this
	admin.Get("/courses", AC.AllCourses)
	admin.Get("/nonParentCourses", AC.NonParentCourses)
	admin.Get("/courses/:id<int>", AC.CourseByID)
	admin.Post("/courses", AC.CreateCourse)
	admin.Get("/courses/subjects", AC.AllSubjects)
	admin.Get("/courses/:courseID<int>/subjects", AC.CourseSubjects)
	admin.Get("/courses/addFromWoocommerce", AC.ImportCoursesFromWooCommerce)

	admin.Get("/users", AC.IndexUser)
	admin.Get("/users/:email<string>", AC.UserByEmail)
	admin.Post("/users", AC.AddUser)
	admin.Put("/users/:id<int>", AC.EditUser)
	admin.Delete("/users/:id<int>", AC.DeleteUser)

	admin.Get("/questions", AC.AllQuestions)
	admin.Post("/questions/singleSelect", AC.CreateSingleSelectQuestion)
	admin.Put("/questions/singleSelect/:questionID<int>", AC.EditSingleSelectQuestion)
	admin.Post("/questions/multipleSelect", AC.CreateMultipleSelectQuestion)
	admin.Put("/questions/multipleSelect/:questionID<int>", AC.EditMultipleSelectQuestion)
	admin.Post("/questions/nextGeneration", AC.CreateNextGenerationQuestion)
	// admin.Put("/questions/nextGeneration/:questionID<int>", AC.EditNextGenerationQuestion)
	admin.Get("/questions/:id<int>", AC.QuestionByID)
	admin.Post("/uploadImages", AC.UploadImage)

	admin.Get("/systems/:systemID<int>/questions", AC.SystemQuestions)
	admin.Delete("/systems/:systemID<int>/questions", AC.DeleteSystemQuestions)

	admin.Get("/changeImageURLsInDescription", AC.ChangeImageURLsInDescription)
	// ! authentication routes
	auth := router.Group("/auth")
	auth.Post("/signup", C.SignUpUser)
	auth.Post("/login", C.Login)
	auth.Get("/logout", C.Logout)
	// ! devs route
	dev := router.Group("/devs")
	dev.Get("/translation", C.TranslationTest)
	dev.Get("/allUsers", C.DevAllUsers) // ?: send limit and page in the query
	dev.Get("/panic", func(c *fiber.Ctx) error { panic("PANIC!") })
	dev.Post("/upload", C.UploadFile)
	dev.Post("/fileExistenaceCheck", C.ExistenceCheck)
	dev.Post("/gormUnique", C.GormG)
	dev.Get("/resetMemory", C.ResetMemory)
	dev.Get("/wcProducts", C.WCProducts)
	dev.Get("/allQuestions", C.AllQuestions)
	// practical dev routes
	dev.Get("/autoMigrate", C.AutoMigrate)
	dev.Get("/answerCorrection", C.AnswerCorrection)
	dev.Get("/questionCorrect", C.QuestionCorrect)
	dev.Get("/migrateNewSubjects", C.MigrateNewSubjects)
	devPanel := dev.Group("/admin")
	devPanel.Get("/structInfo", C.StructInfo)
}
