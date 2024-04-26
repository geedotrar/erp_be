package routes

import (
	"github.com/geedotrar/erp-api/handlers"
	"github.com/gin-gonic/gin"
)

type PositionRouter interface {
	Mount()
}

type positionRouterImpl struct {
	v       *gin.RouterGroup
	handler handlers.PositionHandler
}

func NewPositionRouter(v *gin.RouterGroup, handler handlers.PositionHandler) PositionRouter {
	return &positionRouterImpl{v: v, handler: handler}
}

func (p *positionRouterImpl) Mount() {
	// u.v.Use(middleware.CheckAuthBearer)

	p.v.GET("/", p.handler.GetPosition)
	p.v.GET("/:id", p.handler.GetPositionByID)

	p.v.POST("/create", p.handler.CreatePosition)
	p.v.PUT("/update/:id", p.handler.UpdatePosition)
	p.v.DELETE("/delete/:id", p.handler.DeletePosition)

}
