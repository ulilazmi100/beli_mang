package repo

import (
	"beli_mang/db/entities"
	"context"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MerchantRepo interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	GetMerchant(ctx context.Context, tx pgx.Tx, merchantId string) (string, error)
	CreateMerchant(ctx context.Context, tx pgx.Tx, merchant *entities.MerchantRegistrationPayload) (string, error)
	GetMerchants(ctx context.Context, filter entities.GetMerchantQueries) ([]entities.GetMerchantResponse, int, error)
	CreateItem(ctx context.Context, tx pgx.Tx, item *entities.ItemRegistrationPayload) (string, error)
	GetItem(ctx context.Context, filter entities.GetItemQueries) ([]entities.GetItemResponse, int, error)
}

type merchantRepo struct {
	db *pgxpool.Pool
}

func NewMerchantRepo(db *pgxpool.Pool) MerchantRepo {
	return &merchantRepo{db}
}

func (r *merchantRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

func (r *merchantRepo) GetMerchant(ctx context.Context, tx pgx.Tx, merchantId string) (string, error) {
	var id string
	query := "SELECT id FROM merchants WHERE id = $1"

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, merchantId)
	} else {
		row = r.db.QueryRow(ctx, query, merchantId)
	}
	err := row.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *merchantRepo) CreateMerchant(ctx context.Context, tx pgx.Tx, merchant *entities.MerchantRegistrationPayload) (string, error) {
	var id string
	statement := "INSERT INTO merchants (name, merchant_category, image_url, latitude, longitude) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, statement, merchant.Name, merchant.MerchantCategory, merchant.ImageUrl, merchant.Location.Latitude, merchant.Location.Longitude)
	} else {
		row = r.db.QueryRow(ctx, statement, merchant.Name, merchant.MerchantCategory, merchant.ImageUrl, merchant.Location.Latitude, merchant.Location.Longitude)
	}
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *merchantRepo) GetMerchants(ctx context.Context, filter entities.GetMerchantQueries) ([]entities.GetMerchantResponse, int, error) {
	var merchants []entities.GetMerchantResponse
	var createdAt time.Time
	var totalCount int
	query := "SELECT id, name, merchant_category, image_url, latitude, longitude, created_at FROM merchants"

	query += getMerchantConstructWhereQuery(filter)

	if filter.CreatedAt != "" {
		if filter.CreatedAt == "asc" {
			query += " ORDER BY created_at ASC"
		} else if filter.CreatedAt == "desc" {
			query += " ORDER BY created_at DESC"
		}
	} else {
		query += " ORDER BY created_at DESC"
	}

	query += " limit $1 offset $2"

	rows, err := r.db.Query(ctx, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		merchant := entities.GetMerchantResponse{}
		err := rows.Scan(&merchant.MerchantId, &merchant.Name, &merchant.MerchantCategory, &merchant.ImageUrl, &merchant.Location.Latitude, &merchant.Location.Longitude, &createdAt)
		if err != nil {
			return nil, 0, err
		}
		merchant.CreatedAt = createdAt.Format(time.RFC3339Nano)
		merchants = append(merchants, merchant)
	}

	return merchants, totalCount, nil
}

func (r *merchantRepo) CreateItem(ctx context.Context, tx pgx.Tx, item *entities.ItemRegistrationPayload) (string, error) {
	var id string
	statement := "INSERT INTO items (merchant_id, name, product_category, price, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, statement, item.MerchantId, item.Name, item.ProductCategory, item.Price, item.ImageUrl)
	} else {
		row = r.db.QueryRow(ctx, statement, item.MerchantId, item.Name, item.ProductCategory, item.Price, item.ImageUrl)
	}
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *merchantRepo) GetItem(ctx context.Context, filter entities.GetItemQueries) ([]entities.GetItemResponse, int, error) {
	var items []entities.GetItemResponse
	var createdAt time.Time
	var totalCount int
	query := "SELECT id, name, product_category, price, image_url, created_at, count(*) OVER() AS total_count FROM items"

	query += getItemConstructWhereQuery(filter)

	if filter.CreatedAt != "" {
		if filter.CreatedAt == "asc" {
			query += " ORDER BY created_at ASC"
		} else if filter.CreatedAt == "desc" {
			query += " ORDER BY created_at DESC"
		}
	} else {
		query += " ORDER BY created_at DESC"
	}

	query += " limit $1 offset $2"

	rows, err := r.db.Query(ctx, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		item := entities.GetItemResponse{}

		err := rows.Scan(&item.ItemId, &item.Name, &item.ProductCategory, &item.Price, &item.ImageUrl, &createdAt, &totalCount)
		if err != nil {
			return nil, 0, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339Nano)

		items = append(items, item)
	}

	return items, totalCount, nil
}

func getMerchantConstructWhereQuery(filter entities.GetMerchantQueries) string {
	whereSQL := []string{}

	err := validation.Validate(&filter.MerchantCategory, validation.In("SmallRestaurant", "MediumRestaurant", "LargeRestaurant",
		"MerchandiseRestaurant", "BoothKiosk", "ConvenienceStore"))
	if err != nil {
		filter.MerchantCategory = "" // Reset the category if it's invalid
	}

	if filter.MerchantId != "" {
		whereSQL = append(whereSQL, " id = '"+filter.MerchantId+"'")
	}

	if filter.MerchantCategory != "" {
		whereSQL = append(whereSQL, " merchant_category = '"+filter.MerchantCategory+"'")
	}

	if filter.Name != "" {
		whereSQL = append(whereSQL, " name ILIKE '%"+filter.Name+"%'")
	}

	if len(whereSQL) > 0 {
		return " WHERE " + strings.Join(whereSQL, " AND ")
	}

	return ""
}

func getItemConstructWhereQuery(filter entities.GetItemQueries) string {
	whereSQL := []string{}

	// Validate MerchantCategory using Ozzo validation

	err := validation.Validate(&filter.ProductCategory, validation.In("Beverage", "Food", "Snack", "Condiments", "Additions"))
	if err != nil {
		filter.ProductCategory = "" // Reset the category if it's invalid
	}

	whereSQL = append(whereSQL, " merchant_id = '"+filter.MerchantId+"'")

	if filter.ItemId != "" {
		whereSQL = append(whereSQL, " id = '"+filter.ItemId+"'")
	}

	if filter.ProductCategory != "" {
		whereSQL = append(whereSQL, " product_category = '"+filter.ProductCategory+"'")
	}

	if filter.Name != "" {
		whereSQL = append(whereSQL, " name ILIKE '%"+filter.Name+"%'")
	}

	if len(whereSQL) > 0 {
		return " WHERE " + strings.Join(whereSQL, " AND ")
	}

	return ""
}
