package repository

import (
	"context"
	"errors"
	"strings"

	"libratrack/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, u model.User) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
	FindByID(ctx context.Context, id uint) (model.User, error)
	List(ctx context.Context) ([]model.User, error)
	UpdateRole(ctx context.Context, id uint, role string) (model.User, error)
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(ctx context.Context, u model.User) (model.User, error) {
	if strings.TrimSpace(u.Email) == "" || u.PasswordHash == "" {
		return model.User{}, ErrValidation
	}
	if u.Role == "" {
		u.Role = "reader"
	}
	if err := r.db.WithContext(ctx).Create(&u).Error; err != nil {
		return model.User{}, mapDBErr(err)
	}
	return u, nil
}

func (r *gormUserRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, ErrNotFound
	}
	return u, err
}

func (r *gormUserRepository) FindByID(ctx context.Context, id uint) (model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, ErrNotFound
	}
	return u, err
}

func (r *gormUserRepository) List(ctx context.Context) ([]model.User, error) {
	var list []model.User
	err := r.db.WithContext(ctx).Order("id").Find(&list).Error
	return list, err
}

func (r *gormUserRepository) UpdateRole(ctx context.Context, id uint, role string) (model.User, error) {
	u, err := r.FindByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	u.Role = role
	if err := r.db.WithContext(ctx).Model(&u).Update("role", role).Error; err != nil {
		return model.User{}, err
	}
	return u, nil
}
