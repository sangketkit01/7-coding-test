package main

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func (app *App) routes() *fiber.App{
	router := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		ErrorHandler: errorHandler,
	})

	router.Get("/", func(c *fiber.Ctx) error{
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message" : "Hello world"})
	})

	router.Use(app.LoggingMiddleware())
	router.Post("/grpc/create-user", app.CreateUserViaGrpc)

	router.Post("/create-user", app.CreateUser)
	router.Post("/login-user", app.LoginUser)

	authRouter := router.Group("/", app.AuthMiddleware())
	authRouter.Get("/get-user/:id", app.FetchUserById)
	authRouter.Get("/all-users", app.ListAllUsers)
	authRouter.Put("/update-user", app.UpdateUser)
	authRouter.Delete("/delete-user", app.DeleteUser)

	authRouter.Get("/grpc/get-user/:id", app.GetUserViaGrpc)

	return router
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error": msg,
	})
}