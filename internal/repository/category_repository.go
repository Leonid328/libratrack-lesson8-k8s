package repository

import (
	"context"
	"errors"
	"strings"

	"libratrack/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, c model.Category) (model.Category, error)
	FindByID(ctx context.Context, id uint) (model.Category, error)
	List(ctx context.Context) ([]model.Category, error)
}

type gormCategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &gormCategoryRepository{db: db}
}

func (r *gormCategoryRepository) Create(ctx context.Context, c model.Category) (model.Category, error) {
	if strings.TrimSpace(c.Name) == "" {
		return model.Category{}, ErrValidation
	}
	if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
		return model.Category{}, mapDBErr(err)
	}
	return c, nil
}

func (r *gormCategoryRepository) FindByID(ctx context.Context, id uint) (model.Category, error) {
	var c model.Category
	err := r.db.WithContext(ctx).First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Category{}, ErrNotFound
	}
	return c, err
}

func (r *gormCategoryRepository) List(ctx context.Context) ([]model.Category, error) {
	var list []model.Category
	err := r.db.WithContext(ctx).Order("id").Find(&list).Error
	return list, err
}
