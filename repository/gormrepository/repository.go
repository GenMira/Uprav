package gormrepository

import (
	"context"
	"uprav/model"

	"github.com/google/uuid"
	//"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *model.Task) error
	GetAllTasks(ctx context.Context, uid int) ([]model.Task, error)
	DeleteTask(ctx context.Context, id int, uid int) error
	UpdateTask(ctx context.Context, task *model.Task) error
	GetTask(ctx context.Context, id int, uid int) (*model.Task, error)
	GetTags(ctx context.Context, uid int) ([]string, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByName(ctx context.Context, name string) (*model.User, error)
	GetUserByUID(ctx context.Context,uid int)(*model.User, error)
}

type GroupRepository interface {
	CreateGroup(ctx context.Context, name string, membersID []uint) (*model.Group, error)
	GetGroup(ctx context.Context, groupID uuid.UUID) (*model.Group, error)
	UpdateGroup(ctx context.Context, groupID uuid.UUID, name string, membersID []uint) (*model.Group, error)
	DeleteGroup(ctx context.Context, groupID uuid.UUID) error
	GetGroups(ctx context.Context, userID int) ([]model.Group, error)
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