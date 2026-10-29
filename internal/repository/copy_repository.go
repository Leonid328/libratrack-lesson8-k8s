package repository

import (
	"context"
	"errors"
	"strings"

	"libratrack/internal/model"

	"gorm.io/gorm"
)

type CopyRepository interface {
	Create(ctx context.Context, bookID uint, inventoryNo, status string) (model.Copy, error)
	ListByBook(ctx context.Context, bookID uint) ([]model.Copy, error)
	Availability(ctx context.Context, bookID uint) (model.Availability, error)
}

type gormCopyRepository struct {
	db *gorm.DB
}

func NewCopyRepository(db *gorm.DB) CopyRepository {
	return &gormCopyRepository{db: db}
}

func (r *gormCopyRepository) Create(ctx context.Context, bookID uint, inventoryNo, status string) (model.Copy, error) {
	if strings.TrimSpace(inventoryNo) == "" {
		return model.Copy{}, ErrValidation
	}
	if status == "" {
		status = "available"
	}
	if status != "available" && status != "borrowed" && status != "lost" {
		return model.Copy{}, ErrValidation
	}
	var n int64
	if err := r.db.WithContext(ctx).Model(&model.Book{}).Where("id = ?", bookID).Count(&n).Error; err != nil {
		return model.Copy{}, err
	}
	if n == 0 {
		return model.Copy{}, ErrNotFound
	}
	cp := model.Copy{BookID: bookID, InventoryNo: inventoryNo, Status: status}
	if err := r.db.WithContext(ctx).Create(&cp).Error; err != nil {
		return model.Copy{}, mapDBErr(err)
	}
	return cp, nil
}

func (r *gormCopyRepository) ListByBook(ctx context.Context, bookID uint) ([]model.Copy, error) {
	if _, err := NewBookRepository(r.db).FindByID(ctx, bookID); err != nil {
		return nil, err
	}
	var list []model.Copy
	err := r.db.WithContext(ctx).Where("book_id = ?", bookID).Order("id").Find(&list).Error
	return list, err
}

func (r *gormCopyRepository) Availability(ctx context.Context, bookID uint) (model.Availability, error) {
	if _, err := NewBookRepository(r.db).FindByID(ctx, bookID); err != nil {
		return model.Availability{}, err
	}
	var total int64
	var avail int64
	r.db.WithContext(ctx).Model(&model.Copy{}).Where("book_id = ?", bookID).Count(&total)
	r.db.WithContext(ctx).Model(&model.Copy{}).Where("book_id = ? AND status = ?", bookID, "available").Count(&avail)
	return model.Availability{BookID: bookID, TotalCopies: int(total), AvailableCopies: int(avail)}, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
