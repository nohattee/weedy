package v1

import (
	"context"
	"fmt"
	"net/http"

	"weedy/internal/entity"
	"weedy/internal/repository"
	"weedy/pkg/httpserver"
)

type userController struct {
	UserRepo interface {
		List(ctx context.Context, params *repository.ListParams) ([]*entity.User, error)
	}
}

type listUsersRequest struct {
}

func NewUserController(userRepo *repository.UserRepo) httpserver.Controller {
	return &userController{UserRepo: userRepo}
}

func (c *userController) Routes() []httpserver.Route {
	return []httpserver.Route{
		{
			Name:    "",
			Path:    "/",
			Method:  http.MethodGet,
			Request: &listUsersRequest{},
			Handler: c.List,
		},
	}
}

func (c *userController) List(ctx context.Context) httpserver.Response {
	users, err := c.UserRepo.List(ctx, &repository.ListParams{})
	if err != nil {
		return httpserver.Response{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("c.UserRepo.List err: %v", err),
		}
	}
	return httpserver.Response{
		StatusCode: http.StatusOK,
		Data:       users,
	}
}
