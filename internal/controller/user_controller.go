package controller

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go-mongoDb-starter/internal/models"
	"go-mongoDb-starter/internal/service"
	"go-mongoDb-starter/pkg"
	"net/http"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) RegisterUser(ctx echo.Context) error {
	var req models.RegisterInput
	if err := ctx.Bind(&req); err != nil {
		return pkg.ResponseJson(ctx, http.StatusBadRequest, nil, "invalid request format : "+err.Error())
	}

	// validate request
	if err := req.Validate(); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			formatterErrors := pkg.FormatValidationError(&req, validationErrors)
			return pkg.ResponseJson(ctx, http.StatusBadRequest, nil, formatterErrors)
		}
		return pkg.ResponseJson(ctx, http.StatusBadRequest, nil, err.Error())
	}

	user, err := c.userService.RegisterUser(ctx.Request().Context(), req)
	if err != nil {
		return pkg.ResponseJson(ctx, http.StatusInternalServerError, nil, err.Error())
	}

	response := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"address":  user.Address,
		"age":      user.Age,
		"phone":    user.Phone,
	}

	return pkg.ResponseJson(ctx, http.StatusCreated, response, "user created successfully")
}

func (c *UserController) LoginUser(ctx echo.Context) error {
	var req models.LoginInput
	if err := ctx.Bind(&req); err != nil {
		return pkg.ResponseJson(ctx, http.StatusBadRequest, nil, "invalid request format : "+err.Error())
	}

	// validate request
	if err := req.Validate(); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			formatterErrors := pkg.FormatValidationError(&req, validationErrors)
			return pkg.ResponseJson(ctx, http.StatusBadRequest, nil, formatterErrors)
		}
		return pkg.ResponseJson(ctx, http.StatusBadRequest, nil, err.Error())
	}

	user, err := c.userService.LoginUser(ctx.Request().Context(), &req)
	if err != nil {
		return pkg.ResponseJson(ctx, http.StatusInternalServerError, nil, err.Error())
	}

	response := models.UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Address:  user.Address,
		Age:      user.Age,
		Phone:    user.Phone,
	}

	response.UserToken.Token = user.Token

	return pkg.ResponseJson(ctx, http.StatusOK, response, "user login successfully")
}
