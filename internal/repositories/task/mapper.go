package task

import (
	"github.com/google/uuid"
	"github.com/graphzc/sdd-task-management-example/internal/domain/entities"
	"github.com/graphzc/sdd-task-management-example/internal/domain/enums"
)

func FromTaskEntity(entity *entities.Task) (*Model, error) {
	if entity == nil {
		return nil, ErrNullTask
	}

	taskUUID, err := uuid.Parse(entity.ID)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(entity.UserID)
	if err != nil {
		return nil, err
	}

	return &Model{
		ID:          taskUUID,
		UserID:      userUUID,
		Title:       entity.Title,
		Description: entity.Description,
		Priority:    entity.Priority.Int(),
		Status:      entity.Status.String(),
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}, nil
}

func (m *Model) ToTaskEntity() *entities.Task {
	return &entities.Task{
		ID:          m.ID.String(),
		UserID:      m.UserID.String(),
		Title:       m.Title,
		Description: m.Description,
		Priority:    enums.TaskPriority(m.Priority),
		Status:      enums.TaskStatus(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
