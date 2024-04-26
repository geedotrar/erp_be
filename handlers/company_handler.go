package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/geedotrar/erp-api/models"
	"github.com/geedotrar/erp-api/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CompanyHandler interface {
	GetCompany(ctx *gin.Context)
	GetCompanyByID(ctx *gin.Context)

	CreateCompany(ctx *gin.Context)
	UpdateCompany(ctx *gin.Context)

	DeleteCompany(ctx *gin.Context)
	RestoreCompany(ctx *gin.Context)
}

type companyHandlerImpl struct {
	svc service.CompanyService
}

func NewCompanyHandler(svc service.CompanyService) CompanyHandler {
	return &companyHandlerImpl{svc: svc}
}

func (u *companyHandlerImpl) GetCompany(ctx *gin.Context) {
	company, err := u.svc.GetCompany(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.CompaniesResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to get company",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// company not found
	if len(company) == 0 {
		ctx.JSON(http.StatusNotFound, models.CompaniesResponse{
			Status:  http.StatusNotFound,
			Message: "company not found",
			Data:    nil,
			Error:   false,
		})
		return
	}
	ctx.JSON(http.StatusOK, models.CompaniesResponse{
		Status:  http.StatusOK,
		Message: "success to get company",
		Data:    &company,
		Error:   false,
	})
}

func (c *companyHandlerImpl) GetCompanyByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.CompanyResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	company, err := c.svc.GetCompanyByID(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.CompanyResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to get company",
			Data:    nil,
			Error:   true,
		})
		return
	}

	if company.ID == 0 {
		ctx.JSON(http.StatusNotFound, models.CompanyResponse{
			Status:  http.StatusNotFound,
			Message: "company not found",
			Data:    nil,
			Error:   false,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.CompanyResponse{
		Status:  http.StatusOK,
		Message: "success to get company",
		Data:    &company,
		Error:   false,
	})
}

func (c *companyHandlerImpl) CreateCompany(ctx *gin.Context) {
	companyCreate := models.CompanyRequest{}
	// check req JSON OK
	if err := ctx.ShouldBindJSON(&companyCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, models.CompanyResponse{
			Status:  http.StatusBadRequest,
			Message: "failed to create company: unable to parse request body",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// validasi data ke request jika req salah atau data tidak diisi
	validate := validator.New()
	if err := validate.Struct(companyCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, models.CompanyResponse{
			Status:  http.StatusBadRequest,
			Message: "failed to create company: invalid data",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// check is data already exists in soft deleted
	softDeleted, err := c.svc.IsCompanySoftDeleted(ctx, companyCreate.CompanyName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.CompanyResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to create company: internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}
	if softDeleted {
		ctx.JSON(http.StatusConflict, models.CompanyResponse{
			Status:  http.StatusConflict,
			Message: "failed to create company: company already exists in soft deleted",
			Data:    nil,
			Error:   true,
		})
		return
	}

	company, err := c.svc.CreateCompany(ctx, companyCreate)
	if err != nil {
		if strings.Contains(err.Error(), "company already exists") {
			ctx.JSON(http.StatusConflict, models.CompanyResponse{
				Status:  http.StatusConflict,
				Message: "failed to create company: company already exists",
				Data:    nil,
				Error:   true,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.CompanyResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to create company: internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.CompanyResponse{
		Status:  http.StatusCreated,
		Message: "company created successfully",
		Data:    company.Data,
		Error:   false,
	})
}

func (c *companyHandlerImpl) UpdateCompany(ctx *gin.Context) {
	// parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.CompanyResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// check id exists
	company, _ := c.svc.GetCompanyByID(ctx, uint64(id))
	if company.ID == 0 {
		ctx.JSON(http.StatusNotFound, models.CompanyResponse{
			Status:  http.StatusNotFound,
			Message: "company not found",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// check req JSON OK
	var companyEdit models.CompanyRequest
	if err := ctx.ShouldBindJSON(&companyEdit); err != nil {
		ctx.JSON(http.StatusBadRequest, models.CompanyResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// validasi data ke request jika req salah atau data tidak diisi
	validate := validator.New()
	if err := validate.Struct(companyEdit); err != nil {
		ctx.JSON(http.StatusBadRequest, models.CompanyResponse{
			Status:  http.StatusBadRequest,
			Message: "failed to update company: invalid data",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// check is data exists in soft deleted
	softDeleted, err := c.svc.IsCompanySoftDeleted(ctx, companyEdit.CompanyName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.CompanyResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to update company: internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}
	if softDeleted {
		ctx.JSON(http.StatusConflict, models.CompanyResponse{
			Status:  http.StatusConflict,
			Message: "failed to update company: company already exists in soft deleted",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// call service to edit company
	updatedCompany, err := c.svc.UpdateCompany(ctx, uint64(id), companyEdit)
	if err != nil {
		if strings.Contains(err.Error(), "company already exists") {
			ctx.JSON(http.StatusConflict, models.CompanyResponse{
				Status:  http.StatusConflict,
				Message: "failed to update company: company already exists",
				Data:    nil,
				Error:   true,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.CompanyResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to update company: internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// Return updated company data
	ctx.JSON(http.StatusOK, models.CompanyResponse{
		Status:  http.StatusOK,
		Message: "updated company successfully",
		Data:    updatedCompany.Data,
		Error:   false,
	})
}

func (c *companyHandlerImpl) DeleteCompany(ctx *gin.Context) {
	// get id company
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.CompanyResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	company, err := c.svc.DeleteCompany(ctx, uint64(id))
	if err != nil {
		if strings.Contains(err.Error(), "company is still in use by user") {
			ctx.JSON(http.StatusBadRequest, models.CompanyResponse{
				Status:  http.StatusBadRequest,
				Message: "failed delete company: company is still in use by user",
				Data:    nil,
				Error:   true,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.CompanyResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed Delete Company: internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}
	if company.ID == 0 {
		ctx.JSON(http.StatusNotFound, models.CompanyResponse{
			Status:  http.StatusNotFound,
			Message: "company not found",
			Data:    nil,
			Error:   false,
		})
		return
	}
	ctx.JSON(http.StatusOK, models.CompanyResponse{
		Status:  http.StatusOK,
		Message: "company deleted successfully",
		Data:    &company,
		Error:   false,
	})
}

func (c *companyHandlerImpl) RestoreCompany(ctx *gin.Context) {
	// Get company ID from URL parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.CompanyResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// Restore company
	err = c.svc.RestoreCompany(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.CompanyResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to restore company: internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// Response success message
	ctx.JSON(http.StatusOK, models.CompanyResponse{
		Status:  http.StatusOK,
		Message: "company restored successfully",
		Data:    nil,
		Error:   false,
	})
}
