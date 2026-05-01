package gormrepository

import (
	//"gorm.io/gorm"
	"uprav/model"
	"context"

)

func (r *Repository) CreateTask(ctx context.Context, task *model.Task) error {
	if err := r.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllTasks(ctx context.Context) ([]model.Task, error) {
    var tasks []model.Task
    if err := r.db.WithContext(ctx).Find(&tasks).Error;err != nil {
			return nil, err
		}
    return tasks, nil
}