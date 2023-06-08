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
	// ! admin Route
	admin := router.Group("/admin")
	admin.Get("/users", AC.IndexUser)
	admin.Get("/users/:id<int>", AC.UserByID)
	admin.Put("/users/:id<int>", AC.EditUser)
	admin.Post("/users", AC.AddUser)
	admin.Delete("/users/:id<int>", AC.DeleteUser)
	admin.Get("/laws", AC.IndexLaw)
	admin.Get("/laws/search", AC.LawSearch)
	admin.Get("laws/:id<int>", C.LawByID)
	admin.Post("/laws", AC.CreateLaw)
	admin.Put("/laws/:id<int>", AC.UpdateLaw)
	admin.Delete("/laws/:id<int>", AC.DeleteLaw)
	admin.Delete("/laws/:id<int>/files/:fileID<int>", AC.DeleteFile) // ! TODO : file az storage ham bayad paak she
	// ! authentication routes
	router.Post("/signup", C.SignUpUser)
	router.Post("/login", C.Login)
	router.Get("/logout", C.Logout)
	// ! messaging
	msg := router.Group("correspondence")
	msg.Use(encryptcookie.New(encryptcookie.Config{
		// ! only base64 charasters
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