package controller

import (
	"beli_mang/db/entities"
	"beli_mang/responses"
	"beli_mang/svc"
	"net/http"

	"github.com/labstack/echo/v4"
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

func (c *UserController) AdminRegister(ctx echo.Context) error {
	var newUser entities.RegistrationPayload
	if err := ctx.Bind(&newUser); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	accessToken, err := c.svc.AdminRegister(ctx.Request().Context(), newUser)
	if err != nil {
		return err
	}

	respData := registerResponse{
		AccessToken: accessToken,
	}

	return ctx.JSON(http.StatusCreated, respData)
}

func (c *UserController) AdminLogin(ctx echo.Context) error {
	var user entities.RegistrationPayload
	if err := ctx.Bind(&user); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	loginPayload := entities.Credential{
		Username: user.Username,
		Password: user.Password,
	}

	accessToken, err := c.svc.AdminLogin(ctx.Request().Context(), loginPayload)
	if err != nil {
		return err
	}

	respData := loginResponse{
		AccessToken: accessToken,
	}

	return ctx.JSON(http.StatusOK, respData)
}

func (c *UserController) UserRegister(ctx echo.Context) error {
	var newUser entities.RegistrationPayload
	if err := ctx.Bind(&newUser); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	accessToken, err := c.svc.UserRegister(ctx.Request().Context(), newUser)
	if err != nil {
		return err
	}

	respData := registerResponse{
		AccessToken: accessToken,
	}

	return ctx.JSON(http.StatusCreated, respData)
}

func (c *UserController) UserLogin(ctx echo.Context) error {
	var user entities.RegistrationPayload
	if err := ctx.Bind(&user); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	loginPayload := entities.Credential{
		Username: user.Username,
		Password: user.Password,
	}

	accessToken, err := c.svc.UserLogin(ctx.Request().Context(), loginPayload)
	if err != nil {
		return err
	}

	respData := loginResponse{
		AccessToken: accessToken,
	}

	return ctx.JSON(http.StatusOK, respData)
}
