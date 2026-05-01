package model

import (
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Task struct{
	gorm.Model
	Uid uuid.UUID
	Name string
	Priority int
	Tag string
	Deadline time.Time
	Period time.Time
	Assumption []uuid.UUID
	Group uuid.UUID
	Assign uuid.UUID
}