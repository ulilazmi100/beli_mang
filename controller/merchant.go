package controller

import (
	"beli_mang/db/entities"
	"beli_mang/svc"

	"github.com/gofiber/fiber/v2"
)

type MerchantController struct {
	svc svc.MerchantSvc
}

type registerMerchantResponse struct {
	MerchantId string `json:"merchantId"`
}

type registerItemResponse struct {
	ItemId string `json:"itemId"`
}

type getResponse struct {
	Data     interface{} `json:"data"`
	Metadata metadata    `json:"meta"`
}

type metadata struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

func NewMerchantController(svc svc.MerchantSvc) *MerchantController {
	return &MerchantController{svc: svc}
}

func (c *MerchantController) RegisterMerchant(ctx *fiber.Ctx) error {
	var newMerchant entities.MerchantRegistrationPayload
	if err := ctx.BodyParser(&newMerchant); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	id, err := c.svc.RegisterMerchant(ctx.Context(), newMerchant)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(registerMerchantResponse{
		MerchantId: id,
	})
}

func (c *MerchantController) GetMerchant(ctx *fiber.Ctx) error {
	var merchantQuery entities.GetMerchantQueries
	if err := ctx.QueryParser(&merchantQuery); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if merchantQuery.Limit == 0 {
		merchantQuery.Limit = 5
	}

	if merchantQuery.Limit < 0 || merchantQuery.Offset < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid query param")
	}
	resp, total, err := c.svc.GetMerchant(ctx.Context(), merchantQuery)
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(getResponse{
			Data: []interface{}{},
			Metadata: metadata{
				Limit:  merchantQuery.Limit,
				Offset: merchantQuery.Offset,
				Total:  0,
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(getResponse{
		Data: resp,
		Metadata: metadata{
			Limit:  merchantQuery.Limit,
			Offset: merchantQuery.Offset,
			Total:  total,
		},
	})
}

func (c *MerchantController) RegisterItem(ctx *fiber.Ctx) error {
	var newItem entities.ItemRegistrationPayload
	if err := ctx.BodyParser(&newItem); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	newItem.MerchantId = ctx.Params("merchantId")

	id, err := c.svc.RegisterItem(ctx.Context(), newItem)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(registerItemResponse{
		ItemId: id,
	})
}

func (c *MerchantController) GetItem(ctx *fiber.Ctx) error {
	var itemQuery entities.GetItemQueries

	if err := ctx.QueryParser(&itemQuery); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	itemQuery.MerchantId = ctx.Params("merchantId")

	if itemQuery.Limit == 0 {
		itemQuery.Limit = 5
	}

	if itemQuery.Limit < 0 || itemQuery.Offset < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid query param")
	}

	resp, total, err := c.svc.GetItem(ctx.Context(), itemQuery)
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(getResponse{
			Data: []interface{}{},
			Metadata: metadata{
				Limit:  itemQuery.Limit,
				Offset: itemQuery.Offset,
				Total:  0,
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(getResponse{
		Data: resp,
		Metadata: metadata{
			Limit:  itemQuery.Limit,
			Offset: itemQuery.Offset,
			Total:  total,
		},
	})
}
