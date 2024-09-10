package main

import (
	"github.com/labstack/echo/v4"
	"go-mongoDb-starter/internal/controller"
	"go-mongoDb-starter/internal/repository"
	"go-mongoDb-starter/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type Controller struct {
	UserController *controller.UserController

	// add another controller here
}

// NewController create new controller
func NewController(db *mongo.Database, jwtSecret string) *Controller {
	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo, jwtSecret)
	userController := controller.NewUserController(userService)

	return &Controller{
		UserController: userController,
	}
}

// Routes all routes application API
func (app *Config) Routes(e *echo.Echo, c *Controller, jwtSecret []byte) {
	r := e.Group("/api/v1")
	r.POST("/users/register", c.UserController.RegisterUser)
	r.POST("/users/login", c.UserController.LoginUser)
	// add another routes here

}
