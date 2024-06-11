package controller

import (
	"beli_mang/db/entities"
	"beli_mang/svc"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type PurchaseController struct {
	svc svc.PurchaseSvc
}

type registerPurchaseResponse struct {
	PurchaseId string `json:"purchaseId"`
}

func NewPurchaseController(svc svc.PurchaseSvc) *PurchaseController {
	return &PurchaseController{svc: svc}
}

func (c *PurchaseController) GetNearbyMerchant(ctx *fiber.Ctx) error {

	var merchantQuery entities.GetNearbyMerchantQueries
	if err := ctx.QueryParser(&merchantQuery); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	var err error

	latlong := ctx.Params("latlong")
	parts := strings.Split(latlong, ",")
	if len(parts) != 2 {
		return ctx.Status(fiber.StatusBadRequest).JSON("Invalid latlong format")
	}

	merchantQuery.Latitude, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("Invalid latitude")
	}

	merchantQuery.Longitude, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("Invalid longitude")
	}

	if merchantQuery.Limit == 0 {
		merchantQuery.Limit = 5
	}

	if merchantQuery.Limit < 0 || merchantQuery.Offset < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid query param")
	}
	resp, total, err := c.svc.GetNearbyMerchant(ctx.Context(), merchantQuery)
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

func (c *PurchaseController) EstimateOrder(ctx *fiber.Ctx) error {

	var getEstimatePayload entities.GetEstimatePayload
	if err := ctx.BodyParser(&getEstimatePayload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	var err error

	resp, err := c.svc.GetOrderEstimation(ctx.Context(), getEstimatePayload, ctx.Locals("user_id").(string))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (c *PurchaseController) PlaceOrder(ctx *fiber.Ctx) error {

	var placeOrderPayload entities.PlaceOrderPayload
	if err := ctx.BodyParser(&placeOrderPayload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	var err error

	resp, err := c.svc.PlaceOrder(ctx.Context(), placeOrderPayload)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (c *PurchaseController) GetOrder(ctx *fiber.Ctx) error {

	var getOrderQuery entities.GetUserOrderQueries
	if err := ctx.QueryParser(&getOrderQuery); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	var err error

	if getOrderQuery.Limit == 0 {
		getOrderQuery.Limit = 5
	}

	if getOrderQuery.Limit < 0 || getOrderQuery.Offset < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid query param")
	}

	getOrderQuery.UserId = ctx.Locals("user_id").(string)

	resp, err := c.svc.GetOrder(ctx.Context(), getOrderQuery)
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return ctx.Status(fiber.StatusOK).JSON([]interface{}{})
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
