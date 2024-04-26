package repository

import (
	"context"

	"github.com/geedotrar/erp-api/config"
	"github.com/geedotrar/erp-api/models"
	"gorm.io/gorm"
)

type CompanyQuery interface {
	GetCompany(ctx context.Context) ([]models.Company, error)
	GetCompanyByID(ctx context.Context, id uint64) (models.Company, error)
	GetCompanyByCompanyName(ctx context.Context, companyName string) (models.Company, error)

	CreateCompany(ctx context.Context, company models.CompanyRequest) (models.CompanyRequest, error)
	UpdateCompany(ctx context.Context, id uint64, company models.CompanyRequest) (models.CompanyRequest, error)

	DeleteCompany(ctx context.Context, id uint64) error

	// mencari user yang menggunakan companyID
	GetUserByCompanyID(ctx context.Context, companyID uint64) ([]models.User, error)

	GetSoftDeletedCompanies(ctx context.Context) ([]models.Company, error)

	RestoreCompany(ctx context.Context, id uint64) error
}

type companyQueryImpl struct {
	db config.GormPostgres
}

func NewCompanyQuery(db config.GormPostgres) CompanyQuery {
	return &companyQueryImpl{db: db}
}

func (c *companyQueryImpl) GetCompany(ctx context.Context) ([]models.Company, error) {
	db := c.db.GetConnection()
	company := []models.Company{}
	if err := db.
		WithContext(ctx).
		Table("companies").
		Find(&company).Error; err != nil {
		return []models.Company{}, err
	}
	return company, nil
}

func (c *companyQueryImpl) GetCompanyByID(ctx context.Context, id uint64) (models.Company, error) {
	db := c.db.GetConnection()
	company := models.Company{}
	if err := db.
		WithContext(ctx).
		Table("companies").
		Where("id = ?", id).
		Find(&company).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Company{}, nil
		}
		return models.Company{}, err
	}
	return company, nil
}

func (c *companyQueryImpl) GetCompanyByCompanyName(ctx context.Context, companyName string) (models.Company, error) {
	db := c.db.GetConnection()
	company := models.Company{}
	if err := db.WithContext(ctx).Where("company_name = ?", companyName).Find(&company).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Company{}, nil
		}
		return models.Company{}, err
	}
	return company, nil
}

// search user using companyID
func (c *companyQueryImpl) GetUserByCompanyID(ctx context.Context, companyID uint64) ([]models.User, error) {
	db := c.db.GetConnection()
	users := []models.User{}
	if err := db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Find(&users).Error; err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (c *companyQueryImpl) GetSoftDeletedCompanies(ctx context.Context) ([]models.Company, error) {
	db := c.db.GetConnection()
	company := []models.Company{}
	if err := db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Find(&company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

func (c companyQueryImpl) CreateCompany(ctx context.Context, company models.CompanyRequest) (models.CompanyRequest, error) {
	db := c.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("companies").
		Save(&company).Error; err != nil {
		return models.CompanyRequest{}, err
	}
	return company, nil
}

func (c *companyQueryImpl) UpdateCompany(ctx context.Context, id uint64, company models.CompanyRequest) (models.CompanyRequest, error) {
	db := c.db.GetConnection()
	updatedCompany := models.CompanyRequest{}
	if err := db.
		WithContext(ctx).
		Table("companies").
		Where("id = ?", id).
		Updates(&company).
		First(&updatedCompany).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.CompanyRequest{}, nil
		}
	}
	return updatedCompany, nil
}

func (c *companyQueryImpl) DeleteCompany(ctx context.Context, id uint64) error {
	db := c.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("companies").
		Delete(&models.Company{ID: id}).
		Error; err != nil {
		return err
	}
	return nil
}

// func (c *companyQueryImpl) RestoreCompany(ctx context.Context, id uint64) error {
// 	db := c.db.GetConnection()
// 	query := "UPDATE companies SET deleted_at = NULL WHERE id = ?"
// 	if err := db.Exec(query, id).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

func (c *companyQueryImpl) RestoreCompany(ctx context.Context, id uint64) error {
	db := c.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Unscoped().
		Model(&models.Company{}).
		Where("id = ?", id).
		Update("deleted_at", nil).
		Error; err != nil {
		return err
	}
	return nil
}
