package models

import (
	"time"

	"gorm.io/gorm"
)

type CompaniesResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    *[]Company `json:"data"`
	Error   bool       `json:"error"`
}

type CompanyResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    *Company `json:"data"`
	Error   bool     `json:"error"`
}

type Company struct {
	ID          uint64         `json:"id" gorm:"primaryKey"`
	CompanyName string         `json:"company_name"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
}
type CompanyRequest struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	CompanyName string    `json:"company_name" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at"`
}
