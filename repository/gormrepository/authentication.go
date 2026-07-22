package gormrepository

import (
	"gorm.io/gorm"

	"uprav/model"
	"context"
	"fmt"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		fmt.Printf("[GORM ERROR] CreateUser failed: %v\n", err)
		return err
	}
	return nil
}

func(r *userRepository) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("name = ?", name).First(&user).Error; err != nil {
		fmt.Printf("[GORM ERROR] GetUserByName failed: %v\n", err)
		return nil, err
	}
	return &user, nil
}

func(r *userRepository) GetUserByUID(ctx context.Context,uid int)(*model.User, error) {
	var user model.User
	if err := r.db.Where("uid = ?", uid).First(&user).Error; err != nil {
		fmt.Printf("[GORM ERROR] GetUserByName failed: %v\n", err)
		return nil, err
	}
	return &user, nil
}