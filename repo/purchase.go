package repo

import (
	"beli_mang/db/entities"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseRepo interface {
	GetNearbyMerchants(ctx context.Context, filter entities.GetNearbyMerchantQueries) ([]entities.GetNearbyMerchantResponse, int, error)
	GetMerchantLocations(ctx context.Context, getEstimatePayload entities.GetEstimatePayload) ([]entities.RoutePoint, error)
	GetTotalItemsPrice(ctx context.Context, getEstimatePayload entities.GetEstimatePayload) (int, error)
	SaveOrderEstimation(ctx context.Context, order entities.OrderInfo) (string, error)
	SaveOrderItems(ctx context.Context, getEstimatePayload entities.GetEstimatePayload, orderId string) error
	PlaceOrder(ctx context.Context, placeOrderPayload entities.PlaceOrderPayload) (pgconn.CommandTag, error)
	GetUserOrders(ctx context.Context, filter entities.GetUserOrderQueries) ([]entities.GetUserOrderResponse, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type purchaseRepo struct {
	db *pgxpool.Pool
}

func NewPurchaseRepo(db *pgxpool.Pool) PurchaseRepo {
	return &purchaseRepo{db}
}

func (r *purchaseRepo) GetNearbyMerchants(ctx context.Context, filter entities.GetNearbyMerchantQueries) ([]entities.GetNearbyMerchantResponse, int, error) {
	var merchants []entities.GetNearbyMerchantResponse
	var totalCount int

	var distance float64
	query := `
    SELECT id, name, merchant_category, image_url, latitude, longitude, created_at, count(*) OVER() AS total_count,
        (acos(
            cos(radians($1)) * cos(radians(latitude)) *
            cos(radians(longitude) - radians($2)) +
            sin(radians($1)) * sin(radians(latitude))
        )) AS distance
    FROM merchants` //No need for earth constant because they all would be multiplied with the same constant

	query += getNearbyMerchantConstructWhereQuery(filter)

	query += " ORDER BY distance"
	query += " LIMIT $3 OFFSET $4"

	rows, err := r.db.Query(ctx, query, filter.Latitude, filter.Longitude, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		merchant := entities.GetNearbyMerchantResponse{}
		var createdAt time.Time
		err := rows.Scan(&merchant.Merchant.MerchantId, &merchant.Merchant.Name, &merchant.Merchant.MerchantCategory, &merchant.Merchant.ImageUrl, &merchant.Merchant.Location.Latitude, &merchant.Merchant.Location.Longitude, &createdAt, &totalCount, &distance)
		if err != nil {
			return nil, 0, err
		}
		merchant.Merchant.CreatedAt = createdAt.Format(time.RFC3339Nano)

		merchants = append(merchants, merchant)
	}

	getItemsQuery := "SELECT id, name, product_category, price, image_url, created_at FROM items WHERE merchant_id = $1"

	var wg sync.WaitGroup
	errChan := make(chan error, len(merchants))
	resultsChan := make(chan entities.GetNearbyMerchantResponse, len(merchants))

	for _, merchant := range merchants {
		wg.Add(1)
		go func(merchant entities.GetNearbyMerchantResponse) {
			defer wg.Done()

			itemRows, err := r.db.Query(ctx, getItemsQuery, merchant.Merchant.MerchantId)
			if err != nil {
				errChan <- err
				return
			}
			defer itemRows.Close()

			for itemRows.Next() {
				item := entities.GetItemResponse{}
				var createdAtItem time.Time

				err := itemRows.Scan(&item.ItemId, &item.Name, &item.ProductCategory, &item.Price, &item.ImageUrl, &createdAtItem)
				if err != nil {
					errChan <- err
					return
				}

				item.CreatedAt = createdAtItem.Format(time.RFC3339Nano)
				merchant.Items = append(merchant.Items, item)
			}

			resultsChan <- merchant
		}(merchant)
	}

	wg.Wait()
	close(errChan)
	close(resultsChan)

	for err := range errChan {
		if err != nil {
			return nil, 0, err
		}
	}

	var resultMerchants []entities.GetNearbyMerchantResponse
	for merchant := range resultsChan {
		resultMerchants = append(resultMerchants, merchant)
	}

	return resultMerchants, totalCount, nil
}

func (r *purchaseRepo) GetMerchantLocations(ctx context.Context, getEstimatePayload entities.GetEstimatePayload) ([]entities.RoutePoint, error) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(getEstimatePayload.Orders))
	startingPointsChan := make(chan entities.RoutePoint, len(getEstimatePayload.Orders))
	nonStartingPointsChan := make(chan entities.RoutePoint, len(getEstimatePayload.Orders))

	for _, order := range getEstimatePayload.Orders {
		wg.Add(1)
		go func(order entities.Order) {
			defer wg.Done()

			var merch entities.RoutePoint
			query := "SELECT latitude, longitude FROM merchants WHERE id = $1"

			row := r.db.QueryRow(ctx, query, order.MerchantId)
			err := row.Scan(&merch.Latitude, &merch.Longitude)
			if err != nil {
				errChan <- err
				return
			}

			if order.IsStartingPoint {
				startingPointsChan <- merch
			} else {
				nonStartingPointsChan <- merch
			}
		}(order)
	}

	wg.Wait()
	close(errChan)
	close(startingPointsChan)
	close(nonStartingPointsChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	var merchants []entities.RoutePoint
	for merch := range startingPointsChan {
		merchants = append([]entities.RoutePoint{merch}, merchants...)
	}
	for merch := range nonStartingPointsChan {
		merchants = append(merchants, merch)
	}

	return merchants, nil

}

func (r *purchaseRepo) GetTotalItemsPrice(ctx context.Context, getEstimatePayload entities.GetEstimatePayload) (int, error) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(getEstimatePayload.Orders))
	priceChan := make(chan int, len(getEstimatePayload.Orders))

	for _, order := range getEstimatePayload.Orders {
		for _, item := range order.Items {
			wg.Add(1)
			go func(order entities.Order, item entities.OrderItem) {
				defer wg.Done()

				var price int
				query := "SELECT price FROM items WHERE id = $1 AND merchant_id = $2"
				err := r.db.QueryRow(ctx, query, item.ItemId, order.MerchantId).Scan(&price)
				if err != nil {
					errChan <- err
					return
				}
				priceChan <- (price * item.Quantity)
			}(order, item)
		}
	}

	wg.Wait()
	close(errChan)
	close(priceChan)

	for err := range errChan {
		if err != nil {
			return 0, err
		}
	}

	var totalPrice int
	for price := range priceChan {
		totalPrice += price
	}

	return totalPrice, nil
}

func (r *purchaseRepo) SaveOrderEstimation(ctx context.Context, order entities.OrderInfo) (string, error) {
	var id string

	statement := "INSERT INTO orders (user_id, total_price, estimated_delivery_time_in_minutes, status) VALUES ($1, $2, $3, $4) RETURNING id"

	row := r.db.QueryRow(ctx, statement, order.UserId, order.TotalPrice, order.EstimatedDeliveryTimeInMinutes, order.Status)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *purchaseRepo) SaveOrderItems(ctx context.Context, getEstimatePayload entities.GetEstimatePayload, orderId string) error {
	statement := "INSERT INTO order_items (order_id, item_id, quantity) VALUES ($1, $2, $3)"
	var wg sync.WaitGroup
	errChan := make(chan error, len(getEstimatePayload.Orders))

	for _, order := range getEstimatePayload.Orders {
		for _, item := range order.Items {
			wg.Add(1)
			go func(orderId string, item entities.OrderItem) {
				defer wg.Done()

				_, err := r.db.Exec(ctx, statement, orderId, item.ItemId, item.Quantity)
				if err != nil {
					errChan <- err
				}
			}(orderId, item)
		}
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *purchaseRepo) PlaceOrder(ctx context.Context, placeOrderPayload entities.PlaceOrderPayload) (pgconn.CommandTag, error) {
	statement := "UPDATE orders SET status = 'ordered' WHERE id = $1"

	res, err := r.db.Exec(ctx, statement, placeOrderPayload.CalculatedEstimateId)

	return res, err
}

func (r *purchaseRepo) GetUserOrders(ctx context.Context, filter entities.GetUserOrderQueries) ([]entities.GetUserOrderResponse, error) {
	whereClause, args := getGetUserOrderConstructWhereQuery(filter)

	query := `
	SELECT
	    o.id as order_id,
	    m.id as merchant_id,
	    m.name as merchant_name,
	    m.image_url as merchant_image_url,
	    m.merchant_category,
	    m.created_at as merchant_created_at,
	    i.id as item_id,
	    i.name as item_name,
	    i.price as item_price,
	    i.image_url as item_image_url,
	    i.product_category,
	    oi.quantity,
	    i.created_at as item_created_at
	FROM
	    orders o
	JOIN
	    order_items oi ON o.id = oi.order_id
	JOIN
	    items i ON oi.item_id = i.id
	JOIN
	    merchants m ON i.merchant_id = m.id
	WHERE ` + whereClause + `
	LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderMap := make(map[string]entities.GetUserOrderResponse)
	for rows.Next() {
		var (
			orderId           string
			merchantId        string
			merchantName      string
			merchantImageUrl  string
			merchantCategory  string
			merchantCreatedAt time.Time
			itemId            string
			itemName          string
			itemPrice         int
			itemImageUrl      string
			productCategory   string
			quantity          int
			itemCreatedAt     time.Time
		)

		err := rows.Scan(
			&orderId,
			&merchantId,
			&merchantName,
			&merchantImageUrl,
			&merchantCategory,
			&merchantCreatedAt,
			&itemId,
			&itemName,
			&itemPrice,
			&itemImageUrl,
			&productCategory,
			&quantity,
			&itemCreatedAt,
		)
		if err != nil {
			return nil, err
		}

		item := entities.ItemResponse{
			ItemId:          itemId,
			Name:            itemName,
			Price:           itemPrice,
			Quantity:        quantity,
			ImageUrl:        itemImageUrl,
			ProductCategory: productCategory,
			CreatedAt:       itemCreatedAt.Format(time.RFC3339Nano),
		}

		if order, exists := orderMap[orderId]; exists {
			order.Orders[0].Items = append(order.Orders[0].Items, item)
			orderMap[orderId] = order
		} else {
			merchant := entities.GetMerchantResponse{
				MerchantId:       merchantId,
				Name:             merchantName,
				ImageUrl:         merchantImageUrl,
				MerchantCategory: merchantCategory,
				CreatedAt:        merchantCreatedAt.Format(time.RFC3339Nano),
			}
			orderMap[orderId] = entities.GetUserOrderResponse{
				OrderId: orderId,
				Orders: []entities.OrderResponse{
					{
						Merchant: merchant,
						Items:    []entities.ItemResponse{item},
					},
				},
			}
		}
	}

	var orders []entities.GetUserOrderResponse
	for _, order := range orderMap {
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *purchaseRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

func getNearbyMerchantConstructWhereQuery(filter entities.GetNearbyMerchantQueries) string {
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

func getGetUserOrderConstructWhereQuery(filter entities.GetUserOrderQueries) (string, []interface{}) {
	whereSQL := []string{"o.user_id = $1", "o.status = 'ordered'"}
	args := []interface{}{filter.UserId}
	argIdx := 2

	err := validation.Validate(&filter.MerchantCategory, validation.In("SmallRestaurant", "MediumRestaurant", "LargeRestaurant",
		"MerchandiseRestaurant", "BoothKiosk", "ConvenienceStore"))
	if err != nil {
		filter.MerchantCategory = "" // Reset the category if it's invalid
	}

	if filter.MerchantId != "" {
		whereSQL = append(whereSQL, fmt.Sprintf("m.id = $%d", argIdx))
		args = append(args, filter.MerchantId)
		argIdx++
	}

	if filter.MerchantCategory != "" {
		whereSQL = append(whereSQL, fmt.Sprintf("m.merchant_category = $%d", argIdx))
		args = append(args, filter.MerchantCategory)
		argIdx++
	}

	if filter.Name != "" {
		whereSQL = append(whereSQL, fmt.Sprintf("(m.name ILIKE $%d OR i.name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Name+"%")
		argIdx++
	}

	return strings.Join(whereSQL, " AND "), args
}
