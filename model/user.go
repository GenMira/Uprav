package model

type User struct {
	Uid      string `gorm:"type:uuid;primaryKey" json:"uid"`
	Name     string `json:"name"`
	Password string `json:"password"`
}