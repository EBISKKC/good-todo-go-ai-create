package router

import (
	"net/http"

	"good-todo-go/internal/presentation/public/api"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, api.HealthResponse{Status: "ok"})
}
