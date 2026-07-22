package model

import "github.com/google/uuid"

type GroupMember struct {
	ID      uint   `gorm:"primaryKey"` // GORM用の主キー（自動インクリメント）
	GroupID uuid.UUID `gorm:"index"`       // 親の Group.ID と紐付けるための外部キー
	UID     uint   `json:"uid"`
	Name    string `json:"name"`
}

type Group struct {
	ID      uuid.UUID   `gorm:"primaryKey"` // 主キーであることを明記
	Name    string        `json:"name"`
	Members []GroupMember `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE;" json:"members"`
}