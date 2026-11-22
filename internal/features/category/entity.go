package category

import (
	"tjus-be/internal/pkg/entity"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Category struct {
	entity.BaseEntity
	Name      string `json:"name" gorm:"column:name;not null"`
	Slug      string `json:"slug" gorm:"column:slug;not null"`
	SortOrder int    `json:"sort_order" gorm:"column:sort_order;default:0"`
	IsActive  bool   `json:"is_active" gorm:"column:is_active;default:true"`
}

func (Category) TableName() string {
	return "tbl_categories"
}

type CategoryAttributes struct {
	entity.BaseEntity
	CategoryId uuid.UUID      `json:"category_id" gorm:"column:category_id"`
	Name       string         `json:"name" gorm:"column:name"`
	DataType   string         `json:"data_type" gorm:"column:data_type"`
	Unit       string         `json:"unit" gorm:"column:unit"`
	Options    pq.StringArray `json:"options" gorm:"column:options;type:text[]"`
	IsRequired bool           `json:"is_required" gorm:"column:is_required"`
	ShowOnCard bool           `json:"show_on_card" gorm:"column:show_on_card"`
	IsActive   bool           `json:"is_active" gorm:"column:is_active"`
	SortOrder  int            `json:"sort_order" gorm:"column:sort_order"`
}

func (CategoryAttributes) TableName() string {
	return "tbl_catageries_attributes"
}

// Request

type CreatedCatRequest struct {
	Name      string `json:"name" binding:"required"`
	SortOrder int    `json:"sort_order" binding:"required"`
	IsActive  bool   `json:"is_active"`
}

type UpdateCatRequest struct {
	CreatedCatRequest
	Id uuid.UUID `json:"id"`
}

type AttributeItemRequest struct {
	Name       string   `json:"name" binding:"required"`
	DataType   string   `json:"data_type" binding:"required"`
	Unit       string   `json:"unit"`
	Options    []string `json:"options"`
	IsRequired bool     `json:"is_required"`
	ShowOnCard bool     `json:"show_on_card"`
	IsActive   bool     `json:"is_active"`
	SortOrder  int      `json:"sort_order"`
}

type CreatedCatAttrRequest struct {
	CategoryId uuid.UUID              `json:"category_id" binding:"required"`
	Attributes []AttributeItemRequest `json:"attributes" binding:"required"`
}

type UpdateCatAttrRequest struct {
	Id         uuid.UUID `json:"id" binding:"required"`
	CategoryId uuid.UUID `json:"category_id" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	DataType   string    `json:"data_type" binding:"required"`
	Unit       string    `json:"unit"`
	Options    []string  `json:"options"`
	IsRequired bool      `json:"is_required"`
	ShowOnCard bool      `json:"show_on_card"`
	IsActive   bool      `json:"is_active"`
	SortOrder  int       `json:"sort_order"`
}

// Response

type CategoryResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
}

type CategoryAttributeResponse struct {
	Id         uuid.UUID `json:"id"`
	CategoryId uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
	DataType   string    `json:"data_type"`
	Unit       string    `json:"unit"`
	Options    []string  `json:"options"`
	IsRequired bool      `json:"is_required"`
	ShowOnCard bool      `json:"show_on_card"`
	IsActive   bool      `json:"is_active"`
	SortOrder  int       `json:"sort_order"`
}
