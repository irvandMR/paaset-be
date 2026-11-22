package asset

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetRepository interface {
	// Transaction
	Transaction(ctx context.Context, fn func(repo AssetRepository) error) error

	// Asset
	InsertAsset(ctx context.Context, asset *Asset) (*Asset, error)
	UpdateAsset(ctx context.Context, asset *Asset) error
	UpdateAssetStatus(ctx context.Context, id uuid.UUID, status string) error
	FindById(ctx context.Context, id uuid.UUID) (*Asset, error)
	FindBySlug(ctx context.Context, slug string) (*Asset, error)
	FindAll(ctx context.Context, filter AssetFilter) ([]Asset, int64, error)
	DeleteAsset(ctx context.Context, id uuid.UUID) error
	IncrementViewCount(ctx context.Context, id uuid.UUID) error

	// Asset Attribute
	BulkInsertAssetAttributes(ctx context.Context, attrs []AssetAttributes) error
	DeleteAttributesByAssetId(ctx context.Context, assetId uuid.UUID) error
	FindAttributesByAssetId(ctx context.Context, assetId uuid.UUID) ([]AssetAttributes, error)

	// Asset Media
	InsertAssetMedia(ctx context.Context, media *AssetMedia) (*AssetMedia, error)
	FindMediaByAssetId(ctx context.Context, assetId uuid.UUID) ([]AssetMedia, error)
	FindMediaById(ctx context.Context, id uuid.UUID) (*AssetMedia, error)
	DeleteMediaById(ctx context.Context, id uuid.UUID) error

	// Asset Owner
	InsertAssetOwner(ctx context.Context, owner *AssetOwner) error
	UpdateAssetOwner(ctx context.Context, owner *AssetOwner) error
	FindOwnerByAssetId(ctx context.Context, assetId uuid.UUID) (*AssetOwner, error)
}

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db: db}
}

func (ar *assetRepository) Transaction(ctx context.Context, fn func(repo AssetRepository) error) error {
	return ar.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&assetRepository{db: tx})
	})
}
