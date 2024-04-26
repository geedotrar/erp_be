package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/geedotrar/erp-api/models"
	"github.com/geedotrar/erp-api/pkg/response"
	"github.com/geedotrar/erp-api/service"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUsers(ctx *gin.Context)
	GetUserByID(ctx *gin.Context)

	CreateUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)

	DeleteUser(ctx *gin.Context)

	UserSignUp(ctx *gin.Context)
	UserLogin(ctx *gin.Context)
}

type userHandlerImpl struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) UserHandler {
	return &userHandlerImpl{svc: svc}
}

func (u *userHandlerImpl) GetUsers(ctx *gin.Context) {
	users, err := u.svc.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.UsersResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get users",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// user not found
	if len(users) == 0 {
		ctx.JSON(http.StatusNotFound, models.UsersResponse{
			Status:  http.StatusNotFound,
			Message: "Users not found",
			Data:    nil,
			Error:   false,
		})
		return
	}
	ctx.JSON(http.StatusOK, models.UsersResponse{
		Status:  http.StatusOK,
		Message: "Success to get users",
		Data:    &users,
		Error:   false,
	})
}

func (u *userHandlerImpl) GetUserByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	user, err := u.svc.GetUserByID(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get user",
			Data:    nil,
			Error:   true,
		})
		return
	}

	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, models.UserResponse{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Error:   false,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.UserResponse{
		Status:  http.StatusOK,
		Message: "Success to get user",
		Data:    &user,
		Error:   false,
	})
}

func (u *userHandlerImpl) CreateUser(ctx *gin.Context) {
	userCreate := models.UserCreateRequest{}
	if err := ctx.ShouldBindJSON(&userCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, models.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to create user: unable to parse request body",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// check validation Email and password
	if err := userCreate.ValidateCreate(); err != nil {
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, models.UserResponse{
				Status:  http.StatusUnprocessableEntity,
				Message: "Failed to create user: " + err.Error(),
				Data:    nil,
				Error:   true,
			})
			return
		}
	}

	user, err := u.svc.CreateUser(ctx, userCreate)
	// check email already exist
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			ctx.JSON(http.StatusConflict, models.UserResponse{
				Status:  http.StatusConflict,
				Message: "Failed to create user: email already exists",
				Data:    nil,
				Error:   true,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create user: internal server error or email has been soft deleted",
			Data:    nil,
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.UserResponse{
		Status:  http.StatusCreated,
		Message: "User created successfully",
		Data:    user.Data,
		Error:   false,
	})
}

func (u *userHandlerImpl) UpdateUser(ctx *gin.Context) {
	// parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// check id exists
	user, _ := u.svc.GetUserByID(ctx, uint64(id))
	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, models.UserResponse{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// check req OK
	var userEdit models.UserEditRequest
	if err := ctx.ShouldBindJSON(&userEdit); err != nil {
		ctx.JSON(http.StatusBadRequest, models.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// check validation Email and password
	if err := userEdit.ValidateUpdate(); err != nil {
		if err != nil {
			ctx.JSON(http.StatusConflict, models.UserResponse{
				Status:  http.StatusConflict,
				Message: "Failed to update user: " + err.Error(),
				Data:    nil,
				Error:   true,
			})
			return
		}
	}

	// Call service to edit user
	updatedUser, err := u.svc.UpdateUser(ctx, uint64(id), userEdit)
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			ctx.JSON(http.StatusConflict, models.UserResponse{
				Status:  http.StatusConflict,
				Message: "Failed to update user: email already exists",
				Data:    nil,
				Error:   true,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Status internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}

	// Return updated user data
	ctx.JSON(http.StatusOK, models.UserResponse{
		Status:  http.StatusOK,
		Message: "Updated user successfully",
		Data:    updatedUser.Data,
		Error:   false,
	})
}

func (u *userHandlerImpl) DeleteUser(ctx *gin.Context) {
	// get id user
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	user, err := u.svc.DeleteUser(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed Delete User: internal server error",
			Data:    nil,
			Error:   true,
		})
		return
	}
	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, models.UserResponse{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Error:   false,
		})
		return
	}
	ctx.JSON(http.StatusOK, models.UserResponse{
		Status:  http.StatusOK,
		Message: "User Deleted successfully",
		Data:    &user,
		Error:   false,
	})
}

func (u *userHandlerImpl) UserSignUp(ctx *gin.Context) {
	userSignUp := models.UserSignUp{}
	if err := ctx.Bind(&userSignUp); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}
	if err := userSignUp.ValidateSignUp(); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}
	user, err := u.svc.SignUp(ctx, userSignUp)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, map[string]any{
		"user": user,
	})
}

func (u *userHandlerImpl) UserLogin(ctx *gin.Context) {
	var userLogin models.UserLogin

	if err := ctx.Bind(&userLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	// Memeriksa kredensial pengguna
	user, err := u.svc.CheckCredentials(ctx, userLogin.Email, userLogin.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	// Menghasilkan token akses untuk pengguna yang berhasil login
	token, err := u.svc.GenerateUserAccessToken(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	// Mengirimkan token akses sebagai respons ke klien
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "success authorization", "token": token})
}
