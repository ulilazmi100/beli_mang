package middleware

import (
	"beli_mang/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if customError, ok := err.(responses.CustomError); ok {
		code = customError.Status()
		c.JSON(code, map[string]interface{}{
			"message": customError.Error(),
		})
		return
	}

	if httpError, ok := err.(*echo.HTTPError); ok {
		code = httpError.Code
		if httpError.Internal != nil {
			err = httpError.Internal
		} else {
			err = httpError
		}
	}

	if !c.Response().Committed {
		c.JSON(code, map[string]interface{}{
			"message": err.Error(),
		})
	}
}
