package svc

import (
	"beli_mang/db/entities"
	"beli_mang/repo"
	"context"

	"github.com/jackc/pgx/v5"
)

type PurchaseSvc interface {
	GetNearbyMerchant(ctx context.Context, getNearbyMerchantQueries entities.GetNearbyMerchantQueries) ([]entities.GetNearbyMerchantResponse, error)
}

type purchaseSvc struct {
	repo repo.PurchaseRepo
}

func NewPurchaseSvc(repo repo.PurchaseRepo) PurchaseSvc {
	return &purchaseSvc{repo}
}

func (s *purchaseSvc) GetNearbyMerchant(ctx context.Context, getNearbyMerchantQueries entities.GetNearbyMerchantQueries) ([]entities.GetNearbyMerchantResponse, error) {

	merchants, err := s.repo.GetNearbyMerchants(ctx, getNearbyMerchantQueries)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []entities.GetNearbyMerchantResponse{}, nil
		}
		return []entities.GetNearbyMerchantResponse{}, err
	}

	return merchants, nil
}
