package presenter

import (
	"net/http"

	"good-todo-go/internal/presentation/public/api"
	"good-todo-go/internal/usecase/output"

	"github.com/labstack/echo/v4"
)

type IUserPresenter interface {
	GetMe(c echo.Context, out *output.UserOutput) error
	UpdateMe(c echo.Context, out *output.UserOutput) error
}

type UserPresenter struct{}

func NewUserPresenter() IUserPresenter {
	return &UserPresenter{}
}

func (p *UserPresenter) GetMe(c echo.Context, out *output.UserOutput) error {
	return c.JSON(http.StatusOK, toUserResponse(out))
}

func (p *UserPresenter) UpdateMe(c echo.Context, out *output.UserOutput) error {
	return c.JSON(http.StatusOK, toUserResponse(out))
}

func toUserResponse(out *output.UserOutput) api.UserResponse {
	return api.UserResponse{
		Id:            out.ID,
		TenantId:      out.TenantID,
		Email:         out.Email,
		Name:          out.Name,
		Role:          api.UserResponseRole(out.Role),
		EmailVerified: out.EmailVerified,
		CreatedAt:     out.CreatedAt,
		UpdatedAt:     out.UpdatedAt,
	}
}
