package setting

import (
	"time"

	"github.com/google/uuid"
)

type Setting struct {
	Id                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	SiteName          string    `gorm:"column:site_name;not null"`
	WaNumber          string    `gorm:"column:wa_number;not null"`
	WaButtonText      string    `gorm:"column:wa_button_text;not null"`
	WaMessageTemplate string    `gorm:"column:wa_message_template;not null"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Setting) TableName() string {
	return "tbl_settings"
}

type SettingsRequest struct {
	SiteName          string `json:"site_name"`
	WaNumber          string `json:"wa_number"`
	WaButtonText      string `json:"wa_button_text"`
	WaMessageTemplate string `json:"wa_message_template"`
}

type UpdateSettingsRequest struct {
	Id uuid.UUID `json:"id"`
	SettingsRequest
}

type SettingsResponse struct {
	Id                uuid.UUID  `json:"id"`
	SiteName          string     `json:"site_name"`
	WaNumber          string     `json:"wa_number"`
	WaButtonText      string     `json:"wa_button_text"`
	WaMessageTemplate string     `json:"wa_message_template"`
	UpdatedAt         *time.Time `json:"updated_at"`
}
