package main

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	authorizationHeader = "authorization"
	bearer              = "bearer"
	payloadHeader       = "payload"
)

func (app *App) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get(authorizationHeader)
		if authorization == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
		}
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != bearer {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header format")
		}

		token := parts[1]
		payload, err := app.jwtMaker.VerifyToken(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
		}

		c.Locals(payloadHeader, payload)

		return c.Next()
	}
}

func (app *App) LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		log.Printf("[%s] %s - %v\n", c.Method(), c.Path(), duration)

		return err
	}
}
