package controller

import (
	"beli_mang/db/entities"
	"beli_mang/svc"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	svc svc.UserSvc
}

func NewUserController(svc svc.UserSvc) *UserController {
	return &UserController{svc: svc}
}

type registerResponse struct {
	AccessToken string `json:"token,omitempty"`
}

type loginResponse struct {
	AccessToken string `json:"token"`
}

type simpleResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

func (c *UserController) AdminRegister(ctx *fiber.Ctx) error {
	var newUser entities.RegistrationPayload
	if err := ctx.BodyParser(&newUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	accessToken, err := c.svc.AdminRegister(ctx.Context(), newUser)
	if err != nil {
		return err
	}

	respData := registerResponse{
		AccessToken: accessToken,
	}

	return ctx.Status(fiber.StatusCreated).JSON(respData)
}

func (c *UserController) AdminLogin(ctx *fiber.Ctx) error {
	var user entities.RegistrationPayload
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	loginPayload := entities.Credential{
		Username: user.Username,
		Password: user.Password,
	}

	accessToken, err := c.svc.AdminLogin(ctx.Context(), loginPayload)
	if err != nil {
		return err
	}

	respData := loginResponse{
		AccessToken: accessToken,
	}

	return ctx.Status(fiber.StatusOK).JSON(respData)
}

func (c *UserController) UserRegister(ctx *fiber.Ctx) error {
	var newUser entities.RegistrationPayload
	if err := ctx.BodyParser(&newUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	accessToken, err := c.svc.UserRegister(ctx.Context(), newUser)
	if err != nil {
		return err
	}

	respData := registerResponse{
		AccessToken: accessToken,
	}

	return ctx.Status(fiber.StatusCreated).JSON(respData)
}

func (c *UserController) UserLogin(ctx *fiber.Ctx) error {
	var user entities.RegistrationPayload
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	loginPayload := entities.Credential{
		Username: user.Username,
		Password: user.Password,
	}

	accessToken, err := c.svc.UserLogin(ctx.Context(), loginPayload)
	if err != nil {
		return err
	}

	respData := loginResponse{
		AccessToken: accessToken,
	}

	return ctx.Status(fiber.StatusOK).JSON(respData)
}
