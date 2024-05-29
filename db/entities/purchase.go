package entities

type GetNearbyMerchantQueries struct {
	MerchantId       string  `db:"id" json:"merchantId" query:"merchantId"`
	Limit            int     `json:"limit" query:"limit"`
	Offset           int     `json:"offset" query:"offset"`
	Name             string  `db:"name" json:"name" query:"name"`
	MerchantCategory string  `db:"merchant_category" json:"merchantCategory" query:"merchantCategory"`
	Latitude         float64 `db:"latitude" json:"latitude"`
	Longitude        float64 `db:"longitude" json:"longitude"`
}

type GetNearbyMerchantResponse struct {
	Merchant GetMerchantResponse `json:"merchant"`
	Items    []GetItemResponse   `json:"items"`
}
