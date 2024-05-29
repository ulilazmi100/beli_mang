package controller

import (
	"beli_mang/db/entities"
	"beli_mang/responses"
	"beli_mang/svc"
	"net/http"

	"github.com/labstack/echo/v4"
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
	Metadata metadata    `json:"metadata"`
}

type metadata struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

func NewMerchantController(svc svc.MerchantSvc) *MerchantController {
	return &MerchantController{svc: svc}
}

func (c *MerchantController) RegisterMerchant(ctx echo.Context) error {
	var newMerchant entities.MerchantRegistrationPayload
	if err := ctx.Bind(&newMerchant); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	id, err := c.svc.RegisterMerchant(ctx.Request().Context(), newMerchant)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, registerMerchantResponse{
		MerchantId: id,
	})
}

func (c *MerchantController) GetMerchant(ctx echo.Context) error {
	var total int
	var merchantQuery entities.GetMerchantQueries
	if err := ctx.Bind(&merchantQuery); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	if merchantQuery.Limit == 0 {
		merchantQuery.Limit = 5
	}

	if merchantQuery.Limit < 0 || merchantQuery.Offset < 0 {
		return responses.NewBadRequestError("invalid query param")
	}
	resp, err := c.svc.GetMerchant(ctx.Request().Context(), merchantQuery)
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return ctx.JSON(http.StatusOK, getResponse{
			Data: []interface{}{},
			Metadata: metadata{
				Limit:  merchantQuery.Limit,
				Offset: merchantQuery.Offset,
				Total:  total,
			},
		})
	}

	return ctx.JSON(http.StatusOK, getResponse{
		Data: resp,
		Metadata: metadata{
			Limit:  merchantQuery.Limit,
			Offset: merchantQuery.Offset,
			Total:  total,
		},
	})
}

func (c *MerchantController) RegisterItem(ctx echo.Context) error {
	var newItem entities.ItemRegistrationPayload
	if err := ctx.Bind(&newItem); err != nil {
		return responses.NewBadRequestError(err.Error())
	}
	newItem.MerchantId = ctx.Param("merchantId")

	id, err := c.svc.RegisterItem(ctx.Request().Context(), newItem)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, registerItemResponse{
		ItemId: id,
	})
}

func (c *MerchantController) GetItem(ctx echo.Context) error {
	var total int
	var itemQuery entities.GetItemQueries

	if err := ctx.Bind(&itemQuery); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	itemQuery.MerchantId = ctx.Param("merchantId")

	if itemQuery.Limit == 0 {
		itemQuery.Limit = 5
	}

	if itemQuery.Limit < 0 || itemQuery.Offset < 0 {
		return responses.NewBadRequestError("invalid query param")
	}

	resp, err := c.svc.GetItem(ctx.Request().Context(), itemQuery)
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return ctx.JSON(http.StatusOK, getResponse{
			Data: []interface{}{},
			Metadata: metadata{
				Limit:  itemQuery.Limit,
				Offset: itemQuery.Offset,
				Total:  total,
			},
		})
	}

	return ctx.JSON(http.StatusOK, getResponse{
		Data: resp,
		Metadata: metadata{
			Limit:  itemQuery.Limit,
			Offset: itemQuery.Offset,
			Total:  total,
		},
	})
}
