package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sangketkit01/7-coding-test/internal/db"
	"github.com/sangketkit01/7-coding-test/internal/token"
)

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,alphanum"`
}

func (app *App) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user request")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err := app.model.Insert(db.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"message": "create user successfully."})
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,alphanum"`
}

type LoginUserResponse struct {
	Token     string    `json:"token"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (app *App) LoginUser(c *fiber.Ctx) error {
	var req LoginUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user request")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user := db.User{
		Email:    req.Email,
		Password: req.Password,
	}

	loggedInUser, err := user.LoginUser()
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("invalid credentials: %v", err))
	}

	token, payload, err := app.jwtMaker.CreateToken(loggedInUser.ID, time.Hour*24*7)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := LoginUserResponse{
		Token:     token,
		Email:     loggedInUser.Email,
		IssuedAt:  payload.IssuedAt,
		ExpiredAt: payload.ExpiredAt,
	}

	return c.JSON(response)
}

func (app *App) FetchUserById(c *fiber.Ctx) error {
	p := c.Locals(payloadHeader)

	_, ok := p.(*token.Payload)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid payload")
	}

	userId := c.Params("id", "")
	if userId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "user id is not provided.")
	}

	user, err := app.model.FetchUserByID(userId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"user": user})
}

func (app *App) ListAllUsers(c *fiber.Ctx) error {
	p := c.Locals(payloadHeader)

	_, ok := p.(*token.Payload)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid payload")
	}

	users, err := app.model.ListAllUsers()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"users": users})
}

type UpdateUserRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
	Name  string `json:"name" validate:"omitempty"`
}


type UpdateUserResponse struct {
	NewToken  string    `json:"new_token"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (app *App) UpdateUser(c *fiber.Ctx) error {
	p := c.Locals(payloadHeader)

	payload, ok := p.(*token.Payload)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid payload")
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user request")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	
	user, err := app.model.FetchUserByID(payload.ID.Hex())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	
	if user.ID != payload.ID {
		fmt.Println(user.Email, req.Email)
		return fiber.NewError(fiber.StatusForbidden, "you are not allowed to update this user")
	}

	user.Email = req.Email
	user.Name = req.Name

	err = user.UpdateUser()
	if err != nil {
		log.Printf("updated user failed: %v\n", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	newToken, newPayload, err := app.jwtMaker.CreateToken(user.ID, time.Hour * 24 * 7)
	if err != nil{
		log.Printf("create token failed: %v\n", err)
		return fiber.NewError(fiber.StatusInternalServerError, "cannot create new token")
	}

	newUser, err := app.model.FetchUserByID(newPayload.ID.Hex())
	if err != nil{
		log.Println("failed to fetch user's data:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "cannot fetch user data")
	}

	c.Locals("payload", newPayload)

	response := UpdateUserResponse{
		NewToken: newToken,
		Email: newUser.Email,
		Name: newUser.Name,
		IssuedAt: newPayload.IssuedAt,
		ExpiredAt: newPayload.ExpiredAt,
	}

	return c.JSON(response)
}

func (app *App) DeleteUser(c *fiber.Ctx) error{
	p := c.Locals(payloadHeader)

	payload, ok := p.(*token.Payload)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid payload")
	}

	user, err := app.model.FetchUserByID(payload.ID.Hex())
	if err != nil{
		return fiber.NewError(fiber.StatusInternalServerError, "cannot get user to delete")
	}

	if err = user.DeleteUser() ; err != nil{
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to delete user: %v\n", err))
	}

	return c.JSON(fiber.Map{"message" : "Delete user successfully."})
}

