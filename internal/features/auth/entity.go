package auth

import "tjus-be/internal/pkg/entity"

type AdminUser struct {
	entity.BaseEntity
	Email        string  `json:"email" gorm:"column:email;uniqueIndex;not null"`
	PasswordHash string  `json:"-" gorm:"column:password_hash"`
	Name         string  `json:"name" gorm:"column:name"`
	Token        *string `json:"token" gorm:"column:token"`
}

func (AdminUser) TableName() string {
	return "tbl_admin_users"
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LogoutRequest struct {
	Token string `json:"token"`
}
