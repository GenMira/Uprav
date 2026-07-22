package gormrepository

import (
	"gorm.io/gorm"
	"uprav/model"
	"context"
	"fmt"
	"log"
	"github.com/google/uuid"
)

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *groupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) CreateGroup(ctx context.Context, name string, membersID []uint) (*model.Group, error) {
	groupID := uuid.New()

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userMap := make(map[uint]string)
		if len(membersID) > 0 {
			var users []model.User
			if err := tx.Where("uid IN ?", membersID).Find(&users).Error; err != nil {
				return fmt.Errorf("failed to fetch users for group members: %w", err)
			}
			for _, u := range users {
				userMap[uint(u.UID)] = u.Name
			}
		}

		var newGroupMembers []model.GroupMember
		for _, uid := range membersID {
			newGroupMembers = append(newGroupMembers, model.GroupMember{
				GroupID: groupID,
				UID:     uid,
				Name:    userMap[uid], 
			})
		}

		newGroup := model.Group{
			ID:      groupID,
			Name:    name,
			Members: newGroupMembers,
		}

		if err := tx.Create(&newGroup).Error; err != nil {
			return fmt.Errorf("failed to create group: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var createdGroup model.Group
	if err := r.db.WithContext(ctx).
		Preload("Members").
		Where("id = ?", groupID).
		First(&createdGroup).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch created group: %w", err)
	}

	return &createdGroup, nil
}

func (r *groupRepository) GetGroups(ctx context.Context, userID int) ([]model.Group, error) {
	var groups []model.Group

	subQuery := r.db.Model(&model.GroupMember{}).Select("group_id").Where("uid = ?", userID)

	if err := r.db.WithContext(ctx).
		Preload("Members").
		Where("id IN (?)", subQuery).
		Find(&groups).Error; err != nil {
		return nil, err
	}
	
	log.Printf("[DEBUG] GetGroups found %d groups", len(groups))
	return groups, nil
}

func (r *groupRepository) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", groupID).Delete(&model.Group{}).Error; err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	return nil
}

func (r *groupRepository) UpdateGroup(ctx context.Context, groupID uuid.UUID, name string, membersID []uint) (*model.Group, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userMap := make(map[uint]string)
		if len(membersID) > 0 {
			var users []model.User
			if err := tx.Where("uid IN ?", membersID).Find(&users).Error; err != nil {
				return fmt.Errorf("failed to fetch users for group members: %w", err)
			}
			for _, u := range users {
				userMap[uint(u.UID)] = u.Name
			}
		}

		var newMembers []model.GroupMember
		for _, uid := range membersID {
			newMembers = append(newMembers, model.GroupMember{
				GroupID: groupID,
				UID:     uid,
				Name:    userMap[uid], 
			})
		}

		if err := tx.Model(&model.Group{}).Where("id = ?", groupID).Update("name", name).Error; err != nil {
			return fmt.Errorf("failed to update group name: %w", err)
		}

		if err := tx.Where("group_id = ?", groupID).Delete(&model.GroupMember{}).Error; err != nil {
			return fmt.Errorf("failed to delete old group members: %w", err)
		}

		if len(newMembers) > 0 {
			if err := tx.Create(&newMembers).Error; err != nil {
				return fmt.Errorf("failed to insert new group members: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var updatedGroup model.Group
	if err := r.db.WithContext(ctx).
		Preload("Members").
		Where("id = ?", groupID).
		First(&updatedGroup).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated group: %w", err)
	}

	return &updatedGroup, nil
}

func (r *groupRepository) GetGroup(ctx context.Context, groupID uuid.UUID) (*model.Group, error) {
	var group model.Group
	if err := r.db.WithContext(ctx).
		Preload("Members").
		Where("id = ?", groupID).
		First(&group).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch group: %w", err)
	}
	return &group, nil
}