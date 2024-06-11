package repo

import (
	"beli_mang/db/entities"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo interface {
	GetUser(ctx context.Context, username string) (*entities.User, error)
	GetUserByUsernameOrMailAndRole(ctx context.Context, username string, email string, role string) (*entities.User, error)
	GetUserByUsernameOrMailAndRoleTx(ctx context.Context, tx pgx.Tx, username string, email string, role string) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.RegistrationPayload, hashPassword string, role string) (string, error)
	CreateUserTx(ctx context.Context, tx pgx.Tx, user *entities.RegistrationPayload, hashPassword string, role string) (string, error)
	GetUsernameById(ctx context.Context, id string) (*entities.User, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepo {
	return &userRepo{db}
}

func (r *userRepo) GetUser(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	query := "SELECT id, email, password, role FROM users WHERE username = $1"

	row := r.db.QueryRow(ctx, query, username)
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	user.Username = username

	return &user, nil
}

func (r *userRepo) GetUserByUsernameOrMailAndRole(ctx context.Context, username string, email string, role string) (*entities.User, error) {
	var user entities.User
	query := "SELECT id, email, password, role FROM users WHERE username = $1 OR (email = $2 AND role = $3)"

	row := r.db.QueryRow(ctx, query, username, email, role)
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	user.Username = username

	return &user, nil
}

func (r *userRepo) GetUserByUsernameOrMailAndRoleTx(ctx context.Context, tx pgx.Tx, username string, email string, role string) (*entities.User, error) {
	var user entities.User
	query := "SELECT id, email, password, role FROM users WHERE username = $1 OR (email = $2 AND role = $3)"

	row := tx.QueryRow(ctx, query, username, email, role)
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	user.Username = username

	return &user, nil
}

func (r *userRepo) CreateUser(ctx context.Context, user *entities.RegistrationPayload, hashPassword string, role string) (string, error) {
	var id string
	statement := "INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id"

	row := r.db.QueryRow(ctx, statement, user.Username, user.Email, hashPassword, role)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepo) CreateUserTx(ctx context.Context, tx pgx.Tx, user *entities.RegistrationPayload, hashPassword string, role string) (string, error) {
	var id string
	statement := "INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id"

	row := tx.QueryRow(ctx, statement, user.Username, user.Email, hashPassword, role)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepo) GetUsernameById(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	query := "SELECT username FROM users WHERE id = $1"

	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(&user.Username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}
