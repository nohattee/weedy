package v1

import (
	"context"
	"database/sql"
	"net/http"

	"weedy/pkg/httpserver"
)

type userController struct {
	UserUseCase interface {
	}
}

type listUsersRequest struct {
}

func NewUserController(db *sql.DB) httpserver.Controller {
	return &userController{}
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
	return httpserver.Response{
		StatusCode: http.StatusOK,
	}
}
