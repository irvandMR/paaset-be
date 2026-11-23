package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseEntity struct {
	Id        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (b *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	if b.Id == uuid.Nil {
		b.Id = uuid.New()
	}
	return nil
}
