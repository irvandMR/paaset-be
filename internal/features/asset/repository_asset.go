package asset

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ar *assetRepository) InsertAsset(ctx context.Context, asset *Asset) (*Asset, error) {
	err := ar.db.WithContext(ctx).Create(asset).Error
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (ar *assetRepository) FindById(ctx context.Context, id uuid.UUID) (*Asset, error) {
	var asset Asset
	err := ar.db.WithContext(ctx).Where("id = ?", id).First(&asset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &asset, nil
}

func (ar *assetRepository) FindBySlug(ctx context.Context, slug string) (*Asset, error) {
	var asset Asset
	err := ar.db.WithContext(ctx).Where("slug = ?", slug).First(&asset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &asset, nil
}

func (ar *assetRepository) FindAll(ctx context.Context, filter AssetFilter) ([]Asset, int64, error) {
	query := ar.db.WithContext(ctx).Model(&Asset{})

	// 1. Status filter (Default: TERSEDIA jika nil/kosong)
	if filter.Status != nil && *filter.Status != "" {
		if strings.ToUpper(*filter.Status) != "ALL" {
			query = query.Where("status = ?", *filter.Status)
		}
	} else {
		query = query.Where("status = ?", "TERSEDIA")
	}

	// 2. Category filter
	if filter.CategoryId != nil && *filter.CategoryId != "" {
		query = query.Where("category_id = ?", *filter.CategoryId)
	}

	// 3. Search query (title, location, code)
	if strings.TrimSpace(filter.Search) != "" {
		keyword := "%" + strings.TrimSpace(filter.Search) + "%"
		query = query.Where("title ILIKE ? OR location ILIKE ? OR code ILIKE ?", keyword, keyword, keyword)
	}

	// 4. Min Price filter
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}

	// 5. Max Price filter
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}

	// Hitung total data untuk pagination meta
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderBySQL := mapSortOrder(filter.Sort)

	// Ambil data aset sesuai limit, offset, dan opsi pengurutan
	var assets []Asset
	err := query.Order(orderBySQL).Limit(filter.Limit).Offset(filter.Offset).Find(&assets).Error
	if err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}

// mapSortOrder memetakan query param sort ke klausa ORDER BY yang aman (whitelist anti SQL Injection)
func mapSortOrder(sortParam string) string {
	switch strings.ToLower(strings.TrimSpace(sortParam)) {
	case "price_asc":
		return "price ASC, created_at DESC"
	case "price_desc":
		return "price DESC, created_at DESC"
	case "popular":
		return "view_count DESC NULLS LAST, created_at DESC"
	case "oldest":
		return "created_at ASC"
	case "newest":
		return "created_at DESC"
	default:
		return "created_at DESC"
	}
}

func (ar *assetRepository) UpdateAsset(ctx context.Context, asset *Asset) error {
	now := time.Now()
	asset.UpdatedAt = &now

	return ar.db.WithContext(ctx).
		Model(&Asset{}).
		Where("id = ?", asset.Id).
		Updates(map[string]any{
			"title":            asset.Title,
			"slug":             asset.Slug,
			"location":         asset.Location,
			"description":      asset.Description,
			"badge":            asset.Badge,
			"price_label":      asset.PriceLabel,
			"price":            asset.Price,
			"price_negotiable": asset.PriceNegotiable,
			"updated_at":       asset.UpdatedAt,
		}).Error
}

func (ar *assetRepository) UpdateAssetStatus(ctx context.Context, id uuid.UUID, status string) error {
	now := time.Now()
	return ar.db.WithContext(ctx).
		Model(&Asset{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     status,
			"updated_at": now,
		}).Error
}

func (ar *assetRepository) DeleteAsset(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return ar.db.WithContext(ctx).
		Model(&Asset{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     "DELETED",
			"updated_at": now,
		}).Error
}

func (ar *assetRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return ar.db.WithContext(ctx).
		Model(&Asset{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
