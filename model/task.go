package model

import (
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Task struct{
	gorm.Model
	Uid int
	Name string
	Priority int
	Tag string
	Deadline time.Time
	Period time.Time
	IsEveryday bool
	//Assumption []uuid.UUID
	Group uuid.UUID
	Assign uuid.UUID
}