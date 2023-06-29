package route

import (
	"github.com/gofiber/fiber/v2"
	"rest-api/controller"
)

func RouteInit(route *fiber.App) {

	// Books
	route.Get("/book", controller.FindBook)
	route.Post("/book", controller.CreateBook)
	route.Get("/book/:id", controller.FindBookById)
	route.Put("/book", controller.UpdateBook)
	route.Delete("/book", controller.DeleteBook)

}