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

func (r *taskRepository) GetAllTasks(ctx context.Context, uid int) ([]model.Task, error) {
    var tasks []model.Task
    if err := r.db.WithContext(ctx).Where("uid = ?", uid).Find(&tasks).Error; err != nil {
      return nil, err
    }
    return tasks, nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int, uid int) error {
	//if err := r.db.WithContext(ctx).Delete(&model.Task{}, id).Error; err != nil {
	if err := r.db.WithContext(ctx).
		Model(&model.Task{}).
		Where("id = ? AND uid = ?", id, uid). // ★ 他人のタスクを消せないように防御
		Delete(&model.Task{}).
		Error; err != nil {
		return err
	}
	return nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *model.Task) error {
	// 指定されたIDのレコードを探し、渡された構造体のデータで上書きします。
	// Updatesに構造体のポインタを渡すと、GORMは「ゼロ値（0や空文字）」以外のフィールドのみを自動で更新してくれます。
	if err := r.db.WithContext(ctx).
		Model(&model.Task{}).
		Where("id = ? AND uid = ?", task.ID, task.Uid). // 本人のタスクのみ上書きできるようUIDも条件に入れると安全です
		Updates(task).
		Error; err != nil {
		return err
	}	
	return nil
}