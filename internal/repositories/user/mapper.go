package user

import (
	"github.com/google/uuid"
	"github.com/graphzc/sdd-task-management-example/internal/domain/entities"
)

func FromUserEntity(entity *entities.User) (*Model, error) {
	if entity == nil {
		return nil, ErrNullUser
	}

	userUUID, err := uuid.Parse(entity.ID)
	if err != nil {
		return nil, err
	}

	return &Model{
		ID:        userUUID,
		Name:      entity.Name,
		Email:     entity.Email,
		Password:  entity.Password,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}

func (m *Model) ToUserEntity() *entities.User {
	return &entities.User{
		ID:        m.ID.String(),
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
