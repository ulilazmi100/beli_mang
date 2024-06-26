package svc

import (
	"beli_mang/db/entities"
	"beli_mang/repo"
	"beli_mang/responses"
	"beli_mang/utils"
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
)

type PurchaseSvc interface {
	GetNearbyMerchant(ctx context.Context, getNearbyMerchantQueries entities.GetNearbyMerchantQueries) ([]entities.GetNearbyMerchantResponse, int, error)
	GetOrderEstimation(ctx context.Context, getEstimatePayload entities.GetEstimatePayload, userId string) (entities.GetEstimateResponse, error)
	PlaceOrder(ctx context.Context, placeOrderPayload entities.PlaceOrderPayload) (entities.PlaceOrderResponse, error)
	GetOrder(ctx context.Context, getUserOrderQueries entities.GetUserOrderQueries) ([]entities.GetUserOrderResponse, error)
}

type purchaseSvc struct {
	repo repo.PurchaseRepo
}

func NewPurchaseSvc(repo repo.PurchaseRepo) PurchaseSvc {
	return &purchaseSvc{repo}
}

func (s *purchaseSvc) GetNearbyMerchant(ctx context.Context, getNearbyMerchantQueries entities.GetNearbyMerchantQueries) ([]entities.GetNearbyMerchantResponse, int, error) {

	merchants, totalCount, err := s.repo.GetNearbyMerchants(ctx, getNearbyMerchantQueries)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []entities.GetNearbyMerchantResponse{}, 0, nil
		}
		return []entities.GetNearbyMerchantResponse{}, 0, err
	}

	return merchants, totalCount, nil
}

func (s *purchaseSvc) GetOrderEstimation(ctx context.Context, getEstimatePayload entities.GetEstimatePayload, userId string) (entities.GetEstimateResponse, error) {
	if err := getEstimatePayload.Validate(); err != nil {
		if strings.Contains(err.Error(), "invalid UUID") {
			return entities.GetEstimateResponse{}, responses.NewNotFoundError(err.Error())
		}
		return entities.GetEstimateResponse{}, responses.NewBadRequestError(err.Error())
	}

	locations, err := s.repo.GetMerchantLocations(ctx, getEstimatePayload)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.GetEstimateResponse{}, responses.NewNotFoundError("merchantId not found")
		}
		return entities.GetEstimateResponse{}, err
	}

	locations = append(locations,
		entities.RoutePoint{
			Latitude:  getEstimatePayload.Location.Latitude,
			Longitude: getEstimatePayload.Location.Longitude,
		})

	isWithin3Km2, err := utils.IsWithin3Km2(locations)
	if err != nil {
		return entities.GetEstimateResponse{}, err
	}
	if !isWithin3Km2 {
		return entities.GetEstimateResponse{}, responses.NewBadRequestError("location isn't within 3km squared")
	}

	totalPrice, err := s.repo.GetTotalItemsPrice(ctx, getEstimatePayload)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.GetEstimateResponse{}, responses.NewNotFoundError("itemId not found")
		}
		return entities.GetEstimateResponse{}, err
	}

	_, distance, err := utils.NearestNeighborTSP(locations)
	if err != nil {
		return entities.GetEstimateResponse{}, err
	}

	estimatedDeliveryTime := utils.EstimatedDeliveryTimeInMinutes(distance)

	var getEstimateResponse entities.GetEstimateResponse

	getEstimateResponse.EstimatedDeliveryTimeInMinutes = estimatedDeliveryTime
	getEstimateResponse.TotalPrice = totalPrice

	orderId, err := s.repo.SaveOrderEstimation(ctx, entities.OrderInfo{
		UserId:                         userId,
		TotalPrice:                     totalPrice,
		EstimatedDeliveryTimeInMinutes: estimatedDeliveryTime,
		Status:                         "estimated",
	})
	if err != nil {
		return entities.GetEstimateResponse{}, err
	}

	err = s.repo.SaveOrderItems(ctx, getEstimatePayload, orderId)
	if err != nil {
		return entities.GetEstimateResponse{}, err
	}

	return entities.GetEstimateResponse{
		TotalPrice:                     totalPrice,
		EstimatedDeliveryTimeInMinutes: estimatedDeliveryTime,
		CalculatedEstimateId:           orderId,
	}, nil
}

func (s *purchaseSvc) PlaceOrder(ctx context.Context, placeOrderPayload entities.PlaceOrderPayload) (entities.PlaceOrderResponse, error) {
	if err := placeOrderPayload.Validate(); err != nil {
		return entities.PlaceOrderResponse{}, responses.NewBadRequestError(err.Error())
	}

	res, err := s.repo.PlaceOrder(ctx, placeOrderPayload)
	if res.RowsAffected() == 0 {
		return entities.PlaceOrderResponse{}, responses.NewNotFoundError(err.Error())
	}

	if err != nil {
		return entities.PlaceOrderResponse{}, responses.NewInternalServerError(err.Error())
	}

	return entities.PlaceOrderResponse{
		OrderId: placeOrderPayload.CalculatedEstimateId,
	}, nil
}

func (s *purchaseSvc) GetOrder(ctx context.Context, getUserOrderQueries entities.GetUserOrderQueries) ([]entities.GetUserOrderResponse, error) {

	orders, err := s.repo.GetUserOrders(ctx, getUserOrderQueries)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []entities.GetUserOrderResponse{}, nil
		}
		return []entities.GetUserOrderResponse{}, err
	}

	return orders, nil
}
