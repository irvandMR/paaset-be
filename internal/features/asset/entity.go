package asset

import (
	"time"
	"tjus-be/internal/pkg/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Entity

type Asset struct {
	entity.BaseEntity
	Code            string    `json:"code" gorm:"column:code;not null;uniqueIndex"`
	Slug            string    `json:"slug" gorm:"column:slug;not null;uniqueIndex"`
	CategoryId      uuid.UUID `json:"category_id" gorm:"column:category_id"`
	Title           string    `json:"title" gorm:"column:title;not null"`
	Location        string    `json:"location" gorm:"column:location;not null"`
	Description     string    `json:"description" gorm:"column:description"`
	Badge           string    `json:"badge" gorm:"column:badge"`
	PriceLabel      string    `json:"price_label" gorm:"column:price_label"`
	Price           int64     `json:"price" gorm:"column:price;not null"`
	PriceNegotiable bool      `json:"price_negotiable" gorm:"column:price_negotiable"`
	ViewCount       int32     `json:"view_count" gorm:"column:view_count;default:0"`
	Status          string    `json:"status" gorm:"column:status"`
}

func (Asset) TableName() string {
	return "tbl_assets"
}

type AssetAttributes struct {
	Id                  uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	AssetId             uuid.UUID `json:"asset_id" gorm:"column:asset_id"`
	CategoryAttributeId uuid.UUID `json:"attribute_id" gorm:"column:attribute_id"`
	Name                *string   `json:"name,omitempty" gorm:"->"`
	DataType            *string   `json:"data_type,omitempty" gorm:"->"`
	Unit                *string   `json:"unit,omitempty" gorm:"->"`
	CustomeLabel        string    `json:"custome_label" gorm:"column:custome_label"`
	ValueText           string    `json:"value_text" gorm:"column:value_text"`
	ValueNumber         int64     `json:"value_number" gorm:"column:value_number"`
	SortOrder           int       `json:"sort_order" gorm:"column:sort_order"`
}

func (AssetAttributes) TableName() string {
	return "tbl_asset_attribute_values"
}

func (a *AssetAttributes) BeforeCreate(tx *gorm.DB) error {
	if a.Id == uuid.Nil {
		a.Id = uuid.New()
	}
	return nil
}

type AssetMedia struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	AssetId   uuid.UUID `json:"asset_id" gorm:"column:asset_id"`
	TypeMedia string    `json:"type_media" gorm:"column:type_media"`
	FilePath  string    `json:"file_path" gorm:"column:file_path"`
	SortOrder int       `json:"sort_order" gorm:"column:sort_order"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (AssetMedia) TableName() string {
	return "tbl_asset_media"
}

func (m *AssetMedia) BeforeCreate(tx *gorm.DB) error {
	if m.Id == uuid.Nil {
		m.Id = uuid.New()
	}
	return nil
}

type AssetOwner struct {
	AssetId       uuid.UUID `json:"asset_id" gorm:"column:asset_id;primaryKey"`
	OwnerName     string    `json:"owner_name" gorm:"column:owner_name;not null"`
	OwnerContact  string    `json:"owner_contact" gorm:"column:owner_contact;not null"`
	InternalNotes string    `json:"internal_notes" gorm:"column:internal_notes"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (AssetOwner) TableName() string {
	return "tbl_asset_owner_info"
}

// Request

type AssetRequest struct {
	CategoryId      uuid.UUID          `json:"category_id" binding:"required"`
	Title           string             `json:"title" binding:"required"`
	Location        string             `json:"location" binding:"required"`
	Badge           string             `json:"badge"`
	PriceLabel      string             `json:"price_label"`
	Price           int64              `json:"price" binding:"required"`
	PriceNegotiable bool               `json:"price_negotiable"`
	Owner           AssetOwnerRequest  `json:"owner"`
	Attributes      []AttributeRequest `json:"attributes"`
	Description     string             `json:"description"`
}

type AttributeRequest struct {
	CategoryAttributeId uuid.UUID `json:"attribute_id"`
	CustomeLabel        string    `json:"custome_label"`
	ValueText           string    `json:"value_text"`
	ValueNumber         int64     `json:"value_number"`
	SortOrder           int       `json:"sort_order"`
}

type AssetOwnerRequest struct {
	OwnerName     string `json:"owner_name"`
	OwnerContact  string `json:"owner_contact"`
	InternalNotes string `json:"internal_notes"`
}

type AssetMediaRequest struct {
	AssetId   uuid.UUID `json:"asset_id"`
	TypeMedia string    `json:"type_media"`
	FilePath  string    `json:"file_path"`
	SortOrder int       `json:"sort_order"`
}

type UpdateAssetRequest struct {
	Id              uuid.UUID `json:"id" binding:"required"`
	GeneralUpdate   bool      `json:"general_update"`
	OwnerUpdate     bool      `json:"owner_update"`
	AttributeUpdate bool      `json:"attribute_update"`
	StatusUpdate    bool      `json:"status_update"`

	// Fields for GeneralUpdate
	Title           string `json:"title"`
	Location        string `json:"location"`
	Badge           string `json:"badge"`
	PriceLabel      string `json:"price_label"`
	Price           int64  `json:"price"`
	PriceNegotiable bool   `json:"price_negotiable"`
	Description     string `json:"description"`

	// Field for StatusUpdate
	Status string `json:"status"`

	// Field for OwnerUpdate
	Owner AssetOwnerRequest `json:"owner"`

	// Field for AttributeUpdate
	Attributes []AttributeRequest `json:"attributes"`
}

type AssetFilter struct {
	CategoryId *string `form:"category_id"`
	Search     string  `form:"q"`
	MinPrice   *int64  `form:"min_price"`
	MaxPrice   *int64  `form:"max_price"`
	Status     *string `form:"status"`
	Sort       string  `form:"sort"`
	Page       int     `form:"page"`
	Limit      int     `form:"limit"`
	Offset     int     `form:"-"`
}

// Response

type AssetOwnerResponse struct {
	OwnerName     string `json:"owner_name"`
	OwnerContact  string `json:"owner_contact"`
	InternalNotes string `json:"internal_notes"`
}

type AssetAttributeResponse struct {
	CategoryAttributeId uuid.UUID `json:"category_attribute_id"`
	Name                *string   `json:"name,omitempty"`
	DataType            *string   `json:"data_type,omitempty"`
	Unit                *string   `json:"unit,omitempty"`
	CustomeLabel        string    `json:"custome_label"`
	ValueText           string    `json:"value_text"`
	ValueNumber         int64     `json:"value_number"`
	SortOrder           int       `json:"sort_order"`
}

type AssetMediaResponse struct {
	Id        uuid.UUID `json:"id"`
	TypeMedia string    `json:"type_media"`
	FilePath  string    `json:"file_path"`
	SortOrder int       `json:"sort_order"`
}

type AssetResponse struct {
	Id              uuid.UUID                `json:"id"`
	Code            string                   `json:"code"`
	Slug            string                   `json:"slug"`
	CategoryId      uuid.UUID                `json:"category_id"`
	Title           string                   `json:"title"`
	Location        string                   `json:"location"`
	Description     string                   `json:"description"`
	Badge           string                   `json:"badge"`
	PriceLabel      string                   `json:"price_label"`
	Price           int64                    `json:"price"`
	PriceNegotiable bool                     `json:"price_negotiable"`
	ViewCount       int32                    `json:"view_count"`
	Status          string                   `json:"status"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       *time.Time               `json:"updated_at,omitempty"`
	Owner           *AssetOwnerResponse      `json:"owner,omitempty"`
	Attributes      []AssetAttributeResponse `json:"attributes,omitempty"`
	Media           []AssetMediaResponse     `json:"media,omitempty"`
}
