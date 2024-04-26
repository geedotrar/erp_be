package models

import (
	"errors"
	"time"

	"github.com/geedotrar/erp-api/helper"
	"gorm.io/gorm"
)

type UsersResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    *[]User `json:"data"`
	// Data    []User `json:"data,omitempty"`
	Error bool `json:"error"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    *User  `json:"data"`
	Error   bool   `json:"error"`
}

type User struct {
	ID           uint64         `json:"id" gorm:"primaryKey"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Email        string         `json:"email"`
	Password     string         `json:"-"`
	Role         string         `json:"role"`
	PhoneNumber  string         `json:"phone_number"`
	PositionName string         `json:"position_name"`
	Company      string         `json:"company"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
}

type UserCreateRequest struct { // for validate password
	ID           uint64    `json:"id" gorm:"primaryKey"`
	FirstName    string    `json:"first_name" binding:"required"`
	LastName     string    `json:"last_name" binding:"required"`
	Email        string    `json:"email" binding:"required"`
	Password     string    `json:"password" binding:"required"`
	Role         string    `json:"role" binding:"required"`
	PhoneNumber  string    `json:"phone_number" binding:"required"`
	PositionName string    `json:"position_name" binding:"required"`
	Company      string    `json:"company" binding:"required"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserEditRequest struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Role         string    `json:"role"`
	PhoneNumber  string    `json:"phone_number"`
	PositionName string    `json:"position_name"`
	Company      string    `json:"company"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserView struct {
	ID       uint64    `json:"id"`
	Username string    `json:"username" binding:"required"`
	Password string    `json:"-" binding:"required"`
	Email    string    `json:"email" binding:"required"`
	Dob      time.Time `json:"dob" binding:"required"`
}

type UserSignUp struct {
	// ID       uint64    `json:"id" gorm:"primaryKey"`

	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u UserSignUp) ValidateSignUp() error {
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if !helper.IsValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

func (u UserCreateRequest) ValidateCreate() error {
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if u.Email == "" {
		return errors.New("email cannot be empty")
	}

	if !helper.IsValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

func (u UserEditRequest) ValidateUpdate() error {
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if u.Email == "" {
		return errors.New("email cannot be empty")
	}

	if !helper.IsValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	return nil
}
