package asset

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ar *assetRepository) InsertAssetMedia(ctx context.Context, media *AssetMedia) (*AssetMedia, error) {
	err := ar.db.WithContext(ctx).Create(media).Error
	if err != nil {
		return nil, err
	}
	return media, nil
}

func (ar *assetRepository) FindMediaByAssetId(ctx context.Context, assetId uuid.UUID) ([]AssetMedia, error) {
	var results []AssetMedia
	err := ar.db.WithContext(ctx).
		Where("asset_id = ?", assetId).
		Order("sort_order ASC").
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (ar *assetRepository) FindMediaById(ctx context.Context, id uuid.UUID) (*AssetMedia, error) {
	var media AssetMedia
	err := ar.db.WithContext(ctx).Where("id = ?", id).First(&media).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &media, nil
}

func (ar *assetRepository) DeleteMediaById(ctx context.Context, id uuid.UUID) error {
	return ar.db.WithContext(ctx).Where("id = ?", id).Delete(&AssetMedia{}).Error
}
