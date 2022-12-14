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
		AllowHeaders: "Cache-Control,Content-Type",
	}))

	controller := Controller{
		s: &Service{
			repo: &Repository{
				db: db,
			},
			broadcaster: NewBroadcaster(),
		},
	}

	app.Get("/api/sse", controller.UpdateState)
	app.Post("/api/clear", controller.ClearRoom)
	app.Post("/api/player", controller.UpdatePlayer)
	app.Post("/api/room", controller.UpdateRoom)

	return app

}

func main() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&Users{})
	if err != nil {
		panic(err)
	}
	user := Users{ID: 1, RoomID: 1}
	db.Save(&user)
	app := InitApp(db)
	app.Listen(":3001")
}
