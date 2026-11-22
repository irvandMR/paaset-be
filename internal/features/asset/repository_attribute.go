package asset

import (
	"context"

	"github.com/google/uuid"
)

func (ar *assetRepository) BulkInsertAssetAttributes(ctx context.Context, attrs []AssetAttributes) error {
	if len(attrs) == 0 {
		return nil
	}
	return ar.db.WithContext(ctx).Create(&attrs).Error
}

func (ar *assetRepository) DeleteAttributesByAssetId(ctx context.Context, assetId uuid.UUID) error {
	return ar.db.WithContext(ctx).Where("asset_id = ?", assetId).Delete(&AssetAttributes{}).Error
}

func (ar *assetRepository) FindAttributesByAssetId(ctx context.Context, assetId uuid.UUID) ([]AssetAttributes, error) {
	var results []AssetAttributes
	err := ar.db.WithContext(ctx).
		Table("tbl_asset_attribute_values aav").
		Select("aav.id, aav.asset_id, aav.attribute_id, ca.name, ca.data_type, ca.unit, aav.custome_label, aav.value_text, aav.value_number, aav.sort_order").
		Joins("LEFT JOIN tbl_catageries_attributes ca ON ca.id = aav.attribute_id").
		Where("aav.asset_id = ?", assetId).
		Order("aav.sort_order ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
