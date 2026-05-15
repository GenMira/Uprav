package gormrepository

import (
	"context"
	"uprav/model"

	//"gorm.io/gorm"

)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *model.Task) error
	GetAllTasks(ctx context.Context) ([]model.Task, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByName(ctx context.Context, name string) (*model.User, error)
}

// type Repository struct {
// 	db *gorm.DB
// }

// func NewGormRepository(db *gorm.DB) (*Repository, error) {
// 	repo := &Repository{db: db}
// 	// err := migration.Migrate(db)

// 	return repo, nil
// }

// func (r *Repository) Transaction(
// 	ctx context.Context,
// 	fn func(tx *Repository) error,
// ) error {
// 	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		return fn(&Repository{db: tx})
// 	})
// }