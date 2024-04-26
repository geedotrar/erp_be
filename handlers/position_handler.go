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

type PositionHandler interface {
	GetPosition(ctx *gin.Context)
	GetPositionByID(ctx *gin.Context)

	CreatePosition(ctx *gin.Context)
	UpdatePosition(ctx *gin.Context)

	DeletePosition(ctx *gin.Context)
}

type positionHandlerImpl struct {
	svc service.PositionService
}

func NewPositionHandler(svc service.PositionService) PositionHandler {
	return &positionHandlerImpl{svc: svc}
}

func (p *positionHandlerImpl) GetPosition(ctx *gin.Context) {
	position, err := p.svc.GetPosition(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.PositionsResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to get position",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// position not found
	if len(position) == 0 {
		ctx.JSON(http.StatusNotFound, models.PositionsResponse{
			Status:  http.StatusNotFound,
			Message: "position not found",
			Data:    nil,
			Error:   false,
		})
		return
	}
	ctx.JSON(http.StatusOK, models.PositionsResponse{
		Status:  http.StatusOK,
		Message: "success to get position",
		Data:    &position,
		Error:   false,
	})
}

func (p *positionHandlerImpl) GetPositionByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.PositionResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	position, err := p.svc.GetPositionByID(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.PositionResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to get position",
			Data:    nil,
			Error:   true,
		})
		return
	}

	if position.ID == 0 {
		ctx.JSON(http.StatusNotFound, models.PositionResponse{
			Status:  http.StatusNotFound,
			Message: "position not found",
			Data:    nil,
			Error:   false,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.PositionResponse{
		Status:  http.StatusOK,
		Message: "success to get position",
		Data:    &position,
		Error:   false,
	})
}

func (p *positionHandlerImpl) CreatePosition(ctx *gin.Context) {
	positionCreate := models.PositionCreateRequest{}
	if err := ctx.ShouldBindJSON(&positionCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, models.PositionResponse{
			Status:  http.StatusBadRequest,
			Message: "failed to create position: unable to parse request body",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// validasi data ke request jika req salah atau data tidak diisi
	validate := validator.New()
	if err := validate.Struct(positionCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, models.PositionResponse{
			Status:  http.StatusBadRequest,
			Message: "failed to create position: invalid data",
			Data:    nil,
			Error:   true,
		})
		return
	}

	softDeleted, err := p.svc.IsPositionSoftDeleted(ctx, positionCreate.PositionName, positionCreate.PositionCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.PositionResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to create position: internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}
	if softDeleted {
		ctx.JSON(http.StatusConflict, models.PositionResponse{
			Status:  http.StatusConflict,
			Message: "failed to create position: position name or position code already exists in soft deleted",
			Data:    nil,
			Error:   true,
		})
		return
	}

	position, err := p.svc.CreatePosition(ctx, positionCreate)
	if err != nil {
		if strings.Contains(err.Error(), "position already exists") {
			ctx.JSON(http.StatusConflict, models.PositionResponse{
				Status:  http.StatusConflict,
				Message: "failed to create position: position already exists",
				Data:    nil,
				Error:   true,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.PositionResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to create position: internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.PositionResponse{
		Status:  http.StatusCreated,
		Message: "position created successfully",
		Data:    position.Data,
		Error:   false,
	})
}

func (p *positionHandlerImpl) UpdatePosition(ctx *gin.Context) {
	// parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.PositionResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// check id exists
	position, _ := p.svc.GetPositionByID(ctx, uint64(id))
	if position.ID == 0 {
		ctx.JSON(http.StatusNotFound, models.PositionResponse{
			Status:  http.StatusNotFound,
			Message: "position not found",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// check req OK
	var positionEdit models.PositionUpdateRequest
	if err := ctx.ShouldBindJSON(&positionEdit); err != nil {
		ctx.JSON(http.StatusBadRequest, models.PositionResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// validasi data ke request jika req salah atau data tidak diisi
	validate := validator.New()
	if err := validate.Struct(positionEdit); err != nil {
		ctx.JSON(http.StatusBadRequest, models.PositionResponse{
			Status:  http.StatusBadRequest,
			Message: "failed to update position: invalid data",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// check is data exists in soft deleted
	softDeleted, err := p.svc.IsPositionSoftDeleted(ctx, positionEdit.PositionName, positionEdit.PositionCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.PositionResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to update position: internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}
	if softDeleted {
		ctx.JSON(http.StatusConflict, models.PositionResponse{
			Status:  http.StatusConflict,
			Message: "failed to update position: position name or position code already exists in soft deleted",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// Call service to edit position
	updatedPosition, err := p.svc.UpdatePosition(ctx, uint64(id), positionEdit)
	if err != nil {
		if strings.Contains(err.Error(), "position already exists") {
			ctx.JSON(http.StatusConflict, models.PositionResponse{
				Status:  http.StatusConflict,
				Message: "failed to update position: position already exists",
				Data:    nil,
				Error:   true,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.PositionResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to update position: internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// Return updated position data
	ctx.JSON(http.StatusOK, models.PositionResponse{
		Status:  http.StatusOK,
		Message: "updated position successfully",
		Data:    updatedPosition.Data,
		Error:   false,
	})
}

func (p *positionHandlerImpl) DeletePosition(ctx *gin.Context) {
	// get id position
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.PositionResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	position, err := p.svc.DeletePosition(ctx, uint64(id))
	if err != nil {
		if strings.Contains(err.Error(), "position is still in use by user") {
			ctx.JSON(http.StatusBadRequest, models.PositionResponse{
				Status:  http.StatusBadRequest,
				Message: "failed delete position: position is still in use by user",
				Data:    nil,
				Error:   true,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.PositionResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed Delete Position: internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}

	if position.ID == 0 {
		ctx.JSON(http.StatusNotFound, models.PositionResponse{
			Status:  http.StatusNotFound,
			Message: "position not found",
			Data:    nil,
			Error:   false,
		})
		return
	}
	ctx.JSON(http.StatusOK, models.PositionResponse{
		Status:  http.StatusOK,
		Message: "position Deleted successfully",
		Data:    &position,
		Error:   false,
	})
}
