package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/geedotrar/erp-api/helper"
	"github.com/geedotrar/erp-api/models"
	"github.com/geedotrar/erp-api/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uint64) (models.User, error)

	CreateUser(ctx context.Context, createUser models.UserCreateRequest) (models.UserResponse, error)
	UpdateUser(ctx context.Context, id uint64, updateUser models.UserEditRequest) (models.UserResponse, error)

	DeleteUser(ctx context.Context, id uint64) (models.User, error)

	SignUp(ctx context.Context, userSignUp models.UserSignUp) (models.UserView, error)
	GenerateUserAccessToken(ctx context.Context, user models.User) (token string, err error)
	CheckCredentials(ctx context.Context, email string, password string) (models.User, error)
}

type userServiceImpl struct {
	repo repository.UserQuery
}

func NewUserService(repo repository.UserQuery) UserService {
	return &userServiceImpl{repo: repo}
}

func (u *userServiceImpl) GetUsers(ctx context.Context) ([]models.User, error) {
	users, err := u.repo.GetUsers(ctx)
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (u *userServiceImpl) GetUserByID(ctx context.Context, id uint64) (models.User, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *userServiceImpl) CreateUser(ctx context.Context, createUser models.UserCreateRequest) (models.UserResponse, error) {
	// check email
	existingUser, err := u.repo.GetUserByEmail(ctx, createUser.Email)
	if err != nil {
		return models.UserResponse{}, err
	}
	if existingUser.ID != 0 {
		return models.UserResponse{}, errors.New("email already exists")
	}

	// create req
	user := models.UserCreateRequest{
		FirstName:    createUser.FirstName,
		LastName:     createUser.LastName,
		Email:        createUser.Email,
		Role:         createUser.Role,
		PhoneNumber:  createUser.PhoneNumber,
		PositionName: createUser.PositionName,
		Company:      createUser.Company,
	}

	// Hash password
	pass, err := helper.GenerateHash(createUser.Password)
	if err != nil {
		return models.UserResponse{}, err
	}
	user.Password = pass

	// Store user to database
	createdUser, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return models.UserResponse{}, err
	}

	// response
	response := models.UserResponse{
		Data: &models.User{
			ID:           createdUser.ID,
			FirstName:    createdUser.FirstName,
			LastName:     createdUser.LastName,
			Email:        createdUser.Email,
			Role:         createdUser.Role,
			PhoneNumber:  createdUser.PhoneNumber,
			PositionName: createdUser.PositionName,
			Company:      createdUser.Company,
		}}
	return response, nil
}

func (u *userServiceImpl) UpdateUser(ctx context.Context, id uint64, updateUser models.UserEditRequest) (models.UserResponse, error) {
	// check email
	existingUser, err := u.repo.GetUserByEmail(ctx, updateUser.Email)
	if err != nil {
		return models.UserResponse{}, err
	}
	if existingUser.ID != 0 && existingUser.ID != id {
		return models.UserResponse{}, errors.New("email already exists")
	}

	// update req
	user := models.UserCreateRequest{
		FirstName:    updateUser.FirstName,
		LastName:     updateUser.LastName,
		Email:        updateUser.Email,
		Role:         updateUser.Role,
		PhoneNumber:  updateUser.PhoneNumber,
		PositionName: updateUser.PositionName,
		Company:      updateUser.Company,
	}

	// Hash password
	pass, err := helper.GenerateHash(updateUser.Password)
	if err != nil {
		return models.UserResponse{}, err
	}
	user.Password = pass

	// Store user to database

	updatedUser, err := u.repo.UpdateUser(ctx, id, models.UserEditRequest(user))
	if err != nil {
		return models.UserResponse{}, err
	}

	// response
	response := models.UserResponse{
		Data: &models.User{
			ID:           updatedUser.ID,
			FirstName:    updatedUser.FirstName,
			LastName:     updatedUser.LastName,
			Email:        updatedUser.Email,
			Role:         updatedUser.Role,
			PhoneNumber:  updatedUser.PhoneNumber,
			PositionName: updatedUser.PositionName,
			Company:      updatedUser.Company,
		}}
	return response, nil
}

func (u *userServiceImpl) DeleteUser(ctx context.Context, id uint64) (models.User, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}

	if user.ID == 0 {
		return models.User{}, err
	}

	err = u.repo.DeleteUser(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, err
}

func (u *userServiceImpl) SignUp(ctx context.Context, userSignUp models.UserSignUp) (models.UserView, error) {
	user := models.User{
		Email: userSignUp.Email,
	}
	// encryption password
	// hashing
	pass, err := helper.GenerateHash(userSignUp.Password)
	if err != nil {
		return models.UserView{}, err
	}
	user.Password = pass

	getUserByEmail, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return models.UserView{}, err
	}
	if getUserByEmail.Email == user.Email {
		return models.UserView{}, errors.New("email already exist")
	}

	// store to db
	createdUser, err := u.repo.SignUp(ctx, user)
	if err != nil {
		return models.UserView{}, err
	}
	printUser := models.UserView{
		ID:    createdUser.ID,
		Email: createdUser.Email,
	}

	return printUser, err
}

func (u *userServiceImpl) GenerateUserAccessToken(ctx context.Context, user models.User) (token string, err error) {
	// generate claim
	now := time.Now()

	claim := models.StandardClaim{
		Jti: fmt.Sprintf("%v", time.Now().UnixNano()),
		Iss: "go-middleware",
		Aud: "golang-006",
		Sub: "access-token",
		Exp: uint64(now.Add(time.Hour).Unix()),
		Iat: uint64(now.Unix()),
		Nbf: uint64(now.Unix()),
	}

	userClaim := models.AccessClaim{
		StandardClaim: claim,
		UserID:        user.ID,
	}

	token, err = helper.GenerateToken(userClaim)
	return
}

func (u *userServiceImpl) CheckCredentials(ctx context.Context, email string, password string) (models.User, error) {
	// Retrieve user by email
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, err
	}

	// Check if user exists
	if user.ID == 0 {
		return models.User{}, errors.New("user not found")
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, err
	}

	// Credentials are correct, return user
	return user, nil
}

// func (u *userServiceImpl) CheckSoftDeletedUserByEmail(ctx context.Context, email string) bool {
// 	user, err := u.repo.CheckSoftDeletedUserByEmail(ctx, email)
// 	if err != nil {
// 		return false
// 	}
// 	return user.DeletedAt.Valid
// }
