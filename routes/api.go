package routes

import (
	C "docker/controllers"
	AC "docker/controllers/admin"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
)

func APIInit(router *fiber.App) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"msg":        "freeman was here :)",
			"lastChange": "add static files",
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
	userCourse.Get("/update", C.UpdateUserCourses)
	userCourse.Get("/:courseID<int>/subjects", C.AllSubjects)

	userQuiz := user.Group("/quizzes")
	userQuiz.Get("/", C.AllQuizzes)
	userQuiz.Get("/:id<int>", C.QuizByID)
	userQuiz.Put("/:id<int>", C.UpdateQuiz) // TODO not tested yet
	userQuiz.Post("/", C.CreateQuiz)
	userQuiz.Post("/createFakeQuiz", C.CreateFakeQuiz)
	userNotes := user.Group("/notes")
	userNotes.Get("/", C.AllNotes)
	userNotes.Put("/:id<int>", C.EditNote)
	// ! admin routes
	admin := router.Group("/admin")
	admin.Get("/courses", AC.AllCourses)
	admin.Get("/courses/:id<int>", AC.CourseByID)
	admin.Post("/courses", AC.CreateCourse)
	admin.Get("/courses/:courseID<int>/subjects", AC.AllSubjects)
	admin.Get("/courses/addFromWoocommerce", AC.AddCoursesFromWooCommerce)

	admin.Get("/users/:email<string>", AC.UserByEmail)
	admin.Post("/users", AC.AddUser)
	admin.Put("/users/:id<int>", AC.EditUser)
	admin.Delete("/users/:id<int>", AC.DeleteUser)
	admin.Post("/questions", AC.AddQuestion)
	admin.Get("/questions/:id<int>", AC.QuestionByID)
	admin.Post("/uploadImages", AC.UploadImage)

	admin.Get("/users", AC.IndexUser)
	// ! authentication routes
	auth := router.Group("/auth")
	auth.Post("/signup", C.SignUpUser)
	auth.Post("/login", C.Login)
	auth.Get("/logout", C.Logout)
	// ! messaging
	msg := router.Group("correspondence")
	msg.Use(encryptcookie.New(encryptcookie.Config{
		// ! only base64 characters
		// ! A-Z | a-z | 0-9 | + | /
		Key: "S6e5+xc65+4dfs/nb4/f56+EW+56N4d6",
	}))
	// ! dashboard routes
	dashboard := router.Group("/dashboard", C.AuthMiddleware)
	dashboard.Get("/", C.Dashboard)
	// ! devs route
	dev := router.Group("/devs")
	dev.Get("/autoMigrate", C.AutoMigrate)
	dev.Get("/translation", C.TranslationTest)
	dev.Get("/pagination", C.PaginationTest) // ?: send limit and page in the query
	dev.Get("/allUsers", C.DevAllUsers)      // ?: send limit and page in the query
	dev.Get("/panic", func(c *fiber.Ctx) error { panic("PANIC!") })
	dev.Post("/upload", C.UploadFile)
	dev.Post("/fileExistenaceCheck", C.ExistenceCheck)
	dev.Post("/gormUnique", C.GormG)
	router.Get("/contextMemoryAddress", C.FiberContextMemoryAddress)
	devPanel := dev.Group("/admin")
	devPanel.Get("/structInfo", C.StructInfo)
}
