package middleware

import (
	"errors"
	"strings"

	"beli_mang/crypto"
	"beli_mang/responses"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		if auth == "" {
			return responses.NewUnauthorizedError("token not found")
		}

		splitted := strings.Split(auth, " ")
		if len(splitted) != 2 || splitted[0] != "Bearer" {
			return responses.NewUnauthorizedError("invalid token")
		}

		payload, err := crypto.VerifyToken(splitted[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return responses.NewUnauthorizedError("token expired")
			}
			return responses.NewUnauthorizedError(err.Error())
		}

		if payload.Role == "admin" {
			return responses.NewUnauthorizedError("user is not an admin")
		}

		c.Set("user_id", payload.Id)
		c.Set("username", payload.Username)
		c.Set("role", payload.Role)

		return next(c)
	}
}

func UserAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		if auth == "" {
			return responses.NewUnauthorizedError("token not found")
		}

		splitted := strings.Split(auth, " ")
		if len(splitted) != 2 || splitted[0] != "Bearer" {
			return responses.NewUnauthorizedError("invalid token")
		}

		payload, err := crypto.VerifyToken(splitted[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return responses.NewUnauthorizedError("token expired")
			}
			return responses.NewUnauthorizedError(err.Error())
		}

		if payload.Role == "user" {
			return responses.NewUnauthorizedError("user is not a user")
		}

		c.Set("user_id", payload.Id)
		c.Set("username", payload.Username)
		c.Set("role", payload.Role)

		return next(c)
	}
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		if auth == "" {
			return responses.NewUnauthorizedError("token not found")
		}

		splitted := strings.Split(auth, " ")
		if len(splitted) != 2 || splitted[0] != "Bearer" {
			return responses.NewUnauthorizedError("invalid token")
		}

		payload, err := crypto.VerifyToken(splitted[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return responses.NewUnauthorizedError("token expired")
			}
			return responses.NewUnauthorizedError(err.Error())
		}

		c.Set("user_id", payload.Id)
		c.Set("username", payload.Username)
		c.Set("role", payload.Role)

		return next(c)
	}
}
