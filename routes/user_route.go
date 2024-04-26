package routes

import (
	"github.com/geedotrar/erp-api/handlers"
	"github.com/geedotrar/erp-api/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter interface {
	Mount()
}

type userRouterImpl struct {
	v       *gin.RouterGroup
	handler handlers.UserHandler
}

func NewUserRouter(v *gin.RouterGroup, handler handlers.UserHandler) UserRouter {
	return &userRouterImpl{v: v, handler: handler}
}

func (u *userRouterImpl) Mount() {
	u.v.POST("/register", u.handler.UserSignUp)
	u.v.POST("/login", u.handler.UserLogin)

	u.v.Use(middleware.CheckAuthBearer)

	u.v.GET("/", u.handler.GetUsers)
	u.v.GET("/:id", u.handler.GetUserByID)

	u.v.POST("/", u.handler.CreateUser)
	u.v.PUT("/:id", u.handler.UpdateUser)
	u.v.DELETE("/:id", u.handler.DeleteUser)

}
