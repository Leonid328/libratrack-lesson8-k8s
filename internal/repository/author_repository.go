package repository

import (
	"context"
	"errors"
	"strings"

	"libratrack/internal/model"

	"gorm.io/gorm"
)

type AuthorRepository interface {
	Create(ctx context.Context, a model.Author) (model.Author, error)
	FindByID(ctx context.Context, id uint) (model.Author, error)
	List(ctx context.Context) ([]model.Author, error)
}

type gormAuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &gormAuthorRepository{db: db}
}

func (r *gormAuthorRepository) Create(ctx context.Context, a model.Author) (model.Author, error) {
	if strings.TrimSpace(a.FullName) == "" {
		return model.Author{}, ErrValidation
	}
	if err := r.db.WithContext(ctx).Create(&a).Error; err != nil {
		return model.Author{}, mapDBErr(err)
	}
	return a, nil
}

func (r *gormAuthorRepository) FindByID(ctx context.Context, id uint) (model.Author, error) {
	var a model.Author
	err := r.db.WithContext(ctx).First(&a, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Author{}, ErrNotFound
	}
	return a, err
}

func (r *gormAuthorRepository) List(ctx context.Context) ([]model.Author, error) {
	var list []model.Author
	err := r.db.WithContext(ctx).Order("id").Find(&list).Error
	return list, err
}
