package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Merchant struct {
	Id               string    `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`
	Latitude         float64   `db:"latitude" json:"latitude"`
	Longitude        float64   `db:"longitude" json:"longitude"`
	MerchantCategory string    `db:"merchant_category" json:"merchantCategory"`
	ImageUrl         string    `db:"image_url" json:"imageUrl"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
}

type MerchantRegistrationPayload struct {
	Name             string   `db:"name" json:"name" form:"name"`
	Location         Location `json:"location" form:"location"`
	MerchantCategory string   `db:"merchant_category" json:"merchantCategory" form:"merchantCategory"`
	ImageUrl         string   `db:"image_url" json:"imageUrl" form:"imageUrl"`
}

type Location struct {
	Latitude  float64 `db:"latitude" json:"latitude" form:"latitude"`
	Longitude float64 `db:"longitude" json:"longitude" form:"longitude"`
}

type GetMerchantQueries struct {
	MerchantId       string `db:"id" json:"merchantId" query:"merchantId"`
	Limit            int    `json:"limit" query:"limit"`
	Offset           int    `json:"offset" query:"offset"`
	Name             string `db:"name" json:"name" query:"name"`
	MerchantCategory string `db:"merchant_category" json:"merchantCategory" query:"merchantCategory"`
	CreatedAt        string `db:"created_at" json:"createdAt" query:"createdAt"`
}

type GetMerchantResponse struct {
	MerchantId       string   `db:"id" json:"merchantId"`
	Name             string   `db:"name" json:"name"`
	Location         Location `json:"location"`
	ImageUrl         string   `db:"image_url" json:"imageUrl"`
	MerchantCategory string   `db:"gender" json:"merchantCategory"`
	CreatedAt        string   `db:"created_at" json:"createdAt"`
}

func (u *MerchantRegistrationPayload) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Name,
			validation.Required.Error("name is required"),
			validation.Length(2, 30).Error("name must be between 2 and 30 characters"),
		),
		validation.Field(&u.MerchantCategory,
			validation.Required.Error("merchantCategory is required"),
			validation.In("SmallRestaurant", "MediumRestaurant", "LargeRestaurant", "MerchandiseRestaurant", "BoothKiosk", "ConvenienceStore"),
		),
		validation.Field(&u.ImageUrl,
			validation.Required.Error("imageUrl is required"),
			validation.By(ValidateImageURL),
		),
		validation.Field(&u.Location.Latitude,
			validation.Required.Error("latitude is required"),
		),
		validation.Field(&u.Location.Longitude,
			validation.Required.Error("longitude is required"),
		),
	)

	return err
}
