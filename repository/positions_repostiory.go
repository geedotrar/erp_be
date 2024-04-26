package repository

import (
	"context"

	"github.com/geedotrar/erp-api/config"
	"github.com/geedotrar/erp-api/models"
	"gorm.io/gorm"
)

type PositionQuery interface {
	GetPosition(ctx context.Context) ([]models.Position, error)
	GetPositionByID(ctx context.Context, id uint64) (models.Position, error)
	GetPositionByPositionName(ctx context.Context, positionName string) (models.Position, error)

	CreatePosition(ctx context.Context, position models.PositionCreateRequest) (models.PositionCreateRequest, error)
	UpdatePosition(ctx context.Context, id uint64, position models.PositionUpdateRequest) (models.PositionUpdateRequest, error)

	DeletePosition(ctx context.Context, id uint64) error

	// check user using positionID
	GetUserByPositionID(ctx context.Context, positionID uint64) ([]models.User, error)

	GetSoftDeletedPosition(ctx context.Context) ([]models.Position, error)
}

type positionQueryImpl struct {
	db config.GormPostgres
}

func NewPositionQuery(db config.GormPostgres) PositionQuery {
	return &positionQueryImpl{db: db}
}

func (p *positionQueryImpl) GetPosition(ctx context.Context) ([]models.Position, error) {
	db := p.db.GetConnection()
	position := []models.Position{}
	if err := db.
		WithContext(ctx).
		Table("positions").
		Find(&position).Error; err != nil {
		return []models.Position{}, err
	}
	return position, nil
}

func (p *positionQueryImpl) GetPositionByID(ctx context.Context, id uint64) (models.Position, error) {
	db := p.db.GetConnection()
	position := models.Position{}
	if err := db.
		WithContext(ctx).
		Table("positions").
		Where("id = ?", id).
		Find(&position).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Position{}, nil
		}
		return models.Position{}, err
	}
	return position, nil
}

func (p *positionQueryImpl) GetPositionByPositionName(ctx context.Context, positionName string) (models.Position, error) {
	db := p.db.GetConnection()
	position := models.Position{}
	if err := db.WithContext(ctx).Where("position_name = ?", positionName).Find(&position).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Position{}, nil
		}
		return models.Position{}, err
	}
	return position, nil
}

// search user using positionID
func (p *positionQueryImpl) GetUserByPositionID(ctx context.Context, positionID uint64) ([]models.User, error) {
	db := p.db.GetConnection()
	users := []models.User{}
	if err := db.WithContext(ctx).
		Where("position_id", positionID).
		Find(&users).Error; err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (p *positionQueryImpl) GetSoftDeletedPosition(ctx context.Context) ([]models.Position, error) {
	db := p.db.GetConnection()
	position := []models.Position{}
	if err := db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Find(&position).Error; err != nil {
		return nil, err
	}
	return position, nil
}

func (p positionQueryImpl) CreatePosition(ctx context.Context, position models.PositionCreateRequest) (models.PositionCreateRequest, error) {
	db := p.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("positions").
		Save(&position).Error; err != nil {
		return models.PositionCreateRequest{}, err
	}
	return position, nil
}

func (p *positionQueryImpl) UpdatePosition(ctx context.Context, id uint64, position models.PositionUpdateRequest) (models.PositionUpdateRequest, error) {
	db := p.db.GetConnection()
	updatedPosition := models.PositionUpdateRequest{}
	if err := db.
		WithContext(ctx).
		Table("positions").
		Where("id = ?", id).
		Updates(&position).
		First(&updatedPosition).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.PositionUpdateRequest{}, nil
		}
	}
	return updatedPosition, nil
}

func (c *positionQueryImpl) DeletePosition(ctx context.Context, id uint64) error {
	db := c.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("positions").
		Delete(&models.Position{ID: id}).
		Error; err != nil {
		return err
	}
	return nil
}
