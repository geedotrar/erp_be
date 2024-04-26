package models

import (
	"time"

	"gorm.io/gorm"
)

type PositionsResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    *[]Position `json:"data"`
	Error   bool        `json:"error"`
}

type PositionResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    *Position `json:"data"`
	Error   bool      `json:"error"`
}

type Position struct {
	ID           uint64         `json:"id" gorm:"primaryKey"`
	PositionName string         `json:"position_name"`
	PositionCode string         `json:"position_code"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
}
type PositionCreateRequest struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	PositionName string    `json:"position_name" validate:"required"`
	PositionCode string    `json:"position_code" validate:"required"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PositionUpdateRequest struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	PositionName string    `json:"position_name,omitempty"`
	PositionCode string    `json:"position_code,omitempty"`
	UpdatedAt    time.Time `json:"updated_at"`
}
