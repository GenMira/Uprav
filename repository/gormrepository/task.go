package gormrepository

import (
	"gorm.io/gorm"
	"uprav/model"
	"context"
	"fmt"

)
type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *taskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	if err := r.db.Create(task).Error; err != nil {
		fmt.Printf("[GORM ERROR] CreateTask failed: %v\n", err)
		return err
	}
	return nil
}

func (r *taskRepository) GetAllTasks(ctx context.Context) ([]model.Task, error) {
    var tasks []model.Task
    if err := r.db.WithContext(ctx).Find(&tasks).Error;err != nil {
			return nil, err
		}
    return tasks, nil
}