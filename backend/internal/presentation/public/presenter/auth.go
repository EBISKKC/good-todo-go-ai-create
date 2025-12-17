package presenter

import (
	"net/http"

	"good-todo-go/internal/presentation/public/api"
	"good-todo-go/internal/usecase/output"

	"github.com/labstack/echo/v4"
)

type IAuthPresenter interface {
	Register(c echo.Context, out *output.RegisterOutput) error
	Login(c echo.Context, out *output.LoginOutput) error
	VerifyEmail(c echo.Context, out *output.VerifyEmailOutput) error
	RefreshToken(c echo.Context, out *output.RefreshTokenOutput) error
}

type AuthPresenter struct{}

func NewAuthPresenter() IAuthPresenter {
	return &AuthPresenter{}
}

func (p *AuthPresenter) Register(c echo.Context, out *output.RegisterOutput) error {
	return c.JSON(http.StatusCreated, api.RegisterResponse{
		UserId:   out.UserID,
		TenantId: out.TenantID,
		Email:    out.Email,
		Message:  out.Message,
	})
}

func (p *AuthPresenter) Login(c echo.Context, out *output.LoginOutput) error {
	return c.JSON(http.StatusOK, api.LoginResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		User: api.UserResponse{
			Id:            out.UserID,
			TenantId:      out.TenantID,
			Email:         out.Email,
			Name:          out.Name,
			Role:          api.UserResponseRole(out.Role),
			EmailVerified: true,
		},
	})
}

func (p *AuthPresenter) VerifyEmail(c echo.Context, out *output.VerifyEmailOutput) error {
	return c.JSON(http.StatusOK, api.VerifyEmailResponse{
		Success: out.Success,
		Message: out.Message,
	})
}

func (p *AuthPresenter) RefreshToken(c echo.Context, out *output.RefreshTokenOutput) error {
	return c.JSON(http.StatusOK, api.RefreshTokenResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	})
}
