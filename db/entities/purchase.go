package entities

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetNearbyMerchantQueries struct {
	MerchantId       string  `db:"id" json:"merchantId" query:"merchantId"`
	Limit            int     `json:"limit" query:"limit"`
	Offset           int     `json:"offset" query:"offset"`
	Name             string  `db:"name" json:"name" query:"name"`
	MerchantCategory string  `db:"merchant_category" json:"merchantCategory" query:"merchantCategory"`
	Latitude         float64 `db:"latitude" json:"lat"`
	Longitude        float64 `db:"longitude" json:"long"`
}

type GetNearbyMerchantResponse struct {
	Merchant GetMerchantResponse `json:"merchant"`
	Items    []GetItemResponse   `json:"items"`
}

type GetEstimatePayload struct {
	Location Location `json:"userLocation" form:"userLocation"`
	Orders   []Order  `json:"orders" form:"orders"`
}

type Order struct {
	MerchantId      string      `json:"merchantId" form:"merchantId"`
	IsStartingPoint bool        `json:"isStartingPoint" form:"isStartingPoint"`
	Items           []OrderItem `json:"items" form:"items"`
}

type OrderItem struct {
	ItemId   string `json:"itemId" form:"itemId"`
	Quantity int    `json:"quantity" form:"quantity"`
}

type GetEstimateResponse struct {
	TotalPrice                     int     `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes float64 `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateId           string  `json:"calculatedEstimateId"`
}

type PlaceOrderPayload struct {
	CalculatedEstimateId string `json:"calculatedEstimateId"`
}

type PlaceOrderResponse struct {
	OrderId string `json:"orderId"`
}

type GetUserOrderQueries struct {
	MerchantId       string `db:"id" json:"merchantId" query:"merchantId"`
	Limit            int    `json:"limit" query:"limit"`
	Offset           int    `json:"offset" query:"offset"`
	Name             string `db:"name" json:"name" query:"name"` //fort both either merchant name or item name
	MerchantCategory string `db:"merchant_category" json:"merchantCategory" query:"merchantCategory"`
	UserId           string `json:"userId"`
}

type RoutePoint struct {
	Latitude  float64
	Longitude float64
}

type OrderInfo struct {
	UserId                         string  `db:"user_id" json:"userId"`
	TotalPrice                     int     `db:"total_price" json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes float64 `db:"estimated_delivery_time_in_minutes" json:"estimatedDeliveryTimeInMinutes"`
	Status                         string  `db:"status" json:"status"`
}

type GetUserOrderResponse struct {
	OrderId string          `json:"orderId"`
	Orders  []OrderResponse `json:"orders"`
}

type OrderResponse struct {
	Merchant GetMerchantResponse `json:"merchant"`
	Items    []ItemResponse      `json:"items"`
}

type ItemResponse struct {
	ItemId          string `db:"id" json:"itemId"`
	Name            string `db:"name" json:"name"`
	Price           int    `db:"price" json:"price"`
	Quantity        int    `db:"quantity" json:"quantity"`
	ImageUrl        string `db:"image_url" json:"imageUrl"`
	ProductCategory string `db:"gender" json:"productCategory"`
	CreatedAt       string `db:"created_at" json:"createdAt"`
}

func (u *GetEstimatePayload) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Location,
			validation.Required.Error("location is required"),
			validation.By(ValidateLocation),
		),
		validation.Field(&u.Orders, validation.Required.Error("orders are required"),
			validation.Each(validation.WithContext(validateOrder))),
	)

	if err != nil {
		return err
	}

	// Check that there is exactly one order with IsStartingPoint == true
	startingPointCount := 0
	for _, order := range u.Orders {
		if order.IsStartingPoint {
			startingPointCount++
		}
	}

	if startingPointCount != 1 {
		return errors.New("there must be exactly one starting point order")
	}

	return nil
}

func validateOrder(ctx context.Context, value interface{}) error {
	order, ok := value.(Order)
	if !ok {
		return errors.New("invalid order type")
	}

	err := validation.ValidateStruct(&order,
		validation.Field(&order.MerchantId,
			validation.Required.Error("merchantId is required"),
			validation.By(ValidateUUID),
		),
		validation.Field(&order.Items, validation.Required.Error("items are required"),
			validation.Each(validation.WithContext(validateOrderItem))),
	)

	return err
}

func validateOrderItem(ctx context.Context, value interface{}) error {
	item, ok := value.(OrderItem)
	if !ok {
		return errors.New("invalid item type")
	}

	return validation.ValidateStruct(&item,
		validation.Field(&item.ItemId,
			validation.Required.Error("itemId is required"),
			validation.By(ValidateUUID),
		),
		validation.Field(&item.Quantity,
			validation.Required.Error("quantity is required"),
			validation.Min(1).Error("quantity must be at least 1"),
		),
	)
}

func (u *PlaceOrderPayload) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.CalculatedEstimateId,
			validation.Required.Error("calculatedEstimateId latitude is required"),
		),
	)

	return err
}
