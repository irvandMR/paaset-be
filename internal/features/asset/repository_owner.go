package asset

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ar *assetRepository) InsertAssetOwner(ctx context.Context, owner *AssetOwner) error {
	return ar.db.WithContext(ctx).Create(owner).Error
}

func (ar *assetRepository) UpdateAssetOwner(ctx context.Context, owner *AssetOwner) error {
	owner.UpdatedAt = time.Now()
	return ar.db.WithContext(ctx).
		Model(&AssetOwner{}).
		Where("asset_id = ?", owner.AssetId).
		Updates(map[string]any{
			"owner_name":     owner.OwnerName,
			"owner_contact":  owner.OwnerContact,
			"internal_notes": owner.InternalNotes,
			"updated_at":     owner.UpdatedAt,
		}).Error
}

func (ar *assetRepository) FindOwnerByAssetId(ctx context.Context, assetId uuid.UUID) (*AssetOwner, error) {
	var owner AssetOwner
	err := ar.db.WithContext(ctx).Where("asset_id = ?", assetId).First(&owner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &owner, nil
}
