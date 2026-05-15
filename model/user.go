package model

// type User struct {
// 	Uid      string `gorm:"type:uuid;primaryKey" json:"uid"`
// 	Name     string `json:"name"`
// 	Password string `json:"password"`
// }

type User struct {
    ID   uint   `gorm:"primaryKey"`
    UID  int    `gorm:"uniqueIndex;autoIncrement:10000;not null"`
    Name string `gorm:"uniqueIndex;not null"`
    Password string `gorm:"not null"`
}
