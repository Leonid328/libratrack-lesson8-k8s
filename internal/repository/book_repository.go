package repository

import (
	"context"
	"errors"
	"strings"

	"libratrack/internal/model"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type BookRepository interface {
	Create(ctx context.Context, in model.BookInput) (model.Book, error)
	FindByID(ctx context.Context, id uint) (model.Book, error)
	Update(ctx context.Context, id uint, in model.BookInput) (model.Book, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, categoryID uint, page, limit int) ([]model.Book, error)
}

type gormBookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &gormBookRepository{db: db}
}

func (r *gormBookRepository) Create(ctx context.Context, in model.BookInput) (model.Book, error) {
	if err := validateBookInput(in); err != nil {
		return model.Book{}, err
	}
	authors, err := r.loadAuthors(ctx, in.AuthorIDs)
	if err != nil {
		return model.Book{}, err
	}
	if err := r.ensureCategory(ctx, in.CategoryID); err != nil {
		return model.Book{}, err
	}
	book := model.Book{
		Title:      in.Title,
		ISBN:       in.ISBN,
		CategoryID: in.CategoryID,
		Authors:    authors,
	}
	err = r.db.WithContext(ctx).Create(&book).Error
	if err != nil {
		return model.Book{}, mapDBErr(err)
	}
	return r.FindByID(ctx, book.ID)
}

func (r *gormBookRepository) FindByID(ctx context.Context, id uint) (model.Book, error) {
	var book model.Book
	err := r.db.WithContext(ctx).Preload("Authors").Preload("Category").First(&book, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Book{}, ErrNotFound
	}
	return book, err
}

func (r *gormBookRepository) Update(ctx context.Context, id uint, in model.BookInput) (model.Book, error) {
	if err := validateBookInput(in); err != nil {
		return model.Book{}, err
	}
	book, err := r.FindByID(ctx, id)
	if err != nil {
		return model.Book{}, err
	}
	authors, err := r.loadAuthors(ctx, in.AuthorIDs)
	if err != nil {
		return model.Book{}, err
	}
	if err := r.ensureCategory(ctx, in.CategoryID); err != nil {
		return model.Book{}, err
	}
	book.Title = in.Title
	book.ISBN = in.ISBN
	book.CategoryID = in.CategoryID
	if err := r.db.WithContext(ctx).Model(&book).Updates(map[string]any{
		"title":       in.Title,
		"isbn":        in.ISBN,
		"category_id": in.CategoryID,
	}).Error; err != nil {
		return model.Book{}, mapDBErr(err)
	}
	if err := r.db.WithContext(ctx).Model(&book).Association("Authors").Replace(authors); err != nil {
		return model.Book{}, err
	}
	return r.FindByID(ctx, id)
}

func (r *gormBookRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&model.Book{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *gormBookRepository) List(ctx context.Context, categoryID uint, page, limit int) ([]model.Book, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	q := r.db.WithContext(ctx).Preload("Authors").Preload("Category").Model(&model.Book{})
	if categoryID > 0 {
		q = q.Where("category_id = ?", categoryID)
	}
	var books []model.Book
	err := q.Offset((page - 1) * limit).Limit(limit).Order("id").Find(&books).Error
	return books, err
}

func (r *gormBookRepository) loadAuthors(ctx context.Context, ids []uint) ([]model.Author, error) {
	if len(ids) == 0 {
		return nil, ErrValidation
	}
	var authors []model.Author
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&authors).Error; err != nil {
		return nil, err
	}
	if len(authors) != len(uniqueIDs(ids)) {
		return nil, ErrInvalidRef
	}
	return authors, nil
}

func (r *gormBookRepository) ensureCategory(ctx context.Context, id uint) error {
	var n int64
	if err := r.db.WithContext(ctx).Model(&model.Category{}).Where("id = ?", id).Count(&n).Error; err != nil {
		return err
	}
	if n == 0 {
		return ErrInvalidRef
	}
	return nil
}

func validateBookInput(in model.BookInput) error {
	if strings.TrimSpace(in.Title) == "" || strings.TrimSpace(in.ISBN) == "" || in.CategoryID == 0 {
		return ErrValidation
	}
	if len(in.AuthorIDs) == 0 {
		return ErrValidation
	}
	return nil
}

func uniqueIDs(ids []uint) []uint {
	seen := map[uint]struct{}{}
	out := make([]uint, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func mapDBErr(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return ErrConflict
	}
	if errors.As(err, &pqErr) && pqErr.Code == "23503" {
		return ErrInvalidRef
	}
	return err
}
