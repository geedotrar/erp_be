package repository

import (
	"context"

	"github.com/geedotrar/erp-api/config"
	"github.com/geedotrar/erp-api/models"
	"gorm.io/gorm"
)

type UserQuery interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uint64) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)

	CheckSoftDeletedUserByEmail(ctx context.Context, email string) bool

	CreateUser(ctx context.Context, user models.UserCreateRequest) (models.UserCreateRequest, error)
	UpdateUser(ctx context.Context, id uint64, user models.UserEditRequest) (models.UserEditRequest, error)

	DeleteUser(ctx context.Context, id uint64) error

	SignUp(ctx context.Context, user models.User) (models.User, error)
}

type userQueryImpl struct {
	db config.GormPostgres
}

func NewUserQuery(db config.GormPostgres) UserQuery {
	return &userQueryImpl{db: db}
}

func (u *userQueryImpl) GetUsers(ctx context.Context) ([]models.User, error) {
	db := u.db.GetConnection()
	users := []models.User{}
	if err := db.
		WithContext(ctx).
		Table("users").
		Find(&users).Error; err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (u *userQueryImpl) GetUserByID(ctx context.Context, id uint64) (models.User, error) {
	db := u.db.GetConnection()
	users := models.User{}
	if err := db.
		WithContext(ctx).
		Table("users").
		Where("id = ?", id).
		Find(&users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return users, nil
}

func (u *userQueryImpl) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	db := u.db.GetConnection()
	user := models.User{}
	if err := db.WithContext(ctx).Where("email = ?", email).Find(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return user, nil
}

func (u *userQueryImpl) CreateUser(ctx context.Context, user models.UserCreateRequest) (models.UserCreateRequest, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		Save(&user).Error; err != nil {
		return models.UserCreateRequest{}, err
	}
	return user, nil
}

func (u *userQueryImpl) UpdateUser(ctx context.Context, id uint64, user models.UserEditRequest) (models.UserEditRequest, error) {
	db := u.db.GetConnection()
	updatedUser := models.UserEditRequest{}
	if err := db.
		WithContext(ctx).
		Table("users").
		Where("id = ?", id).
		Updates(&user).
		First(&updatedUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.UserEditRequest{}, nil
		}
	}
	return updatedUser, nil
}

func (u *userQueryImpl) DeleteUser(ctx context.Context, id uint64) error {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		Delete(&models.User{ID: id}).
		Error; err != nil {
		return err
	}
	return nil
}

func (u *userQueryImpl) SignUp(ctx context.Context, user models.User) (models.User, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		Save(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *userQueryImpl) CheckSoftDeletedUserByEmail(ctx context.Context, email string) bool {
	db := u.db.GetConnection()
	var count int64
	if err := db.WithContext(ctx).
		Table("users").
		Where("email = ?", email).
		Where("deleted_at IS NOT NULL").
		Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
