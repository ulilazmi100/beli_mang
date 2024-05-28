package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Item struct {
	ItemId          string    `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	Price           int       `db:"price" json:"price"`
	ProductCategory string    `db:"product_category" json:"productCategory"`
	ImageUrl        string    `db:"image_url" json:"imageUrl"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
}

type ItemRegistrationPayload struct {
	MerchantId      string `db:"merchant_id" json:"merchantId" query:"merchantId"`
	Name            string `db:"name" json:"name" form:"name"`
	Price           int    `json:"price" form:"price"`
	ProductCategory string `db:"product_category" json:"productCategory" form:"productCategory"`
	ImageUrl        string `db:"image_url" json:"imageUrl" form:"imageUrl"`
}

type GetItemQueries struct {
	ItemId          string `db:"id" json:"itemId" query:"itemId"`
	Limit           int    `json:"limit" query:"limit"`
	Offset          int    `json:"offset" query:"offset"`
	Name            string `db:"name" json:"name" query:"name"`
	ProductCategory string `db:"product_category" json:"productCategory" query:"productCategory"`
	CreatedAt       string `db:"created_at" json:"createdAt" query:"createdAt"`
}

type GetItemResponse struct {
	ItemId          string `db:"id" json:"itemId"`
	Name            string `db:"name" json:"name"`
	Price           int    `db:"price" json:"price"`
	ImageUrl        string `db:"image_url" json:"imageUrl"`
	ProductCategory string `db:"gender" json:"productCategory"`
	CreatedAt       string `db:"created_at" json:"createdAt"`
}

func (u *ItemRegistrationPayload) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Name,
			validation.Required.Error("name is required"),
			validation.Length(2, 30).Error("name must be between 2 and 30 characters"),
		),
		validation.Field(&u.ProductCategory,
			validation.Required.Error("productCategory is required"),
			validation.In("Beverage", "Food", "Snack", "Condiments", "Additions"),
		),
		validation.Field(&u.ImageUrl,
			validation.Required.Error("imageUrl is required"),
			validation.By(ValidateImageURL),
		),
		validation.Field(&u.Price,
			validation.Required.Error("latitude is required"),
			validation.Min(1).Error("price min 1"),
		),
	)

	return err
}
