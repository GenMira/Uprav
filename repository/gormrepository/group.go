package gormrepository

import (
	"gorm.io/gorm"
	"uprav/model"
	"context"
	"fmt"

)

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *groupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) CreateGroup(ctx context.Context, group *model.Group) error {
	if err := r.db.Create(group).Error; err != nil {
		fmt.Printf("[GORM ERROR] CreateGroup failed: %v\n", err)
		return err
	}
	return nil
}

func (r *groupRepository) GetGroups(ctx context.Context, userID int) ([]model.Group, error) {
	var groups []model.Group
  if err := r.db.WithContext(ctx).
		Preload("Members"). // レスポンス用にメンバー全員の情報をロードする
		Joins("JOIN group_members ON group_members.group_id = groups.id").
		Where("group_members.uid = ?", userID).
		Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}