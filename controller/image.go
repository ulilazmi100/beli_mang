package controller

import (
	"beli_mang/responses"
	"beli_mang/svc"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ImageController struct {
	svc svc.ImageSvc
}

func NewImageController(svc svc.ImageSvc) *ImageController {
	return &ImageController{
		svc: svc,
	}
}

func (c *ImageController) UploadImage(ctx echo.Context) error {
	fileHeader, err := ctx.FormFile("file")
	if fileHeader == nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewBadRequestError("file should not be empty"))
	}
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, responses.NewInternalServerError("failed to retrieve file"))
	}

	url, err := c.svc.UploadImage(fileHeader)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "File uploaded successfully",
		"data": map[string]interface{}{
			"imageUrl": url,
		},
	})
}
