package category

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	InsertCategory(ctx context.Context, request Category) (*Category, error)
	FindBySlug(ctx context.Context, slug string) (*Category, error)
	FindById(ctx context.Context, id uuid.UUID) (*Category, error)
	UpdateCategory(ctx context.Context, request Category) (*Category, error)
	FindAll(ctx context.Context) ([]Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error

	// Category Attributes
	InsertAttribute(ctx context.Context, attr CategoryAttributes) (*CategoryAttributes, error)
	BulkInsertAttributes(ctx context.Context, attrs []CategoryAttributes) ([]CategoryAttributes, error)
	FindAttributesByCategoryId(ctx context.Context, categoryId uuid.UUID) ([]CategoryAttributes, error)
	FindAttributeById(ctx context.Context, id uuid.UUID) (*CategoryAttributes, error)
	UpdateAttribute(ctx context.Context, attr CategoryAttributes) (*CategoryAttributes, error)
	DeleteAttribute(ctx context.Context, id uuid.UUID) error
}

type categoryRepository struct {
	db *gorm.DB
}

// FindAll implements [CategoryRepository].
func (c *categoryRepository) FindAll(ctx context.Context) ([]Category, error) {
	var categories []Category
	err := c.db.WithContext(ctx).
		Order("sort_order ASC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// UpdateCategory implements [CategoryRepository].
func (c *categoryRepository) UpdateCategory(ctx context.Context, request Category) (*Category, error) {
	now := time.Now()
	request.UpdatedAt = &now

	err := c.db.WithContext(ctx).
		Model(&Category{}).
		Where("id = ?", request.Id).
		Updates(map[string]any{
			"name":       request.Name,
			"slug":       request.Slug,
			"sort_order": request.SortOrder,
			"is_active":  request.IsActive,
			"updated_at": request.UpdatedAt,
		}).Error

	if err != nil {
		return nil, err
	}
	return &request, nil
}

// FindById implements [CategoryRepository].
func (c *categoryRepository) FindById(ctx context.Context, id uuid.UUID) (*Category, error) {
	var category Category
	err := c.db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// FindBySlug implements [CategoryRepository].
func (c *categoryRepository) FindBySlug(ctx context.Context, slug string) (*Category, error) {
	var category Category
	err := c.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// InsertCategory implements [CategoryRepository].
func (c *categoryRepository) InsertCategory(ctx context.Context, request Category) (*Category, error) {
	err := c.db.WithContext(ctx).Create(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// DeleteCategory implements [CategoryRepository].
func (c *categoryRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return c.db.WithContext(ctx).Where("id = ?", id).Delete(&Category{}).Error
}

// InsertAttribute implements [CategoryRepository].
func (c *categoryRepository) InsertAttribute(ctx context.Context, attr CategoryAttributes) (*CategoryAttributes, error) {
	err := c.db.WithContext(ctx).Create(&attr).Error
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

// BulkInsertAttributes menyimpan banyak atribut kategori sekaligus menggunakan GORM batch create
func (c *categoryRepository) BulkInsertAttributes(ctx context.Context, attrs []CategoryAttributes) ([]CategoryAttributes, error) {
	if len(attrs) == 0 {
		return []CategoryAttributes{}, nil
	}

	err := c.db.WithContext(ctx).Create(&attrs).Error
	if err != nil {
		return nil, err
	}

	return attrs, nil
}

// FindAttributesByCategoryId implements [CategoryRepository].
func (c *categoryRepository) FindAttributesByCategoryId(ctx context.Context, categoryId uuid.UUID) ([]CategoryAttributes, error) {
	var attributes []CategoryAttributes
	err := c.db.WithContext(ctx).
		Where("category_id = ?", categoryId).
		Order("sort_order ASC").
		Find(&attributes).Error
	if err != nil {
		return nil, err
	}
	return attributes, nil
}

// FindAttributeById implements [CategoryRepository].
func (c *categoryRepository) FindAttributeById(ctx context.Context, id uuid.UUID) (*CategoryAttributes, error) {
	var attr CategoryAttributes
	err := c.db.WithContext(ctx).Where("id = ?", id).First(&attr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &attr, nil
}

// UpdateAttribute implements [CategoryRepository].
func (c *categoryRepository) UpdateAttribute(ctx context.Context, attr CategoryAttributes) (*CategoryAttributes, error) {
	now := time.Now()
	attr.UpdatedAt = &now

	err := c.db.WithContext(ctx).
		Model(&CategoryAttributes{}).
		Where("id = ?", attr.Id).
		Updates(map[string]any{
			"category_id":  attr.CategoryId,
			"name":         attr.Name,
			"data_type":    attr.DataType,
			"unit":         attr.Unit,
			"options":      attr.Options,
			"is_required":  attr.IsRequired,
			"show_on_card": attr.ShowOnCard,
			"is_active":    attr.IsActive,
			"sort_order":   attr.SortOrder,
			"updated_at":   attr.UpdatedAt,
		}).Error

	if err != nil {
		return nil, err
	}
	return &attr, nil
}

// DeleteAttribute implements [CategoryRepository].
func (c *categoryRepository) DeleteAttribute(ctx context.Context, id uuid.UUID) error {
	return c.db.WithContext(ctx).Where("id = ?", id).Delete(&CategoryAttributes{}).Error
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}
