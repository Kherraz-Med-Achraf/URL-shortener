package main

import "github.com/gofiber/fiber/v2"

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello World en Go avec Fiber ! 🚀")
    })

    app.Listen(":3000")
}