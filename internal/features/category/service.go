package category

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tjus-be/internal/pkg/entity"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CategoryService struct {
	catRepo CategoryRepository
}

func NewCategoryService(categoryRepository CategoryRepository) *CategoryService {
	return &CategoryService{catRepo: categoryRepository}
}

func (c *CategoryService) GetAllCategories(ctx context.Context) ([]CategoryResponse, error) {
	categories, err := c.catRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]CategoryResponse, 0, len(categories))
	for _, cat := range categories {
		responses = append(responses, CategoryResponse{
			Id:        cat.Id,
			Name:      cat.Name,
			Slug:      cat.Slug,
			SortOrder: cat.SortOrder,
			IsActive:  cat.IsActive,
		})
	}
	return responses, nil
}

func (c *CategoryService) CreateCategory(ctx context.Context, req *CreatedCatRequest) (*CategoryResponse, error) {
	baseSlug := generateSlug(req.Name)
	slug := baseSlug

	existing, err := c.catRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("category already exists")
	}

	category, err := c.catRepo.InsertCategory(ctx, Category{
		Name:      req.Name,
		Slug:      slug,
		SortOrder: req.SortOrder,
		IsActive:  req.IsActive,
	})

	if err != nil {
		return nil, err
	}

	return &CategoryResponse{
		Id:        category.Id,
		Name:      category.Name,
		Slug:      category.Slug,
		SortOrder: category.SortOrder,
		IsActive:  category.IsActive,
	}, nil
}

func (c *CategoryService) CreatedCatagory(ctx context.Context, req *CreatedCatRequest) (*CategoryResponse, error) {
	return c.CreateCategory(ctx, req)
}

func (c *CategoryService) EditCategory(ctx context.Context, req *UpdateCatRequest) (*CategoryResponse, error) {
	// 1. Validasi apakah kategori ada
	existingCategory, err := c.validateCategory(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 2. Jika nama berubah, kita buat slug baru
	slug := existingCategory.Slug
	if req.Name != existingCategory.Name {
		baseSlug := generateSlug(req.Name)
		slug = baseSlug

		existingSlug, err := c.catRepo.FindBySlug(ctx, slug)
		if err != nil {
			return nil, err
		}

		if existingSlug != nil && existingSlug.Id != req.Id {
			slug = fmt.Sprintf("%s-%d", baseSlug, time.Now().Unix()%10000)
		}
	}

	// 3. Update category
	category, err := c.catRepo.UpdateCategory(ctx, Category{
		BaseEntity: entity.BaseEntity{Id: req.Id},
		Name:       req.Name,
		Slug:       slug,
		SortOrder:  req.SortOrder,
		IsActive:   req.IsActive,
	})

	if err != nil {
		return nil, err
	}

	return &CategoryResponse{
		Id:        category.Id,
		Name:      category.Name,
		Slug:      category.Slug,
		SortOrder: category.SortOrder,
		IsActive:  category.IsActive,
	}, nil
}

func (c *CategoryService) EditCatagory(ctx context.Context, req *UpdateCatRequest) (*CategoryResponse, error) {
	return c.EditCategory(ctx, req)
}

func (c *CategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	// 1. Validasi apakah kategori ada
	if _, err := c.validateCategory(ctx, id); err != nil {
		return err
	}

	// 2. Eksekusi hapus di repository
	return c.catRepo.DeleteCategory(ctx, id)
}

func (c *CategoryService) CreateAttribute(ctx context.Context, req *CreatedCatAttrRequest) ([]CategoryAttributeResponse, error) {
	// 1. Validasi apakah category_id valid dan ada di database
	if _, err := c.validateCategory(ctx, req.CategoryId); err != nil {
		return nil, err
	}

	// 2. Siapkan slice entity untuk batch insert
	attrEntities := make([]CategoryAttributes, 0, len(req.Attributes))
	for _, item := range req.Attributes {
		options := item.Options
		if options == nil {
			options = []string{}
		}
		attrEntities = append(attrEntities, CategoryAttributes{
			CategoryId: req.CategoryId,
			Name:       item.Name,
			DataType:   item.DataType,
			Unit:       item.Unit,
			Options:    pq.StringArray(options),
			IsRequired: item.IsRequired,
			ShowOnCard: item.ShowOnCard,
			IsActive:   item.IsActive,
			SortOrder:  item.SortOrder,
		})
	}

	// 3. Simpan data sekaligus ke database menggunakan batch insert
	savedAttrs, err := c.catRepo.BulkInsertAttributes(ctx, attrEntities)
	if err != nil {
		return nil, err
	}

	// 4. Susun response
	responses := make([]CategoryAttributeResponse, 0, len(savedAttrs))
	for _, a := range savedAttrs {
		responses = append(responses, CategoryAttributeResponse{
			Id:         a.Id,
			CategoryId: a.CategoryId,
			Name:       a.Name,
			DataType:   a.DataType,
			Unit:       a.Unit,
			Options:    []string(a.Options),
			IsRequired: a.IsRequired,
			ShowOnCard: a.ShowOnCard,
			IsActive:   a.IsActive,
			SortOrder:  a.SortOrder,
		})
	}

	return responses, nil
}

func (c *CategoryService) GetAttributesByCategoryId(ctx context.Context, categoryId uuid.UUID) ([]CategoryAttributeResponse, error) {
	// Validasi apakah kategori ada
	if _, err := c.validateCategory(ctx, categoryId); err != nil {
		return nil, err
	}

	attrs, err := c.catRepo.FindAttributesByCategoryId(ctx, categoryId)
	if err != nil {
		return nil, err
	}

	responses := make([]CategoryAttributeResponse, 0, len(attrs))
	for _, a := range attrs {
		responses = append(responses, CategoryAttributeResponse{
			Id:         a.Id,
			CategoryId: a.CategoryId,
			Name:       a.Name,
			DataType:   a.DataType,
			Unit:       a.Unit,
			Options:    []string(a.Options),
			IsRequired: a.IsRequired,
			ShowOnCard: a.ShowOnCard,
			IsActive:   a.IsActive,
			SortOrder:  a.SortOrder,
		})
	}

	return responses, nil
}

func (c *CategoryService) EditAttribute(ctx context.Context, req *UpdateCatAttrRequest) (*CategoryAttributeResponse, error) {
	// 1. Validasi apakah attribute ada
	existingAttr, err := c.catRepo.FindAttributeById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if existingAttr == nil {
		return nil, errors.New("attribute not found")
	}

	// 2. Validasi apakah category_id ada
	if _, err := c.validateCategory(ctx, req.CategoryId); err != nil {
		return nil, err
	}

	// 3. Update attribute
	attr, err := c.catRepo.UpdateAttribute(ctx, CategoryAttributes{
		BaseEntity: entity.BaseEntity{Id: req.Id},
		CategoryId: req.CategoryId,
		Name:       req.Name,
		DataType:   req.DataType,
		Unit:       req.Unit,
		Options:    pq.StringArray(req.Options),
		IsRequired: req.IsRequired,
		ShowOnCard: req.ShowOnCard,
		IsActive:   req.IsActive,
		SortOrder:  req.SortOrder,
	})
	if err != nil {
		return nil, err
	}

	return &CategoryAttributeResponse{
		Id:         attr.Id,
		CategoryId: attr.CategoryId,
		Name:       attr.Name,
		DataType:   attr.DataType,
		Unit:       attr.Unit,
		Options:    []string(attr.Options),
		IsRequired: attr.IsRequired,
		ShowOnCard: attr.ShowOnCard,
		IsActive:   attr.IsActive,
		SortOrder:  attr.SortOrder,
	}, nil
}

func (c *CategoryService) DeleteAttribute(ctx context.Context, id uuid.UUID) error {
	// 1. Validasi apakah attribute ada
	existingAttr, err := c.catRepo.FindAttributeById(ctx, id)
	if err != nil {
		return err
	}
	if existingAttr == nil {
		return errors.New("attribute not found")
	}

	// 2. Hapus attribute
	return c.catRepo.DeleteAttribute(ctx, id)
}

func (c *CategoryService) validateCategory(ctx context.Context, id uuid.UUID) (*Category, error) {
	category, err := c.catRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
