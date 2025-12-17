package controller

import (
	"net/http"

	"good-todo-go/internal/presentation/public/api"
	"good-todo-go/internal/presentation/public/presenter"
	"good-todo-go/internal/presentation/public/router/context_keys"
	"good-todo-go/internal/usecase"
	"good-todo-go/internal/usecase/input"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userUsecase   usecase.IUserInteractor
	userPresenter presenter.IUserPresenter
}

func NewUserController(
	userUsecase usecase.IUserInteractor,
	userPresenter presenter.IUserPresenter,
) *UserController {
	return &UserController{
		userUsecase:   userUsecase,
		userPresenter: userPresenter,
	}
}

func (ctrl *UserController) GetMe(c echo.Context) error {
	userID, ok := c.Get(context_keys.UserIDContextKey).(string)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	out, err := ctrl.userUsecase.GetMe(c.Request().Context(), userID)
	if err != nil {
		if err == usecase.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctrl.userPresenter.GetMe(c, out)
}

func (ctrl *UserController) UpdateMe(c echo.Context, req api.UpdateUserRequest) error {
	userID, ok := c.Get(context_keys.UserIDContextKey).(string)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	out, err := ctrl.userUsecase.UpdateMe(c.Request().Context(), userID, &input.UpdateUserInput{
		Name: req.Name,
	})
	if err != nil {
		if err == usecase.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctrl.userPresenter.UpdateMe(c, out)
}
