package repo

import (
	"beli_mang/db/entities"
	"context"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseRepo interface {
	GetNearbyMerchants(ctx context.Context, filter entities.GetNearbyMerchantQueries) ([]entities.GetNearbyMerchantResponse, error)
}

type purchaseRepo struct {
	db *pgxpool.Pool
}

func NewPurchaseRepo(db *pgxpool.Pool) PurchaseRepo {
	return &purchaseRepo{db}
}

func (r *purchaseRepo) GetNearbyMerchants(ctx context.Context, filter entities.GetNearbyMerchantQueries) ([]entities.GetNearbyMerchantResponse, error) {
	var merchants []entities.GetNearbyMerchantResponse

	query := "SELECT id, name, merchant_category, image_url, latitude, longitude, created_at FROM merchants"

	query = `
    SELECT id, name, merchant_category, image_url, latitude, longitude, created_at,
        (acos(
            cos(radians($1)) * cos(radians(latitude)) *
            cos(radians(longitude) - radians($2)) +
            sin(radians($1)) * sin(radians(latitude))
        )) AS distance
    FROM merchants
    ORDER BY distance;` //No need for earth constant because they all would be multiplied with the same constant

	query += getNearbyMerchantConstructWhereQuery(filter)

	query += " limit $3 offset $4"

	rows, err := r.db.Query(ctx, query, filter.Latitude, filter.Longitude, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		merchant := entities.GetNearbyMerchantResponse{}
		var createdAt time.Time
		err := rows.Scan(&merchant.Merchant.MerchantId, &merchant.Merchant.Name, &merchant.Merchant.MerchantCategory, &merchant.Merchant.ImageUrl, &merchant.Merchant.Location.Latitude, &merchant.Merchant.Location.Longitude, &createdAt)
		if err != nil {
			return nil, err
		}
		merchant.Merchant.CreatedAt = createdAt.Format(time.RFC3339Nano)

		//TODO: Asynchronize GET items
		getItemsQuery := "SELECT id, name, product_category, price, image_url, created_at FROM items WHERE merchant_id = $1"

		// getItemsQuery += " ORDER BY created_at DESC"

		rows, err := r.db.Query(ctx, getItemsQuery, merchant.Merchant.MerchantId)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			item := entities.GetItemResponse{}
			var createdAtItem time.Time

			err := rows.Scan(&item.ItemId, &item.Name, &item.ProductCategory, &item.Price, &item.ImageUrl, &createdAtItem)
			if err != nil {
				return nil, err
			}

			item.CreatedAt = createdAtItem.Format(time.RFC3339Nano)

			merchant.Items = append(merchant.Items, item)
		}

		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

func getNearbyMerchantConstructWhereQuery(filter entities.GetNearbyMerchantQueries) string {
	whereSQL := []string{}

	if filter.MerchantId != "" {
		whereSQL = append(whereSQL, " id = '"+filter.MerchantId+"'")
	}

	if validation.Validate(&filter.MerchantCategory,
		validation.In("SmallRestaurant", "MediumRestaurant", "LargeRestaurant", "MerchandiseRestaurant", "BoothKiosk", "ConvenienceStore"),
	) == nil {
		whereSQL = append(whereSQL, " purchase_category = '"+filter.MerchantCategory+"'")
	}

	if filter.Name != "" {
		whereSQL = append(whereSQL, " name ILIKE '%"+filter.Name+"%'")
	}

	if len(whereSQL) > 0 {
		return " WHERE " + strings.Join(whereSQL, " AND ")
	}

	return ""
}
