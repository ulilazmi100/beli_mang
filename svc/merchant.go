package svc

import (
	"beli_mang/db/entities"
	"beli_mang/repo"
	"beli_mang/responses"
	"context"

	"github.com/jackc/pgx/v5"
)

type MerchantSvc interface {
	RegisterMerchant(ctx context.Context, newMerchant entities.MerchantRegistrationPayload) (string, error)
	GetMerchant(ctx context.Context, getMerchantQueries entities.GetMerchantQueries) ([]entities.GetMerchantResponse, int, error)
	RegisterItem(ctx context.Context, newItem entities.ItemRegistrationPayload) (string, error)
	GetItem(ctx context.Context, getItemQueries entities.GetItemQueries) ([]entities.GetItemResponse, int, error)
}

type merchantSvc struct {
	repo repo.MerchantRepo
}

func NewMerchantSvc(repo repo.MerchantRepo) MerchantSvc {
	return &merchantSvc{repo}
}

func (s *merchantSvc) RegisterMerchant(ctx context.Context, newMerchant entities.MerchantRegistrationPayload) (string, error) {
	if err := newMerchant.Validate(); err != nil {
		return "", responses.NewBadRequestError(err.Error())
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx) // Ensures rollback in case of an error

	id, err := s.repo.CreateMerchant(ctx, tx, &newMerchant)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return id, nil
}

func (s *merchantSvc) GetMerchant(ctx context.Context, getMerchantQueries entities.GetMerchantQueries) ([]entities.GetMerchantResponse, int, error) {
	merchants, totalCount, err := s.repo.GetMerchants(ctx, getMerchantQueries)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []entities.GetMerchantResponse{}, 0, nil
		}
		return nil, 0, err
	}

	return merchants, totalCount, nil
}

func (s *merchantSvc) RegisterItem(ctx context.Context, newItem entities.ItemRegistrationPayload) (string, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx) // Ensures rollback in case of an error

	_, err = s.repo.GetMerchant(ctx, tx, newItem.MerchantId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", responses.NewNotFoundError("merchantId not found")
		}
		return "", err
	}

	if err := newItem.Validate(); err != nil {
		return "", responses.NewBadRequestError(err.Error())
	}

	id, err := s.repo.CreateItem(ctx, tx, &newItem)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return id, nil
}

func (s *merchantSvc) GetItem(ctx context.Context, getItemQueries entities.GetItemQueries) ([]entities.GetItemResponse, int, error) {
	_, err := s.repo.GetMerchant(ctx, nil, getItemQueries.MerchantId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []entities.GetItemResponse{}, 0, responses.NewNotFoundError("merchantId not found")
		}
		return []entities.GetItemResponse{}, 0, err
	}

	items, totalCount, err := s.repo.GetItem(ctx, getItemQueries)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []entities.GetItemResponse{}, 0, nil
		}
		return nil, 0, err
	}

	return items, totalCount, nil
}
