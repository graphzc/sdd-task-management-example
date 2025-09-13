package task

import (
	"context"
	"database/sql"
	"errors"

	"github.com/graphzc/sdd-task-management-example/internal/domain/entities"
	"github.com/graphzc/sdd-task-management-example/internal/domain/enums"
	"github.com/graphzc/sdd-task-management-example/internal/utils/timeutil"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, task *entities.Task) (string, error)
	FindByID(ctx context.Context, taskID string) (*entities.Task, error)
	FindByUserID(ctx context.Context, userID string) ([]entities.Task, error)
	UpdateByID(ctx context.Context, taskID string, title, description string, priority enums.TaskPriority) error
	UpdateStatusByID(ctx context.Context, taskID string, status enums.TaskStatus) error
	DeleteByID(ctx context.Context, taskID string) error
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

func (r *repository) Create(ctx context.Context, task *entities.Task) (string, error) {
	taskModel, err := FromTaskEntity(task)
	if err != nil {
		return "", err
	}

	query := `
		INSERT INTO tasks (id, user_id, title, description, priority, status, created_at, updated_at)
		VALUES (:id, :user_id, :title, :description, :priority, :status, :created_at, :updated_at)
	`
	result, err := r.db.NamedExecContext(ctx, query, taskModel)
	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		return "", errors.New("no rows affected")
	}

	return task.ID, nil
}

func (r *repository) FindByID(ctx context.Context, taskID string) (*entities.Task, error) {
	query := `
		SELECT 
			id, user_id, title, description, priority, status, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	var taskModel Model
	err := r.db.GetContext(ctx, &taskModel, query, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return taskModel.ToTaskEntity(), nil
}

func (r *repository) FindByUserID(ctx context.Context, userID string) ([]entities.Task, error) {
	query := `
		SELECT 
			id, user_id, title, description, priority, status, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var taskModels []Model
	err := r.db.SelectContext(ctx, &taskModels, query, userID)
	if err != nil {
		return nil, err
	}

	tasks := make([]entities.Task, len(taskModels))
	for i, model := range taskModels {
		tasks[i] = *model.ToTaskEntity()
	}

	return tasks, nil
}

func (r *repository) UpdateByID(ctx context.Context, taskID string, title, description string, priority enums.TaskPriority) error {
	query := `
		UPDATE tasks 
		SET title = $1, description = $2, priority = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query, title, description, priority, timeutil.BangkokNow(), taskID)
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

func (r *repository) UpdateStatusByID(ctx context.Context, taskID string, status enums.TaskStatus) error {
	query := `
		UPDATE tasks 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, status, timeutil.BangkokNow(), taskID)
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

func (r *repository) DeleteByID(ctx context.Context, taskID string) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, taskID)
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
