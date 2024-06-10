package controller

import (
	"beli_mang/responses"
	"beli_mang/svc"

	"github.com/gofiber/fiber/v2"
)

type ImageController struct {
	svc svc.ImageSvc
}

func NewImageController(svc svc.ImageSvc) *ImageController {
	return &ImageController{
		svc: svc,
	}
}

func (c *ImageController) UploadImage(ctx *fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("file")
	if fileHeader == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.NewBadRequestError("file should not be empty"))
	}
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalServerError("failed to retrieve file"))
	}

	url, err := c.svc.UploadImage(fileHeader)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"message": "File uploaded successfully",
		"data": map[string]interface{}{
			"imageUrl": url,
		},
	})
}
