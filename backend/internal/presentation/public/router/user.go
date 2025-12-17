package router

import (
	"good-todo-go/internal/presentation/public/api"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetMe(ctx echo.Context) error {
	return s.userController.GetMe(ctx)
}

func (s *Server) UpdateMe(ctx echo.Context) error {
	var req api.UpdateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	return s.userController.UpdateMe(ctx, req)
}
