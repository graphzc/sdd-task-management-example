package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/graphzc/sdd-task-management-example/internal/domain/entities"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
}

type repository struct {
	db *sqlx.DB
}

// @WireSet("Repository")
func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, user *entities.User) error {
	userModel, err := FromUserEntity(user)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (id, name, email, password, created_at, updated_at)
		VALUES (:id, :name, :email, :password, :created_at, :updated_at)
	`
	result, err := r.db.NamedExecContext(ctx, query, userModel)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT 
			id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var userModel Model
	err := r.db.GetContext(ctx, &userModel, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return userModel.ToUserEntity(), nil
}
