package controller

import (
	"beli_mang/db/entities"
	"beli_mang/responses"
	"beli_mang/svc"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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

func (c *PurchaseController) GetNearbyMerchant(ctx echo.Context) error {

	var merchantQuery entities.GetNearbyMerchantQueries
	if err := ctx.Bind(&merchantQuery); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	var err error

	merchantQuery.Latitude, err = strconv.ParseFloat(ctx.Param("lat"), 64)
	if err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	merchantQuery.Latitude, err = strconv.ParseFloat(ctx.Param("long"), 64)
	if err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	if merchantQuery.Limit == 0 {
		merchantQuery.Limit = 5
	}

	if merchantQuery.Limit < 0 || merchantQuery.Offset < 0 {
		return responses.NewBadRequestError("invalid query param")
	}
	resp, err := c.svc.GetNearbyMerchant(ctx.Request().Context(), merchantQuery)
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return ctx.JSON(http.StatusOK, simpleResponse{
			Data: []interface{}{},
		})
	}

	return ctx.JSON(http.StatusOK, simpleResponse{
		Data: resp,
	})
}

func (c *PurchaseController) EstimateOrder(ctx echo.Context) error {

	var getEstimatePayload entities.GetEstimatePayload
	if err := ctx.Bind(&getEstimatePayload); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	var err error

	resp, err := c.svc.GetOrderEstimation(ctx.Request().Context(), getEstimatePayload, ctx.Get("user_id").(string))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (c *PurchaseController) PlaceOrder(ctx echo.Context) error {

	var placeOrderPayload entities.PlaceOrderPayload
	if err := ctx.Bind(&placeOrderPayload); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	var err error

	resp, err := c.svc.PlaceOrder(ctx.Request().Context(), placeOrderPayload)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (c *PurchaseController) GetOrder(ctx echo.Context) error {

	var getOrderQuery entities.GetUserOrderQueries
	if err := ctx.Bind(&getOrderQuery); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	var err error

	if getOrderQuery.Limit == 0 {
		getOrderQuery.Limit = 5
	}

	if getOrderQuery.Limit < 0 || getOrderQuery.Offset < 0 {
		return responses.NewBadRequestError("invalid query param")
	}

	getOrderQuery.UserId = ctx.Get("user_id").(string)

	resp, err := c.svc.GetOrder(ctx.Request().Context(), getOrderQuery)
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return ctx.JSON(http.StatusOK, simpleResponse{
			Data: []interface{}{},
		})
	}

	return ctx.JSON(http.StatusOK, simpleResponse{
		Data: resp,
	})
}
