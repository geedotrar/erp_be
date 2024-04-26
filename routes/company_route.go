package routes

import (
	"github.com/geedotrar/erp-api/handlers"
	"github.com/gin-gonic/gin"
)

type CompanyRouter interface {
	Mount()
}

type companyRouterImpl struct {
	v       *gin.RouterGroup
	handler handlers.CompanyHandler
}

func NewCompanyRouter(v *gin.RouterGroup, handler handlers.CompanyHandler) CompanyRouter {
	return &companyRouterImpl{v: v, handler: handler}
}

func (c *companyRouterImpl) Mount() {
	// u.v.Use(middleware.CheckAuthBearer)

	c.v.GET("/", c.handler.GetCompany)
	c.v.GET("/:id", c.handler.GetCompanyByID)

	c.v.POST("/create", c.handler.CreateCompany)
	c.v.PUT("/update/:id", c.handler.UpdateCompany)
	c.v.DELETE("/delete/:id", c.handler.DeleteCompany)
	c.v.PUT("/restore/:id", c.handler.RestoreCompany)

}
