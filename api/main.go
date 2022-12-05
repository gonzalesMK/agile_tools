package main

import (
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func InitApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	// CORS for external resources
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Cache-Control",
	}))

	controller := Controller{
		s: &Service{
			repo: &Repository{
				db: db,
			},
			broadcaster: NewBroadcaster(),
		},
	}

	app.Get("/sse", controller.UpdateState)

	return app

}

func main() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&Users{})
	if err != nil {
		panic(err)
	}
	app := InitApp(db)
	app.Listen(":3000")
}
