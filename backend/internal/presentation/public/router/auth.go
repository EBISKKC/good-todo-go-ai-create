package router

import (
	"good-todo-go/internal/presentation/public/api"

	"github.com/labstack/echo/v4"
)

func (s *Server) Register(ctx echo.Context) error {
	var req api.RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	return s.authController.Register(ctx, req)
}

func (s *Server) Login(ctx echo.Context) error {
	var req api.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	return s.authController.Login(ctx, req)
}

func (s *Server) VerifyEmail(ctx echo.Context) error {
	var req api.VerifyEmailRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	return s.authController.VerifyEmail(ctx, req)
}

func (s *Server) RefreshToken(ctx echo.Context) error {
	var req api.RefreshTokenRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	return s.authController.RefreshToken(ctx, req)
}
