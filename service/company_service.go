package service

import (
	"context"
	"errors"

	"github.com/geedotrar/erp-api/models"
	"github.com/geedotrar/erp-api/repository"
)

type CompanyService interface {
	GetCompany(ctx context.Context) ([]models.Company, error)
	GetCompanyByID(ctx context.Context, id uint64) (models.Company, error)

	CreateCompany(ctx context.Context, createCompany models.CompanyRequest) (models.CompanyResponse, error)
	UpdateCompany(ctx context.Context, id uint64, updateCompany models.CompanyRequest) (models.CompanyResponse, error)

	DeleteCompany(ctx context.Context, id uint64) (models.Company, error)

	IsCompanySoftDeleted(ctx context.Context, companyName string) (bool, error)
	RestoreCompany(ctx context.Context, id uint64) error
}

type companyServiceImpl struct {
	repo repository.CompanyQuery
}

func NewCompanyService(repo repository.CompanyQuery) CompanyService {
	return &companyServiceImpl{repo: repo}
}

func (c *companyServiceImpl) GetCompany(ctx context.Context) ([]models.Company, error) {
	company, err := c.repo.GetCompany(ctx)
	if err != nil {
		return []models.Company{}, err
	}
	return company, nil
}

func (c *companyServiceImpl) GetCompanyByID(ctx context.Context, id uint64) (models.Company, error) {
	company, err := c.repo.GetCompanyByID(ctx, id)
	if err != nil {
		return models.Company{}, err
	}
	return company, nil
}

func (c *companyServiceImpl) CreateCompany(ctx context.Context, createCompany models.CompanyRequest) (models.CompanyResponse, error) {
	// check companyName
	existingCompany, err := c.repo.GetCompanyByCompanyName(ctx, createCompany.CompanyName)
	if err != nil {
		return models.CompanyResponse{}, err
	}
	if existingCompany.ID != 0 {
		return models.CompanyResponse{}, errors.New("company already exists")
	}

	// create req
	company := models.CompanyRequest{
		CompanyName: createCompany.CompanyName,
	}

	// Store company to database
	createdCompany, err := c.repo.CreateCompany(ctx, company)
	if err != nil {
		return models.CompanyResponse{}, err
	}

	// response
	response := models.CompanyResponse{
		Data: &models.Company{
			ID:          createdCompany.ID,
			CompanyName: createdCompany.CompanyName,
		}}
	return response, nil
}

func (c *companyServiceImpl) UpdateCompany(ctx context.Context, id uint64, updateCompany models.CompanyRequest) (models.CompanyResponse, error) {
	existingCompany, err := c.repo.GetCompanyByID(ctx, id)
	if err != nil {
		return models.CompanyResponse{}, err
	}
	if existingCompany.CompanyName != updateCompany.CompanyName {
		// Check if the new company name already exists
		newCompany, err := c.repo.GetCompanyByCompanyName(ctx, updateCompany.CompanyName)
		if err != nil {
			return models.CompanyResponse{}, err
		}

		if newCompany.ID != 0 && newCompany.ID != existingCompany.ID {
			return models.CompanyResponse{}, errors.New("company already exists")
		}
	}

	// update req
	company := models.CompanyRequest{
		CompanyName: updateCompany.CompanyName,
	}

	// Store company to database
	updatedCompany, err := c.repo.UpdateCompany(ctx, id, models.CompanyRequest(company))
	if err != nil {
		return models.CompanyResponse{}, err
	}

	// response
	response := models.CompanyResponse{
		Data: &models.Company{
			ID:          updatedCompany.ID,
			CompanyName: updatedCompany.CompanyName,
		}}
	return response, nil
}

func (c *companyServiceImpl) DeleteCompany(ctx context.Context, id uint64) (models.Company, error) {
	// check if company is using in user
	users, err := c.repo.GetUserByCompanyID(ctx, id)
	if err != nil {
		return models.Company{}, err
	}
	if len(users) > 0 {
		return models.Company{}, errors.New("company is still in use by user")
	}
	company, err := c.repo.GetCompanyByID(ctx, id)
	if err != nil {
		return models.Company{}, err
	}

	if company.ID == 0 {
		return models.Company{}, err
	}

	err = c.repo.DeleteCompany(ctx, id)
	if err != nil {
		return models.Company{}, err
	}
	return company, err
}

func (c *companyServiceImpl) IsCompanySoftDeleted(ctx context.Context, companyName string) (bool, error) {
	companies, err := c.repo.GetSoftDeletedCompanies(ctx)
	if err != nil {
		return false, err
	}
	for _, company := range companies {
		if company.CompanyName == companyName {
			return true, nil
		}
	}
	return false, nil
}

func (c *companyServiceImpl) RestoreCompany(ctx context.Context, id uint64) error {
	// Restore soft deleted company
	err := c.repo.RestoreCompany(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
