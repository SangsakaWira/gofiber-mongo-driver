package main

import (
    "github.com/gofiber/fiber/v2"
    "rest-api/route"
)

func main() {
    app := fiber.New()
    route.RouteInit(app)
    app.Listen(":3000")
}
