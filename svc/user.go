package svc

import (
	"beli_mang/crypto"
	"beli_mang/db/entities"
	"beli_mang/repo"
	"beli_mang/responses"
	"context"

	"github.com/jackc/pgx/v5"
)

type UserSvc interface {
	AdminRegister(ctx context.Context, newUser entities.RegistrationPayload) (string, error)
	AdminLogin(ctx context.Context, user entities.Credential) (string, error)
	UserRegister(ctx context.Context, newUser entities.RegistrationPayload) (string, error)
	UserLogin(ctx context.Context, creds entities.Credential) (string, error)
}

type userSvc struct {
	repo repo.UserRepo
}

func NewUserSvc(repo repo.UserRepo) UserSvc {
	return &userSvc{repo}
}

func (s *userSvc) AdminRegister(ctx context.Context, newUser entities.RegistrationPayload) (string, error) {
	if err := newUser.Validate(); err != nil {
		return "", responses.NewBadRequestError(err.Error())
	}

	role := "admin"

	existingUser, err := s.repo.GetUserByUsernameOrMailAndRole(ctx, newUser.Username, newUser.Email, role)
	if err != nil {
		if err != pgx.ErrNoRows {
			return "", err
		}
	}

	if existingUser != nil {
		return "", responses.NewConflictError("user already exists, please change username or email")
	}

	hashedPassword, err := crypto.GenerateHashedPassword(newUser.Password)
	if err != nil {
		return "", err
	}

	id, err := s.repo.CreateUser(ctx, &newUser, hashedPassword, role)
	if err != nil {
		return "", err
	}

	token, err := crypto.GenerateToken(id, newUser.Username, role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userSvc) AdminLogin(ctx context.Context, creds entities.Credential) (string, error) {
	if err := creds.Validate(); err != nil {
		return "", responses.NewBadRequestError(err.Error())
	}

	user, err := s.repo.GetUser(ctx, creds.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", responses.NewNotFoundError("user not found")
		}
		return "", err
	}

	err = crypto.VerifyPassword(creds.Password, user.Password)
	if err != nil {
		return "", responses.NewBadRequestError("wrong password!")
	}

	role := "admin"

	token, err := crypto.GenerateToken(user.Id, user.Username, role)
	if err != nil {
		return "", responses.NewBadRequestError(err.Error())
	}

	return token, nil
}

func (s *userSvc) UserRegister(ctx context.Context, newUser entities.RegistrationPayload) (string, error) {
	if err := newUser.Validate(); err != nil {
		return "", responses.NewBadRequestError(err.Error())
	}

	role := "user"

	existingUser, err := s.repo.GetUserByUsernameOrMailAndRole(ctx, newUser.Username, newUser.Email, role)
	if err != nil {
		if err != pgx.ErrNoRows {
			return "", err
		}
	}

	if existingUser != nil {
		return "", responses.NewConflictError("user already exists")
	}

	hashedPassword, err := crypto.GenerateHashedPassword(newUser.Password)
	if err != nil {
		return "", err
	}

	id, err := s.repo.CreateUser(ctx, &newUser, hashedPassword, role)
	if err != nil {
		return "", err
	}

	token, err := crypto.GenerateToken(id, newUser.Username, role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userSvc) UserLogin(ctx context.Context, creds entities.Credential) (string, error) {
	if err := creds.Validate(); err != nil {
		return "", responses.NewBadRequestError(err.Error())
	}

	user, err := s.repo.GetUser(ctx, creds.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", responses.NewNotFoundError("user not found")
		}
		return "", err
	}

	err = crypto.VerifyPassword(creds.Password, user.Password)
	if err != nil {
		return "", responses.NewBadRequestError("wrong password!")
	}

	role := "user"

	token, err := crypto.GenerateToken(user.Id, user.Username, role)
	if err != nil {
		return "", responses.NewBadRequestError(err.Error())
	}

	return token, nil
}
