package handlers

import "github.com/gofiber/fiber/v2"

func StaticHandler(app *fiber.App) *fiber.App {
	app.Static("/", "./frontend/dist", fiber.Static{
		Compress: true,
	})
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile("./frontend/dist/index.html") })
	return app
}
