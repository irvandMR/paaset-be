package setting

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SettingsRepository interface {
	InsertSettings(ctx context.Context, request Setting) error
	UpdateSettings(ctx context.Context, request Setting) (*Setting, error)
	FindById(ctx context.Context, id uuid.UUID) (*Setting, error)
	Find(ctx context.Context) (*Setting, error)
}

type settingsRepository struct {
	db *gorm.DB
}

func (sr *settingsRepository) InsertSettings(ctx context.Context, request Setting) error {
	return sr.db.WithContext(ctx).Create(&request).Error
}

func (sr *settingsRepository) UpdateSettings(ctx context.Context, request Setting) (*Setting, error) {
	dateTime := time.Now()
	request.UpdatedAt = dateTime

	err := sr.db.WithContext(ctx).
		Model(&Setting{}).
		Where("id = ?", request.Id).
		Updates(map[string]any{
			"site_name":            request.SiteName,
			"wa_number":            request.WaNumber,
			"wa_button_text":       request.WaButtonText,
			"wa_message_template": request.WaMessageTemplate,
			"updated_at":           dateTime,
		}).Error

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (sr *settingsRepository) FindById(ctx context.Context, id uuid.UUID) (*Setting, error) {
	var setting Setting
	err := sr.db.WithContext(ctx).Where("id = ?", id).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &setting, nil
}

func (sr *settingsRepository) Find(ctx context.Context) (*Setting, error) {
	var setting Setting
	err := sr.db.WithContext(ctx).Order("updated_at DESC").First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &setting, nil
}

func NewSettingsRepository(db *gorm.DB) SettingsRepository {
	return &settingsRepository{db: db}
}
