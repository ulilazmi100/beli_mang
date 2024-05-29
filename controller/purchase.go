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

	var total int
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
