package gormrepository

import (
	"context"

	"gorm.io/gorm"

)

type Repository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) (*Repository, error) {
	repo := &Repository{db: db}
	// err := migration.Migrate(db)

	return repo, nil
}

func (r *Repository) Transaction(
	ctx context.Context,
	fn func(tx *Repository) error,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&Repository{db: tx})
	})
}