package service

import (
	"context"
	"errors"

	"github.com/geedotrar/erp-api/models"
	"github.com/geedotrar/erp-api/repository"
)

type PositionService interface {
	GetPosition(ctx context.Context) ([]models.Position, error)
	GetPositionByID(ctx context.Context, id uint64) (models.Position, error)

	CreatePosition(ctx context.Context, createPosition models.PositionCreateRequest) (models.PositionResponse, error)
	UpdatePosition(ctx context.Context, id uint64, UpdatePosition models.PositionUpdateRequest) (models.PositionResponse, error)

	DeletePosition(ctx context.Context, id uint64) (models.Position, error)

	IsPositionSoftDeleted(ctx context.Context, positionName string, positionCode string) (bool, error)
}

type positionServiceImpl struct {
	repo repository.PositionQuery
}

func NewPositionService(repo repository.PositionQuery) PositionService {
	return &positionServiceImpl{repo: repo}
}

func (p *positionServiceImpl) GetPosition(ctx context.Context) ([]models.Position, error) {
	position, err := p.repo.GetPosition(ctx)
	if err != nil {
		return []models.Position{}, err
	}
	return position, nil
}

func (p *positionServiceImpl) GetPositionByID(ctx context.Context, id uint64) (models.Position, error) {
	position, err := p.repo.GetPositionByID(ctx, id)
	if err != nil {
		return models.Position{}, err
	}
	return position, nil
}

func (p *positionServiceImpl) CreatePosition(ctx context.Context, createPosition models.PositionCreateRequest) (models.PositionResponse, error) {
	// check positionName
	existingPosition, err := p.repo.GetPositionByPositionName(ctx, createPosition.PositionName)
	if err != nil {
		return models.PositionResponse{}, err
	}
	if existingPosition.ID != 0 {
		return models.PositionResponse{}, errors.New("position already exists")
	}

	// create req
	position := models.PositionCreateRequest{
		PositionName: createPosition.PositionName,
		PositionCode: createPosition.PositionCode,
	}

	// Store position to database
	createdPosition, err := p.repo.CreatePosition(ctx, position)
	if err != nil {
		return models.PositionResponse{}, err
	}

	// response
	response := models.PositionResponse{
		Data: &models.Position{
			ID:           createdPosition.ID,
			PositionName: createdPosition.PositionName,
			PositionCode: createdPosition.PositionCode,
		}}
	return response, nil
}

func (p *positionServiceImpl) UpdatePosition(ctx context.Context, id uint64, updatePosition models.PositionUpdateRequest) (models.PositionResponse, error) {
	// check email
	existingPosition, err := p.repo.GetPositionByPositionName(ctx, updatePosition.PositionName)
	if err != nil {
		return models.PositionResponse{}, err
	}
	if existingPosition.ID != 0 && existingPosition.ID != id {
		return models.PositionResponse{}, errors.New("position already exists")
	}

	// update req
	position := models.PositionUpdateRequest{
		PositionName: updatePosition.PositionName,
		PositionCode: updatePosition.PositionCode,
	}

	// Store position to database
	updatedPosition, err := p.repo.UpdatePosition(ctx, id, models.PositionUpdateRequest(position))
	if err != nil {
		return models.PositionResponse{}, err
	}

	// response
	response := models.PositionResponse{
		Data: &models.Position{
			ID:           updatedPosition.ID,
			PositionName: updatedPosition.PositionName,
			PositionCode: updatedPosition.PositionCode,
		}}
	return response, nil
}

func (p *positionServiceImpl) DeletePosition(ctx context.Context, id uint64) (models.Position, error) {
	// check position is using in user
	users, err := p.repo.GetUserByPositionID(ctx, id)
	if err != nil {
		return models.Position{}, err
	}
	if len(users) > 0 {
		return models.Position{}, errors.New("position is still in use by users")
	}

	position, err := p.repo.GetPositionByID(ctx, id)
	if err != nil {
		return models.Position{}, err
	}

	if position.ID == 0 {
		return models.Position{}, err
	}

	err = p.repo.DeletePosition(ctx, id)
	if err != nil {
		return models.Position{}, err
	}
	return position, err
}

func (p *positionServiceImpl) IsPositionSoftDeleted(ctx context.Context, positonName string, positionCode string) (bool, error) {
	positions, err := p.repo.GetSoftDeletedPosition(ctx)
	if err != nil {
		return false, err
	}
	for _, position := range positions {
		if position.PositionName == positonName || position.PositionCode == positionCode {
			return true, nil
		}
	}
	return false, nil
}
