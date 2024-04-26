package main

import (
	"log"

	"github.com/geedotrar/erp-api/config"
	"github.com/geedotrar/erp-api/handlers"
	"github.com/geedotrar/erp-api/repository"
	"github.com/geedotrar/erp-api/routes"
	"github.com/geedotrar/erp-api/service"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {
	server()
}

func server() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	g := gin.Default()
	g.Use(gin.Recovery())

	usersGroup := g.Group("/users")
	gorm := config.NewGormPostgres()
	userRepo := repository.NewUserQuery(gorm)
	userSvc := service.NewUserService(userRepo)
	userHdl := handlers.NewUserHandler(userSvc)
	userRouter := routes.NewUserRouter(usersGroup, userHdl)
	userRouter.Mount()

	companyGroup := g.Group("/company")
	companyRepo := repository.NewCompanyQuery(gorm)
	companySvc := service.NewCompanyService(companyRepo)
	companyHdl := handlers.NewCompanyHandler(companySvc)
	companyRouter := routes.NewCompanyRouter(companyGroup, companyHdl)
	companyRouter.Mount()

	positionGroup := g.Group("/positions")
	positionRepo := repository.NewPositionQuery(gorm)
	positionSvc := service.NewPositionService(positionRepo)
	positionHdl := handlers.NewPositionHandler(positionSvc)
	positionRouter := routes.NewPositionRouter(positionGroup, positionHdl)
	positionRouter.Mount()

	g.Run(":8080")
}
